package collectionsrv

import (
	"context"
	"errors"

	collectiongrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/collection"
	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
)

//go:generate mockery --name CollectionStorage --output ./mocks/
type CollectionStorage interface {
	Collection(ctx context.Context, name string) (collectionmodels.Collection, error)
	Collections(ctx context.Context, parent string) []collectionmodels.Collection
	CreateCollection(ctx context.Context, name, description string) (collectionmodels.Collection, error)
	DeleteCollection(ctx context.Context, name string) error
	UpdateCollection(ctx context.Context, collection collectionmodels.Collection) error
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
	return s.collectionStore.Collection(ctx, name)
}

func (s *Service) Collections(ctx context.Context, parent string, size int32, token string) ([]collectionmodels.Collection, string, error) {
	all := s.collectionStore.Collections(ctx, parent)

	if len(all) == 0 {
		return []collectionmodels.Collection{}, "", nil
	}

	startIndex := 0
	if token != "" {
		for i, coll := range all {
			if coll.Name == token {
				startIndex = i + 1
				break
			}
		}
	}

	endIndex := startIndex + int(size)
	if endIndex > len(all) {
		endIndex = len(all)
	}

	page := all[startIndex:endIndex]

	var nextPageToken string
	if endIndex < len(all) {
		nextPageToken = all[endIndex].Name
	}

	return page, nextPageToken, nil
}

func (s *Service) CreateCollection(ctx context.Context, parent string, collectionID string, description string) (collectionmodels.Collection, error) {
	name := parent + "/collections/" + collectionID
	return s.collectionStore.CreateCollection(ctx, name, description)
}

func (s *Service) DeleteCollection(ctx context.Context, name string) error {
	return s.collectionStore.DeleteCollection(ctx, name)
}

func (s *Service) UpdateCollection(ctx context.Context, collection collectionmodels.Collection, paths []string) (collectionmodels.Collection, error) {
	for _, path := range paths {
		switch path {
		case "description":
			err := s.collectionStore.UpdateCollection(ctx, collection)
			if err != nil {
				return collectionmodels.Collection{}, err
			}
		default:
			return collectionmodels.Collection{}, errors.New("unknown field: " + path)
		}
	}
	return collection, nil
}
