package documentstore

import (
	"context"
	"sync"

	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/document"
	"github.com/google/uuid"
)

type Collection struct {
	documents map[uuid.UUID]documentmodels.Document
	mu        sync.RWMutex
}

func NewCollection() *Collection {
	return &Collection{
		documents: make(map[uuid.UUID]documentmodels.Document),
	}
}

func NewCollectionOf(documents map[uuid.UUID]documentmodels.Document) *Collection {
	return &Collection{
		documents: documents,
	}
}

func (c *Collection) Create(ctx context.Context, documentID uuid.UUID, content map[string]any) (documentmodels.Document, error) {
	// if err := ctx.Err(); err != nil {
	// 	return documentmodels.Document{}, err
	// }

	// c.mu.Lock()
	// defer c.mu.Unlock()

	// now := time.Now()
	// document := documentmodels.Document{
	// 	ID:         documentID,
	// 	Content:    content,
	// 	CreateTime: now,
	// 	UpdateTime: now,
	// }
	// c.documents[documentID] = document

	return documentmodels.Document{}, nil
}

func (c *Collection) Document(ctx context.Context, documentID uuid.UUID) (documentmodels.Document, error) {
	// if err := ctx.Err(); err != nil {
	// 	return documentmodels.Document{}, err
	// }

	// c.mu.RLock()
	// defer c.mu.RUnlock()

	// document, exists := c.documents[documentID]
	// if !exists {
	// 	return documentmodels.Document{}, ErrDocumentNotFound
	// }

	return documentmodels.Document{}, nil
}

func (c *Collection) Documents(ctx context.Context) ([]documentmodels.Document, error) {
	// if err := ctx.Err(); err != nil {
	// 	return nil, err
	// }

	// c.mu.RLock()
	// defer c.mu.RUnlock()

	// listed := make([]documentmodels.Document, 0, len(c.documents))
	// for _, doc := range c.documents {
	// 	listed = append(listed, doc.Copy())
	// }

	return make([]documentmodels.Document, 0), nil
}

func (c *Collection) Replace(ctx context.Context, documentID uuid.UUID, content map[string]any) (documentmodels.Document, error) {
	// if err := ctx.Err(); err != nil {
	// 	return documentmodels.Document{}, err
	// }

	// c.mu.Lock()
	// defer c.mu.Unlock()

	// now := time.Now()
	// document, exists := c.documents[documentID]
	// if !exists {
	// 	document = documentmodels.Document{
	// 		ID:         documentID,
	// 		Content:    content,
	// 		CreateTime: now,
	// 	}
	// } else {
	// 	document.Content = content
	// }
	// document.UpdateTime = now
	// c.documents[documentID] = document

	return documentmodels.Document{}, nil
}

func (c *Collection) Delete(ctx context.Context, documentID uuid.UUID) error {
	// if err := ctx.Err(); err != nil {
	// 	return err
	// }

	// c.mu.Lock()
	// defer c.mu.Unlock()

	// if _, exists := c.documents[documentID]; !exists {
	// 	return ErrDocumentNotFound
	// }
	// delete(c.documents, documentID)

	return nil
}
