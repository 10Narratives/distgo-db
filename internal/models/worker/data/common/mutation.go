package commonmodels

type MutationType int

const (
	MutationTypeUnspecified MutationType = iota
	MutationTypeCreate
	MutationTypeUpdate
	MutationTypeDelete
	MutationTypeTransaction
)
