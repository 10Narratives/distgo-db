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
	Begin(context.Context, *dbv1.BeginRequest) (*dbv1.BeginResponse, error)
	Execute(context.Context, *dbv1.ExecuteRequest) (*emptypb.Empty, error)
	Commit(context.Context, *dbv1.CommitRequest) (*emptypb.Empty, error)
	Rollback(context.Context, *dbv1.RollbackRequest) (*emptypb.Empty, error)
}

type serverAPI struct {
	mdbv1.UnimplementedTransactionServiceServer
	coordinator TransactionCoordinator
}

func Register(server *grpc.Server, coordinator TransactionCoordinator) {
	mdbv1.RegisterTransactionServiceServer(server, &serverAPI{coordinator: coordinator})
}

func (s *serverAPI) Begin(ctx context.Context, req *mdbv1.BeginRequest) (*mdbv1.BeginResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := s.coordinator.Begin(ctx, &dbv1.BeginRequest{
		DatabaseName: req.GetDatabaseName(),
	})

	return &mdbv1.BeginResponse{
		TransactionId: resp.GetTransactionId(),
	}, err
}

func (s *serverAPI) Commit(ctx context.Context, req *mdbv1.CommitRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return s.coordinator.Commit(ctx, &dbv1.CommitRequest{
		DatabaseName:  req.GetDatabaseName(),
		TransactionId: req.GetTransactionId(),
	})
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
		DatabaseName:  req.GetDatabaseName(),
		TransactionId: req.GetTransactionId(),
		Operations:    operations,
	})
}

func (s *serverAPI) Rollback(ctx context.Context, req *mdbv1.RollbackRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return s.coordinator.Rollback(ctx, &dbv1.RollbackRequest{
		DatabaseName:  req.GetDatabaseName(),
		TransactionId: req.GetTransactionId(),
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
