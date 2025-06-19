package transactionmodels

import "time"

type TransactionMetadata struct {
	ID          string    `json:"id"`
	Description string    `json:"description"`
	StartedAt   time.Time `json:"started_at"`
}
