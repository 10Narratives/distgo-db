package transactionmodels

import commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"

type TransactionEntry struct {
	TransactionID string
	Operation     commonmodels.Operation
}
