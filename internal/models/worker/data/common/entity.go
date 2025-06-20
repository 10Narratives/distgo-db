package commonmodels

type EntityType int

const (
	EntityTypeUnspecified EntityType = iota
	EntityTypeDatabase
	EntityTypeCollection
	EntityTypeDocument
)
