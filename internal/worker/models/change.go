package models

import (
	"time"

	"github.com/google/uuid"
)

type OperationType int

const (
	InsertOp OperationType = iota
	UpdateOp
	DeleteOp
)

type Change struct {
	Sequence  uint64        `json:"sequence"`
	Type      OperationType `json:"type"`
	Document  *Document     `json:"document"`
	ID        uuid.UUID     `json:"id"`
	Timestamp int64         `json:"timestamp"`
}

func NewChange(seq uint64, opType OperationType, doc *Document, id uuid.UUID) Change {
	return Change{
		Sequence:  seq,
		Type:      opType,
		Document:  doc,
		ID:        id,
		Timestamp: time.Now().UnixNano(),
	}
}

func (ch Change) copy() Change {
	return Change{
		Sequence:  ch.Sequence,
		Type:      ch.Type,
		Document:  ch.Document.DeepCopy(),
		ID:        ch.ID,
		Timestamp: ch.Timestamp,
	}
}
