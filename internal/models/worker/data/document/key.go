package documentmodels

import (
	"errors"
	"fmt"

	"github.com/10Narratives/distgo-db/internal/lib/grpc/utils"
)

var (
	ErrInvalidDocumentName = errors.New("invalid document name format")
)

type Key struct {
	Database   string
	Collection string
	Document   string
}

func NewKey(name string) Key {
	parsed := utils.ParseName(name)
	return Key{
		Database:   parsed.DatabaseID,
		Collection: parsed.CollectionID,
		Document:   parsed.DocumentID,
	}
}

func (k Key) FullName() string {
	return fmt.Sprintf("databases/%s/collections/%s/documents/%s",
		k.Database, k.Collection, k.Document)
}

func (k Key) Validate() error {
	if k.Database == "" || k.Collection == "" || k.Document == "" {
		return ErrInvalidDocumentName
	}
	return nil
}
