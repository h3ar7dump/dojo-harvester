package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/dojo-harvester/backend/internal/storage"
)

func (s *Server) registerUploadRoutes(api *gin.RouterGroup) {
	api.POST("/datasets/:id/upload", s.startUpload)
	api.GET("/uploads/:job_id", s.getUploadStatus)
	api.GET("/uploads/queue", s.getUploadQueue)
}

func (s *Server) startUpload(c *gin.Context) {
	datasetID := c.Param("id")

	dataset, err := s.store.GetDataset(datasetID)
	if err != nil || dataset == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dataset not found"})
		return
	}

	if dataset.Status != "converted" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Dataset must be fully converted and validated before uploading"})
		return
	}

	jobID := uuid.New().String()
	now := time.Now().UTC()

	job := &storage.UploadJob{
		JobID:              jobID,
		DatasetID:          datasetID,
		Status:             "pending",
		PlatformEndpoint:   s.cfg.Platform.URL,
		TotalBytes:         1024 * 1024 * 500, // Mock 500MB
		UploadedBytes:      0,
		ProgressPercentage: 0,
		RetryCount:         0,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.store.SaveUploadJob(job); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create upload job"})
		return
	}

	if err := s.store.AddToQueue("pending", jobID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to queue upload job"})
		return
	}

	// Trigger the uploader service
	if err := s.uploader.StartJob(jobID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to start upload process"})
		return
	}

	c.JSON(http.StatusAccepted, job)
}

func (s *Server) getUploadStatus(c *gin.Context) {
	jobID := c.Param("job_id")

	job, err := s.store.GetUploadJob(jobID)
	if err != nil || job == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Upload job not found"})
		return
	}

	c.JSON(http.StatusOK, job)
}

func (s *Server) getUploadQueue(c *gin.Context) {
	pending, _ := s.store.GetQueue("pending")
	inProgress, _ := s.store.GetQueue("in_progress")
	completed, _ := s.store.GetQueue("completed")
	failed, _ := s.store.GetQueue("failed")

	c.JSON(http.StatusOK, gin.H{
		"pending":     pending,
		"in_progress": inProgress,
		"completed":   completed,
		"failed":      failed,
	})
}
