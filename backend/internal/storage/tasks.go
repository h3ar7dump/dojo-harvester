package storage

import (
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
)

type HarvestTask struct {
	TaskID            string                 `json:"task_id"`
	Title             string                 `json:"title"`
	Description       string                 `json:"description"`
	RobotConfiguration map[string]interface{} `json:"robot_configuration"`
	Objectives        []string               `json:"objectives"`
	RequiredEpisodes  int                    `json:"required_episodes"`
	CompletedEpisodes int                    `json:"completed_episodes"`
	Status            string                 `json:"status"` // pending, in_progress, completed, cancelled
	AssigneeID        *string                `json:"assignee_id,omitempty"`
	Priority          string                 `json:"priority"`
	DueDate           *time.Time             `json:"due_date,omitempty"`
	CreatedBy         string                 `json:"created_by"`
	CreatedAt         time.Time              `json:"created_at"`
	UpdatedAt         time.Time              `json:"updated_at"`
	CompletedAt       *time.Time             `json:"completed_at,omitempty"`
}

func TaskKey(taskID string) []byte {
	return []byte("task:" + taskID)
}

func (s *Store) SaveTask(task *HarvestTask) error {
	data, err := Encode(task)
	if err != nil {
		return fmt.Errorf("failed to encode task: %w", err)
	}

	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(TaskKey(task.TaskID), data)
	})
}

func (s *Store) GetTask(taskID string) (*HarvestTask, error) {
	var task HarvestTask
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(TaskKey(taskID))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return Decode(val, &task)
		})
	})

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get task: %w", err)
	}

	return &task, nil
}

func (s *Store) GetAllTasks() ([]HarvestTask, error) {
	var tasks []HarvestTask
	err := s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := []byte("task:")
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			var t HarvestTask
			err := item.Value(func(v []byte) error {
				return Decode(v, &t)
			})
			if err != nil {
				return err
			}
			tasks = append(tasks, t)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list tasks: %w", err)
	}

	return tasks, nil
}
