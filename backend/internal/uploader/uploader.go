package uploader

import (
	"context"
	"fmt"
	"math"
	"os/exec"
	"time"

	"go.uber.org/zap"

	"github.com/dojo-harvester/backend/internal/config"
	"github.com/dojo-harvester/backend/internal/logger"
	"github.com/dojo-harvester/backend/internal/storage"
)

type Manager struct {
	store *storage.Store
	cfg   *config.Config
}

func NewManager(store *storage.Store, cfg *config.Config) *Manager {
	return &Manager{
		store: store,
		cfg:   cfg,
	}
}

// StartJob initiates the upload process, moving from pending to in_progress
func (m *Manager) StartJob(jobID string) error {
	job, err := m.store.GetUploadJob(jobID)
	if err != nil || job == nil {
		return fmt.Errorf("job not found")
	}

	dataset, err := m.store.GetDataset(job.DatasetID)
	if err != nil || dataset == nil {
		return fmt.Errorf("dataset not found")
	}

	m.store.RemoveFromQueue("pending", jobID)
	m.store.AddToQueue("in_progress", jobID)

	now := time.Now().UTC()
	job.Status = "in_progress"
	job.StartedAt = &now
	m.store.SaveUploadJob(job)

	// Start async processing
	go m.processUpload(job, dataset.StoragePath)

	return nil
}

func (m *Manager) processUpload(job *storage.UploadJob, datasetPath string) {
	// In a real implementation this would stream chunks using go-resty
	// or call the external upload_local.sh script.
	// For this task, we will simulate using the external script logic

	scriptPath := m.cfg.Scripts.UploadLocal
	args := []string{job.JobID, datasetPath, m.cfg.Platform.URL}

	retryLimit := m.cfg.Platform.RetryCount
	if retryLimit == 0 {
		retryLimit = 5
	}

	for attempt := 0; attempt <= retryLimit; attempt++ {
		if attempt > 0 {
			// Exponential backoff: 1s, 2s, 4s, 8s, 16s...
			backoff := time.Duration(math.Pow(2, float64(attempt-1))) * time.Second
			logger.Get().Warn("Retrying upload", zap.String("job_id", job.JobID), zap.Int("attempt", attempt), zap.Duration("backoff", backoff))
			time.Sleep(backoff)
			
			job.RetryCount = attempt
			now := time.Now().UTC()
			job.ResumedAt = &now
			m.store.SaveUploadJob(job)
		}

		ctx, cancel := context.WithTimeout(context.Background(), 1*time.Hour)
		defer cancel()

		cmd := exec.CommandContext(ctx, scriptPath, args...)
		
		err := cmd.Run()
		
		if err == nil {
			// Success
			now := time.Now().UTC()
			job.Status = "completed"
			job.CompletedAt = &now
			job.ProgressPercentage = 100.0
			job.UploadedBytes = job.TotalBytes
			
			m.store.SaveUploadJob(job)
			m.store.RemoveFromQueue("in_progress", job.JobID)
			m.store.AddToQueue("completed", job.JobID)
			
			logger.Get().Info("Upload completed successfully", zap.String("job_id", job.JobID))
			return
		}

		// Handle Failure
		errMsg := err.Error()
		job.LastError = &errMsg
		job.ErrorLog = append(job.ErrorLog, errMsg)
		m.store.SaveUploadJob(job)
		
		logger.Get().Error("Upload failed", zap.String("job_id", job.JobID), zap.Error(err))
	}

	// Exhausted retries
	job.Status = "failed"
	m.store.SaveUploadJob(job)
	m.store.RemoveFromQueue("in_progress", job.JobID)
	m.store.AddToQueue("failed", job.JobID)
}
