package walmodels

import (
	"time"

	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/document"
	"github.com/google/uuid"
)

type Record struct {
	Op         documentmodels.OperationType `json:"op"`
	Timestamp  time.Time                    `json:"timestamp"`
	Collection string                       `json:"collection"`
	DocumentID uuid.UUID                    `json:"document_id"`
	Content    map[string]any               `json:"content,omitempty"`
}

type Entry []byte
