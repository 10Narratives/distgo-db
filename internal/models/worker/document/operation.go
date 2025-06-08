package documentmodels

type OperationType int

const (
	OpCreate OperationType = iota
	OpUpdate
	OpDelete
)

type Operation struct {
	Type       OperationType
	Collection string
	DocumentID string
	Content    map[string]any
}
