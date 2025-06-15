package datastorage

import (
	"context"
	"fmt"
	"time"

	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
)

func (s *Storage) CreateDatabase(ctx context.Context, key databasemodels.Key, displayName string) (databasemodels.Database, error) {
	db := databasemodels.Database{
		Name:        fmt.Sprintf("databases/%s", key.Database),
		DisplayName: displayName,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	_, exists := s.databases.Load(key)
	if exists {
		return databasemodels.Database{}, ErrDatabaseAlreadyExists
	}

	s.databases.Store(key, db)
	return db, nil
}

func (s *Storage) Database(ctx context.Context, key databasemodels.Key) (databasemodels.Database, error) {
	val, ok := s.databases.Load(key)
	if !ok {
		return databasemodels.Database{}, ErrDatabaseNotFound
	}

	return val.(databasemodels.Database), nil
}

func (s *Storage) Databases(ctx context.Context) []databasemodels.Database {
	var result []databasemodels.Database

	s.databases.Range(func(_, value any) bool {
		result = append(result, value.(databasemodels.Database))
		return true
	})

	return result
}

func (s *Storage) DeleteDatabase(ctx context.Context, key databasemodels.Key) error {
	_, exists := s.databases.Load(key)
	if !exists {
		return ErrDatabaseNotFound
	}
	s.databases.Delete(key)
	return nil
}

func (s *Storage) UpdateDatabase(ctx context.Context, key databasemodels.Key, displayName string) error {
	val, ok := s.databases.Load(key)
	if !ok {
		return ErrDatabaseNotFound
	}

	db := val.(databasemodels.Database)
	db.DisplayName = displayName
	db.UpdatedAt = time.Now().UTC()

	s.databases.Store(key, db)
	return nil
}
