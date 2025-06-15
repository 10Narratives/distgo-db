package databasemodels

import "time"

type Database struct {
	Name        string    `json:"name"`
	DisplayName string    `json:"display_name"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
