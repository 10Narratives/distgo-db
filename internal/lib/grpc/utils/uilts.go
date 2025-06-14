package utils

import "strings"

type ParsedName struct {
	DatabaseID   string
	CollectionID string
	DocumentID   string
}

func ParseName(name string) ParsedName {
	parts := strings.Split(strings.TrimSpace(name), "/")

	db, col, doc := "", "", ""
	if len(parts) == 6 {
		db, col, doc = parts[1], parts[3], parts[5]
	} else if len(parts) == 4 {
		db, col = parts[1], parts[3]
	} else if len(parts) == 2 {
		db = parts[1]
	}

	return ParsedName{
		DatabaseID:   db,
		CollectionID: col,
		DocumentID:   doc,
	}
}
