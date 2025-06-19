package transactiongrpc

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

//go:generate mockery --name TransactionService --output ./mocks/
type TransactionService interface {
	Begin(ctx context.Context, description string) (string, time.Time, error)
	Commit(ctx context.Context, transactionID string) (time.Time, error)
	Rollback(ctx context.Context, transactionID string) error
}

type ServerAPI struct {
	dbv1.UnimplementedTransactionServiceServer
	txService TransactionService
}

func New(txService TransactionService) *ServerAPI {
	return &ServerAPI{
		txService: txService,
	}
}

func Register(server *grpc.Server, txService TransactionService) {
	dbv1.RegisterTransactionServiceServer(server, New(txService))
}

func (s *ServerAPI) Begin(ctx context.Context, req *dbv1.BeginRequest) (*dbv1.BeginResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	transactionID, startedAt, err := s.txService.Begin(ctx, req.GetDescription())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &dbv1.BeginResponse{
		TransactionId: transactionID,
		StartedAt:     timestamppb.New(startedAt),
	}, nil
}

func (s *ServerAPI) Commit(ctx context.Context, req *dbv1.CommitRequest) (*dbv1.CommitResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	committedAt, err := s.txService.Commit(ctx, req.GetTransactionId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &dbv1.CommitResponse{
		CommittedAt: timestamppb.New(committedAt),
	}, nil
}

func (s *ServerAPI) Rollback(ctx context.Context, req *dbv1.RollbackRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.txService.Rollback(ctx, req.GetTransactionId()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
