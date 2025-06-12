package workerstore

import "errors"

var (
	ErrCollectionNotFound      error = errors.New("collection not found")
	ErrCollectionAlreadyExists error = errors.New("collection already exists")
	ErrDocumentNotFound        error = errors.New("document not found")
	ErrDocumentAlreadyExists   error = errors.New("document already exists")
)
