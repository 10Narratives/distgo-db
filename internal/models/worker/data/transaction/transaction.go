package transactionmodels

import commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"

type Transaction struct {
	ID         string
	Operations []commonmodels.Operation
}
