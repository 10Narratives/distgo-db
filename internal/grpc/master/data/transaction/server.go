package transactionapi

import (
	"context"

	mdbv1 "github.com/10Narratives/distgo-db/pkg/proto/master/database/v1"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TransactionCoordinator interface {
	Execute(context.Context, *dbv1.ExecuteRequest) (*emptypb.Empty, error)
}

type serverAPI struct {
	mdbv1.UnimplementedTransactionServiceServer
	coordinator TransactionCoordinator
}

func Register(server *grpc.Server, coordinator TransactionCoordinator) {
	mdbv1.RegisterTransactionServiceServer(server, &serverAPI{coordinator: coordinator})
}

func (s *serverAPI) Execute(ctx context.Context, req *mdbv1.ExecuteRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	operations := make([]*dbv1.Operation, 0, len(req.GetOperations()))
	for _, op := range req.GetOperations() {
		operations = append(operations, convertOperationFromGRPC(op))
	}

	return s.coordinator.Execute(ctx, &dbv1.ExecuteRequest{
		DatabaseName: req.GetDatabaseName(),
		Operations:   operations,
	})
}

func convertOperationFromGRPC(src *mdbv1.Operation) *dbv1.Operation {
	if src == nil {
		return nil
	}

	return &dbv1.Operation{
		MutationType: src.GetMutationType(),
		EntityType:   src.GetEntityType(),
		Name:         src.GetName(),
		Value:        src.GetValue(),
	}
}
