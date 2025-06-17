package walmodels

import (
	"time"

	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
)

type WALEntry struct {
	ID        string
	Target    string
	Type      commonmodels.MutationType
	OldValue  string
	NewValue  string
	Timestamp time.Time
}
