package documentmodels

import (
	"time"

	"github.com/google/uuid"
)

type Document struct {
	ID         uuid.UUID `json:"id"`
	Content    string    `json:"content"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
}
