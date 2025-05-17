package models

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/mohae/deepcopy"
)

// Errors for Document operations in a document-oriented database.
var (
	ErrIncorrectPath      = errors.New("incorrect path: path cannot be empty or invalid")
	ErrLevelNotFound      = errors.New("path level not found: the specified field does not exist in the document")
	ErrIntermediateNotMap = errors.New("intermediate value is not a map: expected a nested document at the specified level")
	ErrUnexpectedError    = errors.New("unexpected error: failed to retrieve or update field in the document")
)

// Document represents a document in a document-oriented database.
// It stores data as nested key-value pairs and supports CRUD operations.
type Document struct {
	ID   uuid.UUID      // Unique identifier for the document.
	Data map[string]any // Nested key-value storage representing the document's structure.
}

// NewDocument creates and initializes a new Document with a unique ID.
func NewDocument() Document {
	return Document{
		ID:   uuid.New(),
		Data: make(map[string]any),
	}
}

// Copy creates a deep copy of the Document, including all nested structures.
func (d Document) Copy() Document {
	dataCopy := deepcopy.Copy(d.Data)
	return Document{d.ID, dataCopy.(map[string]any)}
}

// Get retrieves the value at the specified path within the document.
// The path uses '.' as a separator for nested fields (e.g., "user.address.city").
// Returns an error if the path is invalid or the field does not exist.
func (d Document) Get(path string) (any, error) {
	levels, err := splitPath(path)
	if err != nil {
		return nil, err
	}

	current := d.Data
	for i, level := range levels {
		value, contains := current[level]
		if !contains {
			return nil, ErrLevelNotFound
		}

		if i == len(levels)-1 {
			return value, nil
		}

		m, isMap := value.(map[string]any)
		if !isMap {
			return nil, ErrIntermediateNotMap
		}
		current = m
	}

	return nil, ErrUnexpectedError
}

// Set updates or inserts a value at the specified path within the document.
// Creates intermediate nested documents if necessary.
// The path uses '.' as a separator for nested fields (e.g., "user.address.city").
// Returns an error for invalid paths.
func (d *Document) Set(path string, value any) error {
	levels, err := splitPath(path)
	if err != nil {
		return err
	}

	current := d.Data
	for i, level := range levels {
		if i == len(levels)-1 {
			current[level] = value
			return nil
		}

		nextLevel, isMap := current[level].(map[string]any)
		if !isMap {
			nextLevel = make(map[string]any)
			current[level] = nextLevel
		}
		current = nextLevel
	}

	return ErrUnexpectedError
}

// splitPath splits a document path into levels using '.' as a separator.
// Returns an error for empty or invalid paths.
func splitPath(path string) ([]string, error) {
	if path == "" {
		return nil, ErrIncorrectPath
	}

	levels := strings.Split(path, ".")
	if len(levels) == 0 || levels[0] == "" {
		return nil, ErrIncorrectPath
	}

	return levels, nil
}
