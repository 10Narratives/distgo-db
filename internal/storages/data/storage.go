package datastorage

import (
	"sync"

	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
	databasesrv "github.com/10Narratives/distgo-db/internal/services/worker/data/database"
)

type Storage struct {
	databases sync.Map // databasemodels.Key -> databasemodels.Database
	// collections sync.Map // collectionmodels.Key -> collectionmodels.Collection
	// documents   sync.Map // documentmodels.key -> documentmodels.Document
}

func New() *Storage {
	return &Storage{}
}

func NewOf(
	initialDatabases map[databasemodels.Key]databasemodels.Database,
) *Storage {
	var databases sync.Map
	for key, database := range initialDatabases {
		databases.Store(key, database)
	}

	return &Storage{
		databases: databases,
	}
}

var (
	_ databasesrv.DatabaseStorage = &Storage{}
	// _ collectionsrv.CollectionStorage = &Storage{}
	// _ documentsrv.DocumentStorage     = &Storage{}
)
