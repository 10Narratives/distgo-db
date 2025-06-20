package clustermodels

import (
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type Worker struct {
	ID       uuid.UUID
	Database string `json:"database"`
	Address  string `json:"address"`
	Conn     *grpc.ClientConn
}
