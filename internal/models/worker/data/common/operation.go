package commonmodels

import "encoding/json"

type Operation struct {
	Mutation MutationType
	Entity   EntityType
	Name     string
	Value    json.RawMessage
}
