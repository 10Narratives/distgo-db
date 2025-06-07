package documentsrv

import (
	"context"

	documentgrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/document"
	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/document"
	"github.com/google/uuid"
)

//go:generate mockery --name DocumentStorage --output ./mocks/
type DocumentStorage interface {
	Get(ctx context.Context, collection string, documentID uuid.UUID) (documentmodels.Document, error)
	Set(ctx context.Context, collection string, documentID uuid.UUID, content map[string]any)
	List(ctx context.Context, collection string) ([]documentmodels.Document, error)
	Delete(ctx context.Context, collection string, documentID uuid.UUID) error
	Replace(ctx context.Context, collection string, documentID uuid.UUID, content map[string]any) (documentmodels.Document, error)
}

type Service struct {
	storage DocumentStorage
}

func New(storage DocumentStorage) *Service {
	return &Service{storage: storage}
}

var _ documentgrpc.DocumentService = Service{}

func (s Service) Create(ctx context.Context, collection string, content map[string]any) (documentmodels.Document, error) {
	var documentID uuid.UUID = uuid.New()
	s.storage.Set(ctx, collection, documentID, content)
	doc, err := s.storage.Get(ctx, collection, documentID)
	return doc, err
}

func (s Service) Get(ctx context.Context, collection string, documentID string) (documentmodels.Document, error) {
	return s.storage.Get(ctx, collection, uuid.MustParse(documentID))
}

func (s Service) List(ctx context.Context, collection string) ([]documentmodels.Document, error) {
	return s.storage.List(ctx, collection)
}

func (s Service) Delete(ctx context.Context, collection string, documentID string) error {
	return s.storage.Delete(ctx, collection, uuid.MustParse(documentID))
}

func (s Service) Update(ctx context.Context, collection string, documentID string, content map[string]any) (documentmodels.Document, error) {
	return s.storage.Replace(ctx, collection, uuid.MustParse(documentID), content)
}
