package datastorage

import (
	"context"
	"errors"
	"time"

	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/document"
)

var (
	ErrInvalidKey = errors.New("invalid key: missing database, collection or document ID")
)

func (s *Storage) CreateDocument(ctx context.Context, key documentmodels.Key, value string) (documentmodels.Document, error) {
	if key.Database == "" || key.Collection == "" || key.Document == "" {
		return documentmodels.Document{}, ErrInvalidKey
	}

	parentKey := collectionmodels.Key{
		Database:   key.Database,
		Collection: key.Collection,
	}
	if _, err := s.Collection(ctx, parentKey); err != nil {
		return documentmodels.Document{}, err
	}

	now := time.Now().UTC()
	doc := documentmodels.Document{
		Name:      "databases/" + key.Database + "/collections/" + key.Collection + "/documents/" + key.Document,
		Value:     []byte(value),
		CreatedAt: now,
		UpdatedAt: now,
		Parent:    "databases/" + key.Database + "/collections/" + key.Collection,
	}

	s.documents.Store(key, doc)
	return doc, nil
}

func (s *Storage) DeleteDocument(ctx context.Context, key documentmodels.Key) error {
	_, exists := s.documents.Load(key)
	if !exists {
		return ErrDocumentNotFound
	}
	s.documents.Delete(key)
	return nil
}

func (s *Storage) Document(ctx context.Context, key documentmodels.Key) (documentmodels.Document, error) {
	val, ok := s.documents.Load(key)
	if !ok {
		return documentmodels.Document{}, ErrDocumentNotFound
	}
	return val.(documentmodels.Document), nil
}

func (s *Storage) Documents(ctx context.Context, parent collectionmodels.Key) ([]documentmodels.Document, error) {
	_, err := s.Collection(ctx, parent)
	if err != nil {
		return []documentmodels.Document{}, err
	}

	var result []documentmodels.Document

	s.documents.Range(func(key, value any) bool {
		k := key.(documentmodels.Key)
		if k.Database == parent.Database && k.Collection == parent.Collection {
			result = append(result, value.(documentmodels.Document))
		}
		return true
	})

	return result, nil
}

func (s *Storage) UpdateDocument(ctx context.Context, document documentmodels.Document) error {
	key := documentmodels.NewKey(document.Name)
	existing, err := s.Document(ctx, key)
	if err != nil {
		return err
	}

	existing.Value = document.Value
	existing.UpdatedAt = time.Now().UTC()

	s.documents.Store(key, existing)
	return nil
}
