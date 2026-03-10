package storage

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
)

func (s *Store) SaveSession(session *RecordingSession) error {
	data, err := Encode(session)
	if err != nil {
		return fmt.Errorf("failed to encode session: %w", err)
	}

	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(SessionKey(session.SessionID), data)
	})
}

func (s *Store) GetSession(sessionID string) (*RecordingSession, error) {
	var session RecordingSession
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(SessionKey(sessionID))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return Decode(val, &session)
		})
	})

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, nil // Or custom not found error
		}
		return nil, fmt.Errorf("failed to get session: %w", err)
	}

	return &session, nil
}

func (s *Store) SavePrerequisite(prereq *PrerequisiteItem) error {
	data, err := Encode(prereq)
	if err != nil {
		return fmt.Errorf("failed to encode prerequisite: %w", err)
	}

	return s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(PrereqKey(prereq.SessionID, prereq.ItemID), data)
	})
}

func (s *Store) GetPrerequisites(sessionID string) ([]PrerequisiteItem, error) {
	var prereqs []PrerequisiteItem

	err := s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := PrereqPrefix(sessionID)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			var prereq PrerequisiteItem
			err := item.Value(func(v []byte) error {
				return Decode(v, &prereq)
			})
			if err != nil {
				return err
			}
			prereqs = append(prereqs, prereq)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get prerequisites: %w", err)
	}

	return prereqs, nil
}

func (s *Store) GetPrerequisite(sessionID, itemID string) (*PrerequisiteItem, error) {
	var prereq PrerequisiteItem
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(PrereqKey(sessionID, itemID))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return Decode(val, &prereq)
		})
	})

	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get prerequisite: %w", err)
	}

	return &prereq, nil
}
