package commonmodels

import "encoding/json"

type Operation struct {
	Mutation MutationType    `json:"mutation_type"`
	Entity   EntityType      `json:"entity_type"`
	Name     string          `json:"name"`
	Value    json.RawMessage `json:"value"`
}
