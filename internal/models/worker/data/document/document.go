package documentmodels

import (
	"encoding/json"
	"time"
)

type Document struct {
	Name      string
	ID        string
	Value     json.RawMessage
	CreatedAt time.Time
	UpdatedAt time.Time
}
