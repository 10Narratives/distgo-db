package datastorage

import (
	"context"
	"time"

	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
)

func (s *Storage) Collection(ctx context.Context, key collectionmodels.Key) (collectionmodels.Collection, error) {
	val, ok := s.collections.Load(key)
	if !ok {
		return collectionmodels.Collection{}, ErrCollectionNotFound
	}
	return val.(collectionmodels.Collection), nil
}

func (s *Storage) Collections(ctx context.Context, parentKey databasemodels.Key) []collectionmodels.Collection {
	var result []collectionmodels.Collection

	s.collections.Range(func(key, value any) bool {
		k := key.(collectionmodels.Key)
		if k.Database == parentKey.Database {
			result = append(result, value.(collectionmodels.Collection))
		}
		return true
	})

	return result
}

func (s *Storage) CreateCollection(ctx context.Context, key collectionmodels.Key, description string) (collectionmodels.Collection, error) {
	_, exists := s.collections.Load(key)
	if exists {
		return collectionmodels.Collection{}, ErrCollectionNotFound
	}

	now := time.Now()
	coll := collectionmodels.Collection{
		Name:        "databases/" + key.Database + "/collections/" + key.Collection,
		Description: description,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	s.collections.Store(key, coll)
	return coll, nil
}

func (s *Storage) DeleteCollection(ctx context.Context, key collectionmodels.Key) error {
	_, exists := s.collections.Load(key)
	if !exists {
		return ErrCollectionNotFound
	}
	s.collections.Delete(key)
	return nil
}

func (s *Storage) UpdateCollection(ctx context.Context, key collectionmodels.Key, description string) error {
	val, ok := s.collections.Load(key)
	if !ok {
		return ErrCollectionNotFound
	}

	coll := val.(collectionmodels.Collection)
	coll.Description = description
	coll.UpdatedAt = time.Now()

	s.collections.Store(key, coll)
	return nil
}
