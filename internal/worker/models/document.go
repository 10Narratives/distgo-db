// Package models provides data structures for document storage.
package models

import "github.com/google/uuid"

// Document represents a schema-less data entity with unique ID.
type Document struct {
	ID   uuid.UUID      `json:"id"`   // Unique document identifier
	Data map[string]any `json:"data"` // Flexible key-value data storage
}

// NewDocument creates a document with auto-generated UUID and initialized data map.
// Initializes empty map if data is nil.
func NewDocument(data map[string]any) *Document { // FIXME: Replace pointer on simple struct
	if data == nil {
		data = make(map[string]any)
	}

	return &Document{
		ID:   uuid.New(),
		Data: data,
	}
}

// DeepCopy creates a fully independent duplicate of the document.
// Recursively copies all nested data structures.
func (d *Document) DeepCopy() *Document { // TODO: use existing library for making deep copy
	return &Document{
		ID:   d.ID,
		Data: deepCopyMap(d.Data),
	}
}

// deepCopyMap recursively clones a map[string]any structure
func deepCopyMap(src map[string]any) map[string]any {
	dst := make(map[string]any, len(src))
	for k, v := range src {
		dst[k] = deepCopyValue(v)
	}
	return dst
}

// deepCopyValue handles recursive copying of interface values
func deepCopyValue(src any) any {
	switch v := src.(type) {
	case map[string]any:
		return deepCopyMap(v)
	case []any:
		return deepCopySlice(v)
	case []map[string]any:
		return deepCopyMapSlice(v)
	default:
		return v // Direct value copy for primitive types
	}
}

// deepCopySlice clones slice with interface elements
func deepCopySlice(src []any) []any {
	dst := make([]any, len(src))
	for i, v := range src {
		dst[i] = deepCopyValue(v)
	}
	return dst
}

// deepCopyMapSlice clones slice of maps
func deepCopyMapSlice(src []map[string]any) []map[string]any {
	dst := make([]map[string]any, len(src))
	for i, v := range src {
		dst[i] = deepCopyMap(v)
	}
	return dst
}
