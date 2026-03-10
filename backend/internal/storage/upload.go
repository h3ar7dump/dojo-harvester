package storage

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
)

type UploadJob struct {
	JobID              string     `json:"job_id"`
	DatasetID          string     `json:"dataset_id"`
	Status             string     `json:"status"` // pending, in_progress, paused, completed, failed
	PlatformEndpoint   string     `json:"platform_endpoint"`
	TotalBytes         int64      `json:"total_bytes"`
	UploadedBytes      int64      `json:"uploaded_bytes"`
	ProgressPercentage float64    `json:"progress_percentage"`
	RetryCount         int        `json:"retry_count"`
	LastError          *string    `json:"last_error,omitempty"`
	ErrorLog           []string   `json:"error_log"`
	StartedAt          *time.Time `json:"started_at,omitempty"`
	CompletedAt        *time.Time `json:"completed_at,omitempty"`
	ResumedAt          *time.Time `json:"resumed_at,omitempty"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

func UploadKey(jobID string) []byte {
	return []byte("upload:" + jobID)
}

func (s *Store) SaveUploadJob(job *UploadJob) error {
	data, err := Encode(job)
	if err != nil {
		return fmt.Errorf("failed to encode upload job: %w", err)
	}

	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(UploadKey(job.JobID), data)
	})
}

func (s *Store) GetUploadJob(jobID string) (*UploadJob, error) {
	var job UploadJob
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(UploadKey(jobID))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return Decode(val, &job)
		})
	})

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get upload job: %w", err)
	}

	return &job, nil
}

// UploadQueue operations
func QueueKey(queueName string) []byte {
	return []byte("queue:" + queueName)
}

func (s *Store) AddToQueue(queueName, jobID string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		var queue []string
		item, err := txn.Get(QueueKey(queueName))
		if err == nil {
			err = item.Value(func(val []byte) error {
				return json.Unmarshal(val, &queue)
			})
			if err != nil {
				return err
			}
		} else if err != badger.ErrKeyNotFound {
			return err
		}

		// Prevent duplicates
		for _, id := range queue {
			if id == jobID {
				return nil
			}
		}

		queue = append(queue, jobID)
		data, err := json.Marshal(queue)
		if err != nil {
			return err
		}
		return txn.Set(QueueKey(queueName), data)
	})
}

func (s *Store) GetQueue(queueName string) ([]string, error) {
	var queue []string
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(QueueKey(queueName))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &queue)
		})
	})

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return []string{}, nil
		}
		return nil, fmt.Errorf("failed to get queue: %w", err)
	}

	return queue, nil
}

func (s *Store) RemoveFromQueue(queueName, jobID string) error {
	return s.db.Update(func(txn *badger.Txn) error {
		var queue []string
		item, err := txn.Get(QueueKey(queueName))
		if err != nil {
			if err == badger.ErrKeyNotFound {
				return nil
			}
			return err
		}

		err = item.Value(func(val []byte) error {
			return json.Unmarshal(val, &queue)
		})
		if err != nil {
			return err
		}

		var newQueue []string
		for _, id := range queue {
			if id != jobID {
				newQueue = append(newQueue, id)
			}
		}

		data, err := json.Marshal(newQueue)
		if err != nil {
			return err
		}
		return txn.Set(QueueKey(queueName), data)
	})
}
