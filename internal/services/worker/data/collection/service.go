package collectionsrv

import (
	"context"
	"errors"

	collectiongrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/collection"
	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
	commonsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/common"
)

//go:generate mockery --name CollectionStorage --output ./mocks/
type CollectionStorage interface {
	Collection(ctx context.Context, key collectionmodels.Key) (collectionmodels.Collection, error)
	Collections(ctx context.Context, parentKey databasemodels.Key) []collectionmodels.Collection
	CreateCollection(ctx context.Context, key collectionmodels.Key, description string) (collectionmodels.Collection, error)
	UpdateCollection(ctx context.Context, key collectionmodels.Key, description string) error
	DeleteCollection(ctx context.Context, key collectionmodels.Key) error
}

type Service struct {
	collectionStore CollectionStorage
	walStorage      commonsrv.WALStorage
}

var _ collectiongrpc.CollectionService = &Service{}

func New(collectionStore CollectionStorage, walStorage commonsrv.WALStorage) *Service {
	return &Service{
		collectionStore: collectionStore,
		walStorage:      walStorage,
	}
}

func (s *Service) Collection(ctx context.Context, name string) (collectionmodels.Collection, error) {
	key := collectionmodels.NewKey(name)
	return s.collectionStore.Collection(ctx, key)
}

func (s *Service) Collections(ctx context.Context, parent string, size int32, token string) ([]collectionmodels.Collection, string, error) {
	parentKey := databasemodels.NewKey(parent)
	all := s.collectionStore.Collections(ctx, parentKey)

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

	coll, err := s.collectionStore.CreateCollection(ctx, key, description)
	if err != nil {
		return collectionmodels.Collection{}, err
	}

	// entry := walmodels.WALEntry{
	// 	ID:        coll.Name,
	// 	Target:    "collection",
	// 	Type:      commonmodels.MutationTypeCreate,
	// 	NewValue:  coll.Description,
	// 	Timestamp: time.Now(),
	// }
	// if err := s.walStorage.LogEntry(ctx, entry); err != nil {
	// 	return collectionmodels.Collection{}, errors.New("failed to log WAL entry: " + err.Error())
	// }

	return coll, nil
}

func (s *Service) DeleteCollection(ctx context.Context, name string) error {
	key := collectionmodels.NewKey(name)

	// coll, err := s.collectionStore.Collection(ctx, key)
	// if err != nil {
	// 	return err
	// }

	if err := s.collectionStore.DeleteCollection(ctx, key); err != nil {
		return err
	}

	// entry := walmodels.WALEntry{
	// 	ID:        coll.Name,
	// 	Target:    "collection",
	// 	Type:      commonmodels.MutationTypeDelete,
	// 	OldValue:  coll.Description,
	// 	Timestamp: time.Now(),
	// }
	// if err := s.walStorage.LogEntry(ctx, entry); err != nil {
	// 	return errors.New("failed to log WAL entry: " + err.Error())
	// }

	return nil
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

			if err := s.collectionStore.UpdateCollection(ctx, key, collection.Description); err != nil {
				return collectionmodels.Collection{}, err
			}

			// entry := walmodels.WALEntry{
			// 	ID:        collection.Name,
			// 	Target:    "collection",
			// 	Type:      commonmodels.MutationTypeUpdate,
			// 	OldValue:  existingColl.Description,
			// 	NewValue:  collection.Description,
			// 	Timestamp: time.Now(),
			// }
			// if err := s.walStorage.LogEntry(ctx, entry); err != nil {
			// 	return collectionmodels.Collection{}, errors.New("failed to log WAL entry: " + err.Error())
			// }

			existingColl.Description = collection.Description
			return existingColl, nil
		default:
			return collectionmodels.Collection{}, errors.New("unknown field: " + path)
		}
	}

	return collection, nil
}
