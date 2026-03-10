package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/dojo-harvester/backend/internal/logger"
)

func (s *Server) registerRecordingRoutes(api *gin.RouterGroup) {
	api.POST("/sessions/:session_id/start", s.startRecording)
	api.POST("/sessions/:session_id/stop", s.stopRecording)
}

func (s *Server) startRecording(c *gin.Context) {
	sessionID := c.Param("session_id")

	session, err := s.store.GetSession(sessionID)
	if err != nil || session == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	if session.Status != "preparing" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session is not in preparing state"})
		return
	}

	// T045: Start executor
	scriptPath := s.cfg.Scripts.Record
	args := []string{session.SessionID, session.RawStoragePath, "300"} // 5 mins limit

	// T021: Create executor service mapping
	// Need to initialize executor in Server struct, will do that in server.go
	if err := s.exec.StartScript(sessionID, scriptPath, args); err != nil {
		logger.Get().Error("Failed to start recording script", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to launch recording script"})
		return
	}

	// Update state
	session.Status = "recording"
	session.UpdatedAt = time.Now().UTC()
	
	if err := s.store.SaveSession(session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session state"})
		return
	}

	c.JSON(http.StatusOK, session)
}

func (s *Server) stopRecording(c *gin.Context) {
	sessionID := c.Param("session_id")

	session, err := s.store.GetSession(sessionID)
	if err != nil || session == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Session not found"})
		return
	}

	if session.Status != "recording" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Session is not recording"})
		return
	}

	// T046: Stop executor
	if err := s.exec.StopScript(sessionID); err != nil {
		logger.Get().Error("Failed to stop recording script", zap.Error(err))
		// We proceed anyway to update status
	}

	now := time.Now().UTC()
	
	// Update state
	session.Status = "converting" // Transition to converting phase
	session.EndTime = &now
	session.DurationSeconds = int(now.Sub(session.StartTime).Seconds())
	session.UpdatedAt = now

	if err := s.store.SaveSession(session); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update session state"})
		return
	}

	c.JSON(http.StatusOK, session)
}
