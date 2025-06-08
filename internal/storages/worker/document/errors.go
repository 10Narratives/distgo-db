package documentstore

import "errors"

var (
	ErrCollectionNotFound error = errors.New("collection not found")
	ErrDocumentNotFound   error = errors.New("document not found")
)
