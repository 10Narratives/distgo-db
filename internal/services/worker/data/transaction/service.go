package transactionsrv

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	transactiongrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/transaction"
	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
	transactionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/transaction"
	walmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/wal"
	"github.com/google/uuid"
)

type TransactionStorage interface {
	Activate(txID string) error
	Append(ctx context.Context, entry transactionmodels.TransactionEntry) error
	GetAllByTxID(ctx context.Context, txID string) ([]transactionmodels.TransactionEntry, error)
	Remove(txID string) error
}

type DataStorage interface {
	ApplyToDatabase(entry walmodels.WALEntry) error
	ApplyToCollection(entry walmodels.WALEntry) error
	ApplyToDocument(entry walmodels.WALEntry) error
}

type Service struct {
	txStorage   TransactionStorage
	dataStorage DataStorage
}

var _ transactiongrpc.TransactionService = &Service{}

func New(txStorage TransactionStorage, dataStorage DataStorage) *Service {
	return &Service{
		txStorage:   txStorage,
		dataStorage: dataStorage,
	}
}

func (s *Service) Begin(ctx context.Context) string {
	txID := uuid.New().String()
	s.txStorage.Activate(txID)
	return txID
}

func (s *Service) Execute(ctx context.Context, txID string, operations []commonmodels.Operation) error {
	entries := make([]transactionmodels.TransactionEntry, 0, len(operations))
	for _, operation := range operations {
		entries = append(entries, transactionmodels.TransactionEntry{
			TransactionID: txID,
			Operation:     operation,
		})
	}

	for _, entry := range entries {
		err := s.txStorage.Append(ctx, entry)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) Commit(ctx context.Context, txID string) error {
	entries, err := s.getEntriesForTransaction(ctx, txID)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		walEntry, err := convertToWALEntry(entry)
		if err != nil {
			return err
		}

		switch walEntry.Entity {
		case walmodels.EntityTypeDatabase:
			if err := s.dataStorage.ApplyToDatabase(walEntry); err != nil {
				return fmt.Errorf("failed to apply to database: %w", err)
			}
		case walmodels.EntityTypeCollection:
			if err := s.dataStorage.ApplyToCollection(walEntry); err != nil {
				return fmt.Errorf("failed to apply to collection: %w", err)
			}
		case walmodels.EntityTypeDocument:
			if err := s.dataStorage.ApplyToDocument(walEntry); err != nil {
				return fmt.Errorf("failed to apply to document: %w", err)
			}
		default:
			return fmt.Errorf("unknown entity type: %v", walEntry.Entity)
		}
	}

	if err := s.txStorage.Remove(txID); err != nil {
		return fmt.Errorf("failed to remove transaction: %w", err)
	}

	return nil
}

func (s *Service) Rollback(ctx context.Context, txID string) error {
	if err := s.txStorage.Remove(txID); err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}

	return nil
}

func (s *Service) getEntriesForTransaction(ctx context.Context, txID string) ([]transactionmodels.TransactionEntry, error) {
	return s.txStorage.GetAllByTxID(ctx, txID)
}

func convertToWALEntry(entry transactionmodels.TransactionEntry) (walmodels.WALEntry, error) {
	var payload json.RawMessage
	var entity walmodels.EntityType

	switch entry.Operation.Entity {
	case commonmodels.EntityTypeDatabase:
		entity = walmodels.EntityTypeDatabase
		payload = entry.Operation.Value
	case commonmodels.EntityTypeCollection:
		entity = walmodels.EntityTypeCollection
		payload = entry.Operation.Value
	case commonmodels.EntityTypeDocument:
		entity = walmodels.EntityTypeDocument
		payload = entry.Operation.Value
	default:
		return walmodels.WALEntry{}, fmt.Errorf("unknown entity type in operation")
	}

	if len(payload) == 0 {
		return walmodels.WALEntry{}, fmt.Errorf("payload is empty for entity type %v", entity)
	}

	return walmodels.WALEntry{
		ID:        uuid.New(),
		Timestamp: time.Now(),
		Mutation:  entry.Operation.Mutation,
		Payload:   payload,
		Entity:    entity,
	}, nil
}
