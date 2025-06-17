package databasemodels

import (
	"github.com/10Narratives/distgo-db/internal/lib/grpc/utils"
)

type Key struct {
	Database string
}

func NewKey(name string) Key {
	parsed := utils.ParseName(name)
	return Key{
		Database: parsed.DatabaseID,
	}
}
