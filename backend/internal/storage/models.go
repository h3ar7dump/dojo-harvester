package storage

import (
	"encoding/json"
	"time"
)

type RecordingSession struct {
	SessionID          string                 `json:"session_id"`
	TaskID             string                 `json:"task_id"`
	RobotID            string                 `json:"robot_id"`
	OperatorID         string                 `json:"operator_id"`
	Status             string                 `json:"status"` // preparing, recording, converting, uploading, completed, failed
	StartTime          time.Time              `json:"start_time"`
	EndTime            *time.Time             `json:"end_time,omitempty"`
	DurationSeconds    int                    `json:"duration_seconds"`
	RawStoragePath     string                 `json:"raw_storage_path"`
	LerobotStoragePath string                 `json:"lerobot_storage_path"`
	EpisodeID          string                 `json:"episode_id"`
	Metadata           map[string]interface{} `json:"metadata"`
	CreatedAt          time.Time              `json:"created_at"`
	UpdatedAt          time.Time              `json:"updated_at"`
}

type PrerequisiteItem struct {
	ItemID           string                 `json:"item_id"`
	SessionID        string                 `json:"session_id"`
	Name             string                 `json:"name"`
	Description      string                 `json:"description"`
	Category         string                 `json:"category"` // system, robot, camera, environment
	VerificationType string                 `json:"verification_type"` // manual, auto
	Status           string                 `json:"status"` // pending, verifying, verified, failed
	IsRequired       bool                   `json:"is_required"`
	VerificationData map[string]interface{} `json:"verification_data,omitempty"`
	VerifiedAt       *time.Time             `json:"verified_at,omitempty"`
	VerifiedBy       string                 `json:"verified_by,omitempty"`
}

// Key generators
func SessionKey(sessionID string) []byte {
	return []byte("session:" + sessionID)
}

func PrereqKey(sessionID, itemID string) []byte {
	return []byte("prereq:" + sessionID + ":" + itemID)
}

func PrereqPrefix(sessionID string) []byte {
	return []byte("prereq:" + sessionID + ":")
}

// Helpers
func Encode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func Decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}
