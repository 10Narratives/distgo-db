package transactiongrpc

import (
	"context"
	"encoding/json"
	"errors"

	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type TransactionService interface {
	Execute(ctx context.Context, operations []commonmodels.Operation) error
}

type serverAPI struct {
	dbv1.UnimplementedTransactionServiceServer
	transactionService TransactionService
}

func Register(server *grpc.Server, transactionService TransactionService) {
	dbv1.RegisterTransactionServiceServer(server, &serverAPI{
		transactionService: transactionService,
	})
}

func (s *serverAPI) Execute(ctx context.Context, req *dbv1.ExecuteRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	operations, err := operationsFromGRPC(req.GetOperations())
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.transactionService.Execute(ctx, operations); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func operationsFromGRPC(source []*dbv1.Operation) ([]commonmodels.Operation, error) {
	operations := make([]commonmodels.Operation, 0, len(source))
	for _, op := range source {
		mutationType := mutationTypeFromGRPC(op.GetMutationType())
		if mutationType == commonmodels.MutationTypeUnspecified {
			return nil, errors.New("contains operation with unspecified mutation type")
		}

		entityType := entityTypeFromGRPC(op.GetEntityType())
		if entityType == commonmodels.EntityTypeUnspecified {
			return nil, errors.New("contains operation with unspecified entity type")
		}

		operations = append(operations, commonmodels.Operation{
			Mutation: mutationType,
			Entity:   entityType,
			Name:     op.GetName(),
			Value:    json.RawMessage(op.GetValue()),
		})
	}
	return operations, nil
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
	case dbv1.EntityType_ENTITY_TYPE_DOCUMENT:
		return commonmodels.EntityTypeDocument
	case dbv1.EntityType_ENTITY_TYPE_COLLECTION:
		return commonmodels.EntityTypeCollection
	default:
		return commonmodels.EntityTypeUnspecified
	}
}
