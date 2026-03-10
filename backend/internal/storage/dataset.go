package storage

import (
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
)

type Dataset struct {
	DatasetID        string                 `json:"dataset_id"`
	SessionID        string                 `json:"session_id"`
	Status           string                 `json:"status"` // pending, validating, valid, invalid, converting, converted
	StoragePath      string                 `json:"storage_path"`
	TotalFrames      int                    `json:"total_frames"`
	TotalEpisodes    int                    `json:"total_episodes"`
	DurationSeconds  int                    `json:"duration_seconds"`
	ValidationErrors []string               `json:"validation_errors"`
	InfoJSON         map[string]interface{} `json:"info_json"`
	StatsJSON        map[string]interface{} `json:"stats_json"`
	TasksJSON        map[string]interface{} `json:"tasks_json"`
	HasParquetFiles  bool                   `json:"has_parquet_files"`
	HasVideoFiles    bool                   `json:"has_video_files"`
	CreatedAt        time.Time              `json:"created_at"`
	ConvertedAt      *time.Time             `json:"converted_at,omitempty"`
}

func DatasetKey(datasetID string) []byte {
	return []byte("dataset:" + datasetID)
}

func (s *Store) SaveDataset(dataset *Dataset) error {
	data, err := Encode(dataset)
	if err != nil {
		return fmt.Errorf("failed to encode dataset: %w", err)
	}

	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(DatasetKey(dataset.DatasetID), data)
	})
}

func (s *Store) GetDataset(datasetID string) (*Dataset, error) {
	var dataset Dataset
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(DatasetKey(datasetID))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return Decode(val, &dataset)
		})
	})

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get dataset: %w", err)
	}

	return &dataset, nil
}
