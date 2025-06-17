package walstorage

import (
	"context"
	"time"

	walmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/wal"
	walsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/wal"
)

type Storage struct{}

var _ walsrv.WALStorage = &Storage{}

func (s *Storage) Entries(ctx context.Context, size int32, token string, from time.Time, to time.Time) ([]walmodels.WALEntry, string, error) {
	panic("unimplemented")
}

func (s *Storage) Truncate(ctx context.Context, before time.Time) error {
	panic("unimplemented")
}
