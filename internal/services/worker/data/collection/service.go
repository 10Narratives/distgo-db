package collectionsrv

import (
	"context"

	collectiongrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/collection"
	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
)

//go:generate mockery --name CollectionStorage --output ./mocks/
type CollectionStorage interface {
	Collection(ctx context.Context, name string) (collectionmodels.Collection, error)
	Collections(ctx context.Context, parent string) ([]collectionmodels.Collection, error)
	CreateCollection(ctx context.Context, name string) (collectionmodels.Collection, error)
	DeleteCollection(ctx context.Context, name string) error
	UpdateCollection(ctx context.Context, name string) (collectionmodels.Collection, error)
}

type Service struct {
	collectionStore CollectionStorage
}

var _ collectiongrpc.CollectionService = &Service{}

func New(collectionStore CollectionStorage) *Service {
	return &Service{
		collectionStore: collectionStore,
	}
}

func (s *Service) Collection(ctx context.Context, name string) (collectionmodels.Collection, error) {
	panic("unimplemented")
}

func (s *Service) Collections(ctx context.Context, parent string, size int32, token string) ([]collectionmodels.Collection, string, error) {
	panic("unimplemented")
}

// CreateCollection implements collectiongrpc.CollectionService.
func (s *Service) CreateCollection(ctx context.Context, parent string, collectionID string) (collectionmodels.Collection, error) {
	panic("unimplemented")
}

// DeleteCollection implements collectiongrpc.CollectionService.
func (s *Service) DeleteCollection(ctx context.Context, name string) error {
	panic("unimplemented")
}

// UpdateCollection implements collectiongrpc.CollectionService.
func (s *Service) UpdateCollection(ctx context.Context, collection collectionmodels.Collection, paths []string) (collectionmodels.Collection, error) {
	panic("unimplemented")
}
