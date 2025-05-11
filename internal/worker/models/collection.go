package models

import (
	"errors"
	"sync"

	"github.com/google/uuid"
)

// Collection represents a thread-safe named container for documents
// Uses RWMutex for concurrent access control
type Collection struct {
	name      string
	documents map[uuid.UUID]*Document
	mu        sync.RWMutex
}

// NewCollection creates empty document collection with specified name
func NewCollection(name string) *Collection {
	return &Collection{
		name:      name,
		documents: make(map[uuid.UUID]*Document),
	}
}

// Name returns collection identifier string
func (c *Collection) Name() string {
	return c.name
}

// Insert adds document to collection. Thread-safe.
// Returns error if document ID already exists
func (c *Collection) Insert(doc *Document) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.documents[doc.ID]; exists {
		return errors.New("document already exists")
	}

	c.documents[doc.ID] = doc
	return nil
}

// FindByID retrieves document copy by UUID. Thread-safe.
// Returns error if document not found
func (c *Collection) FindByID(id uuid.UUID) (*Document, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	doc, exists := c.documents[id]
	if !exists {
		return nil, errors.New("document not found")
	}

	return doc.DeepCopy(), nil
}

// Update replaces document data with new values. Thread-safe.
// Requires non-nil data map. Returns error if document missing
func (c *Collection) Update(id uuid.UUID, newData map[string]any) error {
	if newData == nil {
		return errors.New("nil data")
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	doc, exists := c.documents[id]
	if !exists {
		return errors.New("document not found")
	}

	doc.Data = deepCopyMap(newData)
	return nil
}

// Delete removes document from collection. Thread-safe.
// Returns error if document not found
func (c *Collection) Delete(id uuid.UUID) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	if _, exists := c.documents[id]; !exists {
		return errors.New("document not found")
	}

	delete(c.documents, id)
	return nil
}

// Exists checks document presence in collection. Thread-safe
func (c *Collection) Exists(id uuid.UUID) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()

	_, exists := c.documents[id]
	return exists
}
