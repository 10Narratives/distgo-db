package commonsrv

import (
	"context"

	walmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/wal"
)

//go:generate mockery --name WALStorage --output ./mocks/
type WALStorage interface {
	LogEntry(ctx context.Context, entry walmodels.WALEntry) error
}
