package walmodels

import (
	"encoding/json"
	"time"

	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/document"
	"github.com/google/uuid"
)

type WALEntry struct {
	ID        uuid.UUID                 `json:"id"`
	Timestamp time.Time                 `json:"timestamp"`
	Mutation  commonmodels.MutationType `json:"mutation"`
	Payload   json.RawMessage           `json:"payload"`
	Entity    commonmodels.EntityType   `json:"entity"`
}

type DatabasePayload struct {
	Key      databasemodels.Key       `json:"key"`
	Database *databasemodels.Database `json:"database,omitempty"`
}

type CollectionPayload struct {
	Key        collectionmodels.Key         `json:"key"`
	Collection *collectionmodels.Collection `json:"collection,omitempty"`
}

type DocumentPayload struct {
	Key      documentmodels.Key       `json:"key"`
	Document *documentmodels.Document `json:"document,omitempty"`
}

type TransactionPayload struct {
	Operations []commonmodels.Operation `json:"operations"`
}
