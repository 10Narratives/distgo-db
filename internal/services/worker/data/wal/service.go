package walsrv

import (
	"context"
	"encoding/json"
	"errors"
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
	// Entries retrieves WAL entries with pagination and time filtering
	// Returns:
	//   - []walmodels.WALEntry: list of WAL entries
	//   - string: next page token
	//   - error: any error that occurred
	Entries(ctx context.Context, size int32, token string, from, to time.Time) ([]walmodels.WALEntry, string, error)

	// Truncate removes WAL entries older than the specified time
	// Returns:
	//   - error: any error that occurred during truncation
	Truncate(ctx context.Context, before time.Time) error

	// Append adds a new entry to the WAL
	// Returns:
	//   - error: any error that occurred during append
	Append(ctx context.Context, entry walmodels.WALEntry) error
}

// Service implements the WAL service operations
type Service struct {
	walStorage WALStorage
}

// Verify Service implements walgrpc.WALService interface
var _ walgrpc.WALService = &Service{}

// New creates a new instance of WAL Service
func New(walStorage WALStorage) *Service {
	return &Service{
		walStorage: walStorage,
	}
}

// Entries retrieves WAL entries with pagination and optional time filtering
// Parameters:
//   - ctx: context for cancellation and timeouts
//   - size: number of entries to return (1-1000)
//   - token: pagination token (empty for first page)
//   - from: start time filter (inclusive)
//   - to: end time filter (inclusive)
//
// Returns:
//   - []walmodels.WALEntry: list of WAL entries
//   - string: next page token
//   - error: any error that occurred
func (s *Service) Entries(ctx context.Context, size int32, token string, from, to time.Time) ([]walmodels.WALEntry, string, error) {
	if size <= 0 || size > 1000 {
		return nil, "", ErrInvalidPageSize
	}

	entries, nextToken, err := s.walStorage.Entries(ctx, size, token, from, to)
	if err != nil {
		return nil, "", errors.Join(ErrStorageOperation, err)
	}

	return entries, nextToken, nil
}

// Truncate removes WAL entries older than the specified time
// Parameters:
//   - ctx: context for cancellation and timeouts
//   - before: cutoff time (exclusive)
//
// Returns:
//   - error: any error that occurred
func (s *Service) Truncate(ctx context.Context, before time.Time) error {
	if before.IsZero() {
		return ErrInvalidBeforeTime
	}

	if err := s.walStorage.Truncate(ctx, before); err != nil {
		return errors.Join(ErrStorageOperation, err)
	}
	return nil
}

// Append adds a new entry to the WAL
// Parameters:
//   - ctx: context for cancellation and timeouts
//   - entityType: type of entity being logged
//   - mutation: type of mutation operation
//   - payload: data payload to store
//
// Returns:
//   - error: any error that occurred
func (s *Service) Append(ctx context.Context, entityType walmodels.EntityType, mutation commonmodels.MutationType, payload interface{}) error {
	payloadData, err := json.Marshal(payload)
	if err != nil {
		return errors.Join(ErrPayloadMarshal, err)
	}

	entry := walmodels.WALEntry{
		ID:        uuid.New(),
		Entity:    entityType,
		Mutation:  mutation,
		Timestamp: time.Now().UTC(),
		Payload:   payloadData,
	}

	if err := s.walStorage.Append(ctx, entry); err != nil {
		return errors.Join(ErrStorageOperation, err)
	}
	return nil
}

// CreateDatabaseEntry creates a WAL entry for database operations
// Parameters:
//   - ctx: context for cancellation and timeouts
//   - mutation: type of mutation operation
//   - key: database key/identifier
//   - db: database object (nil for delete operations)
//
// Returns:
//   - error: any error that occurred
func (s *Service) CreateDatabaseEntry(ctx context.Context, mutation commonmodels.MutationType, key databasemodels.Key, db *databasemodels.Database) error {
	payload := walmodels.DatabasePayload{
		Key:      key,
		Database: db,
	}
	return s.Append(ctx, walmodels.EntityTypeDatabase, mutation, payload)
}

// CreateCollectionEntry creates a WAL entry for collection operations
// Parameters:
//   - ctx: context for cancellation and timeouts
//   - mutation: type of mutation operation
//   - key: collection key/identifier
//   - coll: collection object (nil for delete operations)
//
// Returns:
//   - error: any error that occurred
func (s *Service) CreateCollectionEntry(ctx context.Context, mutation commonmodels.MutationType, key collectionmodels.Key, coll *collectionmodels.Collection) error {
	payload := walmodels.CollectionPayload{
		Key:        key,
		Collection: coll,
	}
	return s.Append(ctx, walmodels.EntityTypeCollection, mutation, payload)
}

// CreateDocumentEntry creates a WAL entry for document operations
// Parameters:
//   - ctx: context for cancellation and timeouts
//   - mutation: type of mutation operation
//   - key: document key/identifier
//   - doc: document object (nil for delete operations)
//
// Returns:
//   - error: any error that occurred
func (s *Service) CreateDocumentEntry(ctx context.Context, mutation commonmodels.MutationType, key documentmodels.Key, doc *documentmodels.Document) error {
	payload := walmodels.DocumentPayload{
		Key:      key,
		Document: doc,
	}
	return s.Append(ctx, walmodels.EntityTypeDocument, mutation, payload)
}
