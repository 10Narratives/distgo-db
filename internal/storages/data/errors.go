package datastorage

import "errors"

var (
	ErrDatabaseNotFound        = errors.New("database not found")
	ErrDatabaseAlreadyExists   = errors.New("database already exists")
	ErrCollectionNotFound      = errors.New("collection not found")
	ErrCollectionAlreadyExists = errors.New("collection already exists")
	ErrDocumentNotFound        = errors.New("document not found")
	ErrDocumentAlreadyExists   = errors.New("document already exists")
)
