package collectionsrv

import (
	"context"
	"errors"
	"time"

	collectiongrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/collection"
	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
)

//go:generate mockery --name CollectionStorage --output ./mocks/
type CollectionStorage interface {
	Collection(ctx context.Context, key collectionmodels.Key) (collectionmodels.Collection, error)
	Collections(ctx context.Context, parentKey databasemodels.Key) ([]collectionmodels.Collection, error)
	CreateCollection(ctx context.Context, key collectionmodels.Key, description string) (collectionmodels.Collection, error)
	UpdateCollection(ctx context.Context, key collectionmodels.Key, description string) error
	DeleteCollection(ctx context.Context, key collectionmodels.Key) error
}

//go:generate mockery --name WAlService --output ./mocks/
type WAlService interface {
	CreateCollectionEntry(ctx context.Context, mutation commonmodels.MutationType, key collectionmodels.Key, coll *collectionmodels.Collection) error
}

type Service struct {
	collectionStore CollectionStorage
	walService      WAlService
}

var _ collectiongrpc.CollectionService = &Service{}

func New(collectionStore CollectionStorage, walService WAlService) *Service {
	return &Service{
		collectionStore: collectionStore,
		walService:      walService,
	}
}

func (s *Service) Collection(ctx context.Context, name string) (collectionmodels.Collection, error) {
	key := collectionmodels.NewKey(name)
	return s.collectionStore.Collection(ctx, key)
}

func (s *Service) Collections(ctx context.Context, parent string, size int32, token string) ([]collectionmodels.Collection, string, error) {
	parentKey := databasemodels.NewKey(parent)
	all, err := s.collectionStore.Collections(ctx, parentKey)
	if err != nil {
		return []collectionmodels.Collection{}, "", err
	}

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
	key := collectionmodels.NewKey(parent + "/collections/" + collectionID)

	err := s.walService.CreateCollectionEntry(ctx, commonmodels.MutationTypeCreate, key, &collectionmodels.Collection{
		Name:        parent + "/collections/" + collectionID,
		Description: description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	})
	if err != nil {
		return collectionmodels.Collection{}, errors.New("failed to create collection WAL entry: " + err.Error())
	}

	coll, err := s.collectionStore.CreateCollection(ctx, key, description)
	if err != nil {
		return collectionmodels.Collection{}, err
	}

	return coll, nil
}

func (s *Service) UpdateCollection(ctx context.Context, collection collectionmodels.Collection, paths []string) (collectionmodels.Collection, error) {
	for _, path := range paths {
		switch path {
		case "description":
			key := collectionmodels.NewKey(collection.Name)

			existingColl, err := s.collectionStore.Collection(ctx, key)
			if err != nil {
				return collectionmodels.Collection{}, err
			}

			err = s.walService.CreateCollectionEntry(ctx, commonmodels.MutationTypeUpdate, key, &collectionmodels.Collection{
				Name:        existingColl.Name,
				Description: collection.Description,
				CreatedAt:   existingColl.CreatedAt,
				UpdatedAt:   time.Now(),
			})
			if err != nil {
				return collectionmodels.Collection{}, errors.New("failed to create WAL entry: " + err.Error())
			}

			if err := s.collectionStore.UpdateCollection(ctx, key, collection.Description); err != nil {
				return collectionmodels.Collection{}, err
			}

			existingColl.Description = collection.Description
			return existingColl, nil
		default:
			return collectionmodels.Collection{}, errors.New("unknown field: " + path)
		}
	}

	return collection, nil
}

func (s *Service) DeleteCollection(ctx context.Context, name string) error {
	key := collectionmodels.NewKey(name)

	err := s.walService.CreateCollectionEntry(ctx, commonmodels.MutationTypeDelete, key, nil)
	if err != nil {
		return errors.New("failed to create WAL entry: " + err.Error())
	}

	if err := s.collectionStore.DeleteCollection(ctx, key); err != nil {
		return err
	}

	return nil
}
