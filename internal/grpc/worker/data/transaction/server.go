package transactiongrpc

import (
	"context"
	"encoding/json"

	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TransactionService interface {
	Begin(ctx context.Context) string
	Execute(ctx context.Context, txID string, operations []commonmodels.Operation) error
	Commit(ctx context.Context, txID string) error
	Rollback(ctx context.Context, txID string) error
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
	dbv1.RegisterTransactionServiceServer(server, &ServerAPI{txService: txService})
}

func (s *ServerAPI) Begin(ctx context.Context, req *dbv1.BeginRequest) (*dbv1.BeginResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	txID := s.txService.Begin(ctx)

	return &dbv1.BeginResponse{
		TransactionId: txID,
	}, nil
}

func (s *ServerAPI) Execute(ctx context.Context, req *dbv1.ExecuteRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.txService.Execute(ctx, req.GetTransactionId(), operationsFromGRPC(req.GetOperations())); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}

func (s *ServerAPI) Commit(ctx context.Context, req *dbv1.CommitRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.txService.Commit(ctx, req.GetTransactionId()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}

func (s *ServerAPI) Rollback(ctx context.Context, req *dbv1.RollbackRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.txService.Rollback(ctx, req.GetTransactionId()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return nil, nil
}

func operationsFromGRPC(ops []*dbv1.Operation) []commonmodels.Operation {
	operations := make([]commonmodels.Operation, 0, len(ops))
	for _, op := range ops {
		if !json.Valid([]byte(op.Value)) {
			continue
		}

		operations = append(operations, commonmodels.Operation{
			Mutation: mutationTypeFromGRPC(op.GetMutationType()),
			Entity:   entityTypeFromGRPC(op.GetEntityType()),
			Name:     op.GetName(),
			Value:    json.RawMessage(op.Value),
		})
	}
	return operations
}

func mutationTypeFromGRPC(typ dbv1.MutationType) commonmodels.MutationType {
	switch typ {
	case dbv1.MutationType_MUTATION_TYPE_CREATE:
		return commonmodels.MutationTypeCreate
	case dbv1.MutationType_MUTATION_TYPE_UPDATE:
		return commonmodels.MutationTypeUpdate
	case dbv1.MutationType_MUTATION_TYPE_DELETE:
		return commonmodels.MutationTypeDelete
	default:
		return commonmodels.MutationTypeUnspecified
	}
}

func entityTypeFromGRPC(typ dbv1.EntityType) commonmodels.EntityType {
	switch typ {
	case dbv1.EntityType_ENTITY_TYPE_DATABASE:
		return commonmodels.EntityTypeDatabase
	case dbv1.EntityType_ENTITY_TYPE_COLLECTION:
		return commonmodels.EntityTypeCollection
	case dbv1.EntityType_ENTITY_TYPE_DOCUMENT:
		return commonmodels.EntityTypeDocument
	default:
		return commonmodels.EntityTypeUnspecified
	}
}
