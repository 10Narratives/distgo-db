package documentstore

import (
	"context"

	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/document"
	documentsrv "github.com/10Narratives/distgo-db/internal/services/worker/document"
	"github.com/google/uuid"
)

type Storage struct {
	collections map[string]*Collection
}

var _ documentsrv.DocumentStorage = Storage{}

func NewStorage() *Storage {
	return &Storage{collections: make(map[string]*Collection)}
}

func NewStorageOf(collections map[string]*Collection) *Storage {
	return &Storage{collections: collections}
}

func (s Storage) Delete(ctx context.Context, collection string, documentID uuid.UUID) error {
	c, exists := s.collections[collection]
	if !exists {
		return ErrCollectionNotFound
	}

	return c.Delete(ctx, documentID)
}

func (s Storage) Get(ctx context.Context, collection string, documentID uuid.UUID) (documentmodels.Document, error) {
	c, exists := s.collections[collection]
	if !exists {
		return documentmodels.Document{}, ErrCollectionNotFound
	}

	return c.Document(ctx, documentID)
}

func (s Storage) List(ctx context.Context, collection string) ([]documentmodels.Document, error) {
	c, exists := s.collections[collection]
	if !exists {
		return []documentmodels.Document{}, ErrCollectionNotFound
	}

	return c.Documents(ctx)
}

func (s Storage) Replace(ctx context.Context, collection string, documentID uuid.UUID, content map[string]any) (documentmodels.Document, error) {
	c, exists := s.collections[collection]
	if !exists {
		return documentmodels.Document{}, ErrCollectionNotFound
	}

	return c.Replace(ctx, documentID, content)
}

func (s Storage) Set(ctx context.Context, collection string, documentID uuid.UUID, content map[string]any) {
	c, exists := s.collections[collection]
	if !exists {
		c = NewCollection()
	}

	c.Replace(ctx, documentID, content)
}
