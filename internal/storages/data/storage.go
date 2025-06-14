package datastorage

import (
	"sync"

	collectionsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/collection"
	databasesrv "github.com/10Narratives/distgo-db/internal/services/worker/data/database"
	documentsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/document"
)

type Storage struct {
	databases   sync.Map
	collections sync.Map
	documents   sync.Map
}

func New() *Storage {
	return &Storage{}
}

var (
	_ databasesrv.DatabaseStorage     = &Storage{}
	_ collectionsrv.CollectionStorage = &Storage{}
	_ documentsrv.DocumentStorage     = &Storage{}
)
