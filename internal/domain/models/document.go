package models

import "time"

type Document struct {
	ID        string
	Content   map[string]any
	CreatedAt time.Time
	UpdatedAt time.Time
}
