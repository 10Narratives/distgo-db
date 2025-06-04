package databasesrv

import (
	"context"
	"log/slog"

	databasegrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/database"
	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/database"
	"github.com/google/uuid"
)

type DocumentStorage interface {
	List(ctx context.Context, collection string) []databasemodels.Document
	Get(ctx context.Context, collection string, documentID uuid.UUID) (databasemodels.Document, error)
	Set(ctx context.Context, collection string, documentID uuid.UUID, content map[string]any)
	Delete(ctx context.Context, collection string, documentID uuid.UUID) error
	Has(ctx context.Context, collection string, documentID uuid.UUID) bool
}

type Service struct {
	log     *slog.Logger
	storage DocumentStorage
}

var _ databasegrpc.DatabaseService = Service{}

func New(log *slog.Logger, storage DocumentStorage) *Service {
	return &Service{
		log:     log,
		storage: storage,
	}
}

func (s Service) CreateDocument(ctx context.Context, collection string, content map[string]any) (databasemodels.Document, error) {
	id := uuid.New()
	s.storage.Set(ctx, collection, id, content)
	return s.storage.Get(ctx, collection, id)
}

func (s Service) DeleteDocument(ctx context.Context, collection string, documentID string) error {
	return s.storage.Delete(ctx, collection, uuid.MustParse(documentID))
}

func (s Service) Document(ctx context.Context, collection string, documentID string) (databasemodels.Document, error) {
	return s.storage.Get(ctx, collection, uuid.MustParse(documentID))
}

func (s Service) Documents(ctx context.Context, collection string) []databasemodels.Document {
	return s.storage.List(ctx, collection)
}

func (s Service) UpdateDocument(ctx context.Context, collection string, documentID string) (databasemodels.Document, error) {
	panic("unimplemented")
}
