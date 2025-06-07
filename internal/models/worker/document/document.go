package documentmodels

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID        uuid.UUID
	Content   map[string]any
	CreatedAt time.Time
	UpdatedAt time.Time
}
