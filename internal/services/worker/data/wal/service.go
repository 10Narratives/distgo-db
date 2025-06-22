package walsrv

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	walgrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/wal"
	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/document"
	walmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/wal"
	"github.com/google/uuid"
)

// Error definitions
var (
	ErrInvalidPageSize   = errors.New("page size must be between 1 and 1000")
	ErrInvalidBeforeTime = errors.New("before time cannot be zero")
	ErrPayloadMarshal    = errors.New("failed to marshal payload")
	ErrStorageOperation  = errors.New("storage operation failed")
)

//go:generate mockery --name WALStorage --output ./mocks/
type WALStorage interface {
	Entries(ctx context.Context, size int32, token string, from, to time.Time) ([]walmodels.WALEntry, string, error)
	Truncate(ctx context.Context, before time.Time) error
	Append(ctx context.Context, entry walmodels.WALEntry) error
}

// Service implements the WAL service operations
type Service struct {
	walStorage WALStorage
}

var _ walgrpc.WALService = &Service{}

func New(walStorage WALStorage) *Service {
	return &Service{
		walStorage: walStorage,
	}
}

func (s *Service) Entries(ctx context.Context, size int32, token string, from, to time.Time) ([]walmodels.WALEntry, string, error) {
	if size <= 0 || size > 1000 {
		return nil, "", ErrInvalidPageSize
	}

	entries, nextToken, err := s.walStorage.Entries(ctx, size, token, from, to)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %v", ErrStorageOperation, err)
	}

	return entries, nextToken, nil
}

func (s *Service) Truncate(ctx context.Context, before time.Time) error {
	if before.IsZero() {
		return ErrInvalidBeforeTime
	}

	if err := s.walStorage.Truncate(ctx, before); err != nil {
		return fmt.Errorf("%w: %v", ErrStorageOperation, err)
	}
	return nil
}

func (s *Service) Append(ctx context.Context, entityType commonmodels.EntityType, mutation commonmodels.MutationType, payload interface{}) error {
	if payload == nil {
		entry := walmodels.WALEntry{
			ID:        uuid.New(),
			Entity:    entityType,
			Mutation:  mutation,
			Timestamp: time.Now().UTC(),
			Payload:   json.RawMessage("null"),
		}
		return s.walStorage.Append(ctx, entry)
	}

	payloadData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("%w: %v", ErrPayloadMarshal, err)
	}

	entry := walmodels.WALEntry{
		ID:        uuid.New(),
		Entity:    entityType,
		Mutation:  mutation,
		Timestamp: time.Now().UTC(),
		Payload:   payloadData,
	}

	return s.walStorage.Append(ctx, entry)
}

func (s *Service) LogTransaction(ctx context.Context, operations []commonmodels.Operation) error {
	payload := walmodels.TransactionPayload{
		Operations: operations,
	}

	return s.Append(ctx, commonmodels.EntityTypeTransaction, commonmodels.MutationTypeTransaction, payload)
}

func (s *Service) CreateDatabaseEntry(ctx context.Context, mutation commonmodels.MutationType, key databasemodels.Key, db *databasemodels.Database) error {
	payload := walmodels.DatabasePayload{
		Key:      key,
		Database: db,
	}
	return s.Append(ctx, commonmodels.EntityTypeDatabase, mutation, payload)
}

func (s *Service) CreateCollectionEntry(ctx context.Context, mutation commonmodels.MutationType, key collectionmodels.Key, coll *collectionmodels.Collection) error {
	payload := walmodels.CollectionPayload{
		Key:        key,
		Collection: coll,
	}
	return s.Append(ctx, commonmodels.EntityTypeCollection, mutation, payload)
}

func (s *Service) CreateDocumentEntry(ctx context.Context, mutation commonmodels.MutationType, key documentmodels.Key, doc *documentmodels.Document) error {
	payload := walmodels.DocumentPayload{
		Key:      key,
		Document: doc,
	}
	return s.Append(ctx, commonmodels.EntityTypeDocument, mutation, payload)
}
