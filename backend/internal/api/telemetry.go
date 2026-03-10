package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/dojo-harvester/backend/internal/logger"
	pb "github.com/dojo-harvester/backend/pkg/proto"
)

func (s *Server) registerTelemetryRoutes(api *gin.RouterGroup) {
	// Endpoint for the external python script to POST telemetry data
	api.POST("/telemetry", s.ingestTelemetry)
}

func (s *Server) ingestTelemetry(c *gin.Context) {
	var payload struct {
		TimestampNs    int64              `json:"timestamp_ns"`
		SessionID      string             `json:"session_id"`
		RobotID        string             `json:"robot_id"`
		JointPositions []float32          `json:"joint_positions"`
		SensorReadings map[string]float32 `json:"sensor_readings"`
		BatteryLevel   float32            `json:"battery_level"`
		ErrorFlags     []string           `json:"error_flags"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// T047: Handle disconnection or errors
	if len(payload.ErrorFlags) > 0 {
		for _, e := range payload.ErrorFlags {
			if e == "ROBOT_DISCONNECTED" {
				logger.Get().Warn("Robot disconnected, aborting session", zap.String("session_id", payload.SessionID))
				s.exec.StopScript(payload.SessionID)
				session, _ := s.store.GetSession(payload.SessionID)
				if session != nil {
					session.Status = "failed"
					s.store.SaveSession(session)
				}
				break
			}
		}
	}

	frame := &pb.RobotTelemetry{
		TimestampNs:    payload.TimestampNs,
		SessionId:      payload.SessionID,
		RobotId:        payload.RobotID,
		JointPositions: payload.JointPositions,
		SensorReadings: payload.SensorReadings,
		BatteryLevel:   payload.BatteryLevel,
		ErrorFlags:     payload.ErrorFlags,
	}

	if err := s.ws.BroadcastTelemetry(frame); err != nil {
		logger.Get().Error("Failed to broadcast telemetry", zap.Error(err))
	}

	c.Status(http.StatusAccepted)
}
