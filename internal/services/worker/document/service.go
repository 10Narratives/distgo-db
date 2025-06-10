package documentsrv

import (
	"context"

	documentgrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/document"
	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/document"
	walmodels "github.com/10Narratives/distgo-db/internal/models/worker/wal"
	"github.com/google/uuid"
)

//go:generate mockery --name DocumentStorage --output ./mocks/
type DocumentStorage interface {
	Get(ctx context.Context, collection string, documentID uuid.UUID) (documentmodels.Document, error)
	Set(ctx context.Context, collection string, documentID uuid.UUID, content map[string]any)
	List(ctx context.Context, collection string) ([]documentmodels.Document, error)
	Delete(ctx context.Context, collection string, documentID uuid.UUID) error
}

//go:generate mockery --name WALStorage --output ./mocks/
type WALStorage interface {
	Write(entry walmodels.Entry) error
	Replay(handler func(walmodels.Entry) error) error
}

type Service struct {
	documentStorage DocumentStorage
	walStorage      WALStorage
}

func New(documentStorage DocumentStorage, walStorage WALStorage) *Service {
	service := &Service{
		documentStorage: documentStorage,
		walStorage:      walStorage,
	}

	return service
}

var _ documentgrpc.DocumentService = Service{}

func (s Service) Collection(ctx context.Context, collectionID string) (documentmodels.Collection, error) {
	panic("unimplemented")
}

func (s Service) Collections(ctx context.Context, size int, token string) ([]documentmodels.Collection, string, int) {
	panic("unimplemented")
}

func (s Service) CreateCollection(ctx context.Context, collectionID string) (documentmodels.Collection, error) {
	panic("unimplemented")
}

func (s Service) CreateDocument(ctx context.Context, collectionID string, content string) (documentmodels.Document, error) {
	panic("unimplemented")
}

func (s Service) DeleteDocument(ctx context.Context, collectionID string, documentID uuid.UUID) error {
	panic("unimplemented")
}

func (s Service) Document(ctx context.Context, collectionID string, documentID uuid.UUID) (documentmodels.Document, error) {
	panic("unimplemented")
}

func (s Service) Documents(ctx context.Context, collectionID string, size int, token string) ([]documentmodels.Document, string, int, error) {
	panic("unimplemented")
}

func (s Service) UpdateDocument(ctx context.Context, collectionID string, documentID uuid.UUID, changes string) (documentmodels.Document, error) {
	panic("unimplemented")
}
