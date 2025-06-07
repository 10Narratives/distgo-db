package documentstore

import (
	"context"
	"errors"
	"time"

	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/document"
	documentsrv "github.com/10Narratives/distgo-db/internal/services/worker/document"
	"github.com/google/uuid"
)

var (
	ErrCollectionNotFound error = errors.New("collection not found")
	ErrDocumentNotFound   error = errors.New("document not found")
)

type Collection map[uuid.UUID]documentmodels.Document

type Storage struct {
	data map[string]Collection
}

var _ documentsrv.DocumentStorage = Storage{}

func New() *Storage {
	return &Storage{data: make(map[string]Collection)}
}

func NewOf(initial map[string]Collection) *Storage {
	return &Storage{data: initial}
}

func (s Storage) Get(ctx context.Context, collection string, documentID uuid.UUID) (documentmodels.Document, error) {
	coll, contains := s.data[collection]
	if !contains {
		return documentmodels.Document{}, ErrCollectionNotFound
	}

	doc, contains := coll[documentID]
	if !contains {
		return documentmodels.Document{}, ErrDocumentNotFound
	}
	return doc, nil
}

func (s Storage) Set(ctx context.Context, collection string, documentID uuid.UUID, content map[string]any) {
	col, exists := s.data[collection]
	if !exists {
		col = make(Collection)
		s.data[collection] = col
	}

	doc, exists := col[documentID]
	if exists {
		doc.Content = content
		doc.UpdatedAt = time.Now()
	} else {
		doc = documentmodels.Document{
			ID:        documentID,
			Content:   content,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
	}
	col[documentID] = doc
}

func (s Storage) List(ctx context.Context, collection string) ([]documentmodels.Document, error) {
	c, exists := s.data[collection]
	if !exists {
		return make([]documentmodels.Document, 0), ErrCollectionNotFound
	}

	listed := make([]documentmodels.Document, 0)
	for _, doc := range c {
		listed = append(listed, doc)
	}

	return listed, nil
}

func (s Storage) Delete(ctx context.Context, collection string, documentID uuid.UUID) error {
	c, exists := s.data[collection]
	if !exists {
		return ErrCollectionNotFound
	}

	_, exists = c[documentID]
	if !exists {
		return ErrDocumentNotFound
	}

	delete(s.data[collection], documentID)
	if len(s.data[collection]) == 0 {
		delete(s.data, collection)
	}

	return nil
}
