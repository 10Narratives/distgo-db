package datastorage

import (
	"context"

	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
)

func (s *Storage) CreateDatabase(ctx context.Context, name string, displayName string) (databasemodels.Database, error) {
	panic("unimplemented")
}

func (s *Storage) Database(ctx context.Context, name string) (databasemodels.Database, error) {
	panic("unimplemented")
}

func (s *Storage) Databases(ctx context.Context) []databasemodels.Database {
	panic("unimplemented")
}

func (s *Storage) DeleteDatabase(ctx context.Context, name string) error {
	panic("unimplemented")
}

func (s *Storage) UpdateDatabase(ctx context.Context, name string, displayName string) error {
	panic("unimplemented")
}
