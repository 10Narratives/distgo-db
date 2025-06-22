package utils

import (
	"strings"
)

type ParsedName struct {
	DatabaseID   string
	CollectionID string
	DocumentID   string
}

func ParseName(name string) ParsedName {
	parts := strings.Split(name, "/")
	result := ParsedName{}

	if len(parts) >= 2 && parts[0] == "databases" {
		result.DatabaseID = parts[1]
	}

	if len(parts) >= 4 && parts[2] == "collections" {
		result.CollectionID = parts[3]
	}

	if len(parts) >= 6 && parts[4] == "documents" {
		result.DocumentID = parts[5]
	}

	return result
}
