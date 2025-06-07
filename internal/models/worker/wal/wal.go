package walmodels

import (
	"time"

	"github.com/google/uuid"
)

type OperationType int

const (
	OpCreate OperationType = iota
	OpUpdate
	OpDelete
)

type Record struct {
	Op         OperationType  `json:"op"`
	Timestamp  time.Time      `json:"timestamp"`
	Collection string         `json:"collection"`
	DocumentID uuid.UUID      `json:"document_id"`
	Content    map[string]any `json:"content,omitempty"`
}

type Entry []byte
