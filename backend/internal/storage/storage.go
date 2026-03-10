package storage

import (
	"fmt"
	"os"

	"github.com/dgraph-io/badger/v4"
	"go.uber.org/zap"
	"github.com/dojo-harvester/backend/internal/config"
	"github.com/dojo-harvester/backend/internal/logger"
)

type Store struct {
	db *badger.DB
}

func NewStore(cfg *config.StorageConfig) (*Store, error) {
	if err := os.MkdirAll(cfg.BadgerPath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create badger directory: %w", err)
	}

	opts := badger.DefaultOptions(cfg.BadgerPath)
	opts.Logger = &badgerLogger{log: logger.Get().Sugar()}

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("failed to open badger database: %w", err)
	}

	// Run value log GC periodically in the background or during close
	// depending on requirements, but for now we just return the store.

	return &Store{db: db}, nil
}

func (s *Store) Close() error {
	if s.db != nil {
		logger.Get().Info("Closing BadgerDB")
		return s.db.Close()
	}
	return nil
}

func (s *Store) DB() *badger.DB {
	return s.db
}

// badgerLogger adapts zap to badger.Logger interface
type badgerLogger struct {
	log *zap.SugaredLogger
}

func (l *badgerLogger) Errorf(f string, v ...interface{}) {
	l.log.Errorf(f, v...)
}
func (l *badgerLogger) Warningf(f string, v ...interface{}) {
	l.log.Warnf(f, v...)
}
func (l *badgerLogger) Infof(f string, v ...interface{}) {
	l.log.Infof(f, v...)
}
func (l *badgerLogger) Debugf(f string, v ...interface{}) {
	l.log.Debugf(f, v...)
}
