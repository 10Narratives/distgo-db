package transactionsrv

import (
	"context"
	"time"

	transactiongrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/transaction"
	transactionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/transaction"
	walmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/wal"
)

//go:generate mockery --name TransactionStorage --output ./mocks/
type TransactionStorage interface {
	Transaction(ctx context.Context, transactionID string) (transactionmodels.TransactionMetadata, error)
	SaveTransaction(ctx context.Context, meta transactionmodels.TransactionMetadata) error
	DeleteTransaction(ctx context.Context, transactionID string) error
}

//go:generate mockery --name WALStorage --output ./mocks/
type WALStorage interface {
	Append(ctx context.Context, entry walmodels.WALEntry) error
	Entries(ctx context.Context, size int32, token string, from time.Time, to time.Time) ([]walmodels.WALEntry, string, error)
	Truncate(ctx context.Context, before time.Time) error
}

//go:generate mockery --name DataStorage --output ./mocks/
type DataStorage interface {
	RecoverDatabase(entry walmodels.WALEntry) error
	RecoverCollection(entry walmodels.WALEntry) error
	RecoverDocument(entry walmodels.WALEntry) error
}

type Service struct {
	txStorage   TransactionStorage
	walStorage  WALStorage
	dataStorage DataStorage
}

func New(
	txStorage TransactionStorage,
	walStorage WALStorage,
	dataStorage DataStorage,
) *Service {
	return &Service{
		txStorage:   txStorage,
		walStorage:  walStorage,
		dataStorage: dataStorage,
	}
}

var _ transactiongrpc.TransactionService = &Service{}

func (s *Service) Begin(ctx context.Context, description string) (string, time.Time, error) {
	return "", time.Time{}, nil
}

func (s *Service) Commit(ctx context.Context, transactionID string) (time.Time, error) {
	return time.Time{}, nil
}

func (s *Service) Rollback(ctx context.Context, transactionID string) error {
	return nil
}
