package documentmodels

import (
	"time"

	"github.com/google/uuid"
	"github.com/mohae/deepcopy"
)

type Document struct {
	ID        uuid.UUID      `json:"id"`
	Content   map[string]any `json:"content"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (d Document) Copy() Document {
	return Document{
		ID:        d.ID,
		Content:   copyMap(d.Content),
		CreatedAt: d.CreatedAt,
		UpdatedAt: d.UpdatedAt,
	}
}

func copyMap(originalMap map[string]any) map[string]any {
	return deepcopy.Copy(originalMap).(map[string]any)
}
