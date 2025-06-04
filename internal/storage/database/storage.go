package databasestore

import (
	"context"
	"errors"
	"time"

	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/database"
	databasesrv "github.com/10Narratives/distgo-db/internal/services/worker/database"
	"github.com/google/uuid"
)

var (
	ErrCollectionNotFound = errors.New("collection not found")
	ErrDocumentNotFound   = errors.New("document not found")
)

type Collection map[uuid.UUID]databasemodels.Document

type Storage struct {
	collections map[string]Collection
}

var _ databasesrv.DocumentStorage = &Storage{}

func New() *Storage {
	return &Storage{collections: make(map[string]Collection)}
}

func (s *Storage) Get(ctx context.Context, collection string, documentID uuid.UUID) (databasemodels.Document, error) {
	c, contains := s.collections[collection]
	if !contains {
		return databasemodels.Document{}, ErrCollectionNotFound
	}

	document, contains := c[documentID]
	if !contains {
		return databasemodels.Document{}, ErrDocumentNotFound
	}

	return document, nil
}

func (s *Storage) Set(ctx context.Context, collection string, documentID uuid.UUID, content map[string]any) {
	if _, ok := s.collections[collection]; !ok {
		s.collections[collection] = make(Collection)
	}

	now := time.Now()

	if existingDoc, ok := s.collections[collection][documentID]; ok {
		existingDoc.Content = content
		existingDoc.UpdatedAt = now
		s.collections[collection][documentID] = existingDoc
	} else {
		newDoc := databasemodels.Document{
			ID:        documentID,
			Content:   content,
			CreatedAt: now,
			UpdatedAt: now,
		}
		s.collections[collection][documentID] = newDoc
	}
}

func (s *Storage) Delete(ctx context.Context, collection string, documentID uuid.UUID) error {
	c, contains := s.collections[collection]
	if !contains {
		return ErrCollectionNotFound
	}

	_, contains = c[documentID]
	if !contains {
		return ErrDocumentNotFound
	}

	delete(s.collections[collection], documentID)
	if len(s.collections[collection]) == 0 {
		delete(s.collections, collection)
	}

	return nil
}

func (s *Storage) Has(ctx context.Context, collection string, documentID uuid.UUID) bool {
	_, ok := s.collections[collection][documentID]
	return ok
}

func (s *Storage) List(ctx context.Context, collection string) ([]databasemodels.Document, error) {
	c, ok := s.collections[collection]
	if !ok {
		return []databasemodels.Document{}, ErrCollectionNotFound
	}
	documents := []databasemodels.Document{}
	for _, doc := range c {
		documents = append(documents, doc)
	}
	return documents, nil
}
