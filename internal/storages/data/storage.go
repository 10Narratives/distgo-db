package datastorage

import (
	"sync"

	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/document"
	collectionsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/collection"
	databasesrv "github.com/10Narratives/distgo-db/internal/services/worker/data/database"
)

type Storage struct {
	databases   sync.Map // databasemodels.Key -> databasemodels.Database
	collections sync.Map // collectionmodels.Key -> collectionmodels.Collection
	documents   sync.Map // documentmodels.key -> documentmodels.Document
}

func New() *Storage {
	return &Storage{}
}

// NewOf creates a new Storage pre-filled with the given initial data.
func NewOf(
	initialDatabases map[databasemodels.Key]databasemodels.Database,
	initialCollections map[collectionmodels.Key]collectionmodels.Collection,
	initialDocuments map[documentmodels.Key]documentmodels.Document,
) *Storage {
	storage := &Storage{
		databases:   sync.Map{},
		collections: sync.Map{},
		documents:   sync.Map{},
	}

	for key, db := range initialDatabases {
		storage.databases.Store(key, db)
	}

	for key, col := range initialCollections {
		storage.collections.Store(key, col)
	}

	for key, doc := range initialDocuments {
		storage.documents.Store(key, doc)
	}

	return storage
}

var (
	_ databasesrv.DatabaseStorage     = &Storage{}
	_ collectionsrv.CollectionStorage = &Storage{}
	// _ documentsrv.DocumentStorage     = &Storage{}
)
