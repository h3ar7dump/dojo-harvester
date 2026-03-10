package api

import (
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/dojo-harvester/backend/internal/logger"
	"github.com/dojo-harvester/backend/internal/storage"
)

func (s *Server) registerDatasetRoutes(api *gin.RouterGroup) {
	api.POST("/datasets/:id/convert", s.startConversion)
	api.GET("/datasets/:id", s.getDatasetStatus)
}

func (s *Server) startConversion(c *gin.Context) {
	sessionID := c.Param("id")

	session, err := s.store.GetSession(sessionID)
	if err != nil || session == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	datasetID := uuid.New().String()
	now := time.Now().UTC()

	dataset := &storage.Dataset{
		DatasetID:   datasetID,
		SessionID:   sessionID,
		Status:      "converting",
		StoragePath: session.LerobotStoragePath,
		CreatedAt:   now,
	}

	if err := s.store.SaveDataset(dataset); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create dataset record"})
		return
	}

	scriptPath := s.cfg.Scripts.Convert
	args := []string{datasetID, session.RawStoragePath, session.LerobotStoragePath}

	if err := s.exec.StartScript(datasetID, scriptPath, args); err != nil {
		logger.Get().Error("Failed to start conversion script", zap.Error(err))
		dataset.Status = "failed"
		s.store.SaveDataset(dataset)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to launch conversion script"})
		return
	}

	// Run background monitor for conversion finish to do validation
	go s.monitorConversion(datasetID, session.LerobotStoragePath)

	c.JSON(http.StatusAccepted, dataset)
}

func (s *Server) monitorConversion(datasetID, storagePath string) {
	// Simple polling to wait for script completion
	for {
		if !s.exec.GetProcessStatus(datasetID) {
			break
		}
		time.Sleep(1 * time.Second)
	}

	dataset, err := s.store.GetDataset(datasetID)
	if err != nil || dataset == nil {
		return
	}

	// T065: Validation logic
	dataset.Status = "validating"
	s.store.SaveDataset(dataset)

	var validationErrs []string

	if _, err := os.Stat(filepath.Join(storagePath, "meta", "info.json")); os.IsNotExist(err) {
		validationErrs = append(validationErrs, "Missing meta/info.json")
	}

	if _, err := os.Stat(filepath.Join(storagePath, "meta", "stats.json")); os.IsNotExist(err) {
		validationErrs = append(validationErrs, "Missing meta/stats.json")
	}

	if _, err := os.Stat(filepath.Join(storagePath, "meta", "tasks.json")); os.IsNotExist(err) {
		validationErrs = append(validationErrs, "Missing meta/tasks.json")
	}

	// Check for at least one parquet and mp4
	parquetMatch, _ := filepath.Glob(filepath.Join(storagePath, "data", "*.parquet"))
	if len(parquetMatch) > 0 {
		dataset.HasParquetFiles = true
	} else {
		validationErrs = append(validationErrs, "No .parquet files found in data/")
	}

	videoMatch, _ := filepath.Glob(filepath.Join(storagePath, "videos", "*.mp4"))
	if len(videoMatch) > 0 {
		dataset.HasVideoFiles = true
	} else {
		validationErrs = append(validationErrs, "No .mp4 files found in videos/")
	}

	now := time.Now().UTC()
	dataset.ConvertedAt = &now
	dataset.ValidationErrors = validationErrs

	if len(validationErrs) > 0 {
		dataset.Status = "invalid"
	} else {
		dataset.Status = "converted"
	}

	s.store.SaveDataset(dataset)
}

func (s *Server) getDatasetStatus(c *gin.Context) {
	datasetID := c.Param("id")

	dataset, err := s.store.GetDataset(datasetID)
	if err != nil || dataset == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Dataset not found"})
		return
	}

	c.JSON(http.StatusOK, dataset)
}
