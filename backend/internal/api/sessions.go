package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"github.com/dojo-harvester/backend/internal/storage"
)

func (s *Server) registerSessionRoutes(api *gin.RouterGroup) {
	api.POST("/sessions", s.createSession)
	api.GET("/sessions/:session_id/prerequisites", s.getPrerequisites)
	api.POST("/prerequisites/:session_id/:item_id/verify", s.verifyPrerequisite)
}

type CreateSessionRequest struct {
	TaskID     string `json:"task_id" binding:"required"`
	RobotID    string `json:"robot_id" binding:"required"`
	OperatorID string `json:"operator_id" binding:"required"`
}

func (s *Server) createSession(c *gin.Context) {
	var req CreateSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionID := uuid.New().String()
	episodeID := uuid.New().String() // Generating a simple episode ID
	now := time.Now().UTC()

	session := &storage.RecordingSession{
		SessionID:          sessionID,
		TaskID:             req.TaskID,
		RobotID:            req.RobotID,
		OperatorID:         req.OperatorID,
		Status:             "preparing",
		StartTime:          now,
		RawStoragePath:     fmtPath("raw", sessionID, episodeID, now),
		LerobotStoragePath: fmtPath("lerobot", sessionID, episodeID, now),
		EpisodeID:          episodeID,
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.store.SaveSession(session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	// Create default prerequisites
	prereqs := []storage.PrerequisiteItem{
		{
			ItemID:           "sys_storage",
			SessionID:        sessionID,
			Name:             "Storage Space",
			Description:      "Check if there is enough free storage space for recording.",
			Category:         "system",
			VerificationType: "auto",
			Status:           "pending",
			IsRequired:       true,
		},
		{
			ItemID:           "rob_connect",
			SessionID:        sessionID,
			Name:             "Robot Connectivity",
			Description:      "Ensure robot is connected and responsive.",
			Category:         "robot",
			VerificationType: "auto",
			Status:           "pending",
			IsRequired:       true,
		},
		{
			ItemID:           "cam_calib",
			SessionID:        sessionID,
			Name:             "Camera Calibration",
			Description:      "Verify all cameras are calibrated and unobstructed.",
			Category:         "camera",
			VerificationType: "manual",
			Status:           "pending",
			IsRequired:       true,
		},
	}

	for _, p := range prereqs {
		if err := s.store.SavePrerequisite(&p); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to initialize prerequisites"})
			return
		}
	}

	c.JSON(http.StatusCreated, session)
}

func (s *Server) getPrerequisites(c *gin.Context) {
	sessionID := c.Param("session_id")

	prereqs, err := s.store.GetPrerequisites(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch prerequisites"})
		return
	}

	c.JSON(http.StatusOK, prereqs)
}

func (s *Server) verifyPrerequisite(c *gin.Context) {
	sessionID := c.Param("session_id")
	itemID := c.Param("item_id")

	prereq, err := s.store.GetPrerequisite(sessionID, itemID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch prerequisite"})
		return
	}
	if prereq == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Prerequisite not found"})
		return
	}

	now := time.Now().UTC()
	prereq.Status = "verified"
	prereq.VerifiedAt = &now

	// T036: Automatic verification logic simulation
	if prereq.VerificationType == "auto" {
		prereq.VerificationData = map[string]interface{}{
			"verified_automatically": true,
			"timestamp":              now,
		}
		if itemID == "sys_storage" {
			prereq.VerificationData["free_space_mb"] = 50000 // Mock value
		}
	}

	if err := s.store.SavePrerequisite(prereq); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update prerequisite"})
		return
	}

	c.JSON(http.StatusOK, prereq)
}

func fmtPath(dataType, sessionID, episodeID string, t time.Time) string {
	dateStr := t.Format("2006-01-02")
	return "/data/" + dateStr + "/" + sessionID + "/" + dataType + "/" + episodeID + "/"
}
