package collectionmodels

import "github.com/10Narratives/distgo-db/internal/lib/grpc/utils"

type Key struct {
	Database   string
	Collection string
}

func NewKey(name string) Key {
	parsed := utils.ParseName(name)
	return Key{
		Database:   parsed.DatabaseID,
		Collection: parsed.CollectionID,
	}
}
