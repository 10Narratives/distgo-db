package documentmodels

import (
	"time"

	"github.com/google/uuid"
)

type TransactionStatus int

const (
	TxActive TransactionStatus = iota
	TxCommit
	TxRollback
)

type Transaction struct {
	ID         uuid.UUID
	Status     TransactionStatus
	Operations []Operation
	StartTime  time.Time
}
