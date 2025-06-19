package walgrpc

import (
	"context"
	"encoding/json"
	"time"

	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
	walmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/wal"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate mockery --name WALService --output ./mocks/
type WALService interface {
	Entries(ctx context.Context, size int32, token string, from, to time.Time) ([]walmodels.WALEntry, string, error)
	Truncate(ctx context.Context, before time.Time) error
}

type ServerAPI struct {
	dbv1.UnimplementedWALServiceServer
	walService WALService
}

func New(walService WALService) *ServerAPI {
	return &ServerAPI{
		walService: walService,
	}
}

func Register(server *grpc.Server, walService WALService) {
	dbv1.RegisterWALServiceServer(server, New(walService))
}

func (s *ServerAPI) ListWALEntries(ctx context.Context, req *dbv1.ListWALEntriesRequest) (*dbv1.ListWALEntriesResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	entries, token, err := s.walService.Entries(ctx, req.GetPageSize(), req.GetPageToken(), req.GetStartTime().AsTime(), req.GetEndTime().AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list wal entries: %s", err.Error())
	}

	listed := make([]*dbv1.WALEntry, 0, len(entries))
	for _, entry := range entries {
		pbEntry := &dbv1.WALEntry{
			Id:            entry.ID.String(),
			EntityType:    entityTypeToGRPC(entry.Entity),
			OperationType: mutationTypeToGRPC(entry.Mutation),
			Timestamp:     timestamppb.New(entry.Timestamp),
		}

		switch entry.Entity {
		case walmodels.EntityTypeDatabase:
			var payload walmodels.DatabasePayload
			if err := json.Unmarshal(entry.Payload, &payload); err != nil {
				return nil, status.Errorf(codes.Internal, "failed to unmarshal database payload: %s", err.Error())
			}
			pbEntry.Payload = &dbv1.WALEntry_DatabasePayload{
				DatabasePayload: &dbv1.DatabasePayload{
					DatabaseId: payload.Key.Database,
					Data:       entry.Payload,
				},
			}

		case walmodels.EntityTypeCollection:
			var payload walmodels.CollectionPayload
			if err := json.Unmarshal(entry.Payload, &payload); err != nil {
				return nil, status.Errorf(codes.Internal, "failed to unmarshal collection payload: %s", err.Error())
			}
			pbEntry.Payload = &dbv1.WALEntry_CollectionPayload{
				CollectionPayload: &dbv1.CollectionPayload{
					DatabaseId:   payload.Key.Database,
					CollectionId: payload.Key.Collection,
					Data:         entry.Payload,
				},
			}

		case walmodels.EntityTypeDocument:
			var payload walmodels.DocumentPayload
			if err := json.Unmarshal(entry.Payload, &payload); err != nil {
				return nil, status.Errorf(codes.Internal, "failed to unmarshal document payload: %s", err.Error())
			}
			pbEntry.Payload = &dbv1.WALEntry_DocumentPayload{
				DocumentPayload: &dbv1.DocumentPayload{
					DatabaseId:   payload.Key.Database,
					CollectionId: payload.Key.Collection,
					DocumentId:   payload.Key.Document,
					Data:         entry.Payload,
				},
			}
		}

		listed = append(listed, pbEntry)
	}

	return &dbv1.ListWALEntriesResponse{
		Entries:       listed,
		NextPageToken: token,
	}, nil
}

func (s *ServerAPI) TruncateWAL(ctx context.Context, req *dbv1.TruncateWALRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.walService.Truncate(ctx, req.GetBefore().AsTime())
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to truncate WAL: %s", err.Error())
	}

	return nil, nil
}

func mutationTypeToGRPC(typ commonmodels.MutationType) dbv1.MutationType {
	switch typ {
	case commonmodels.MutationTypeCreate:
		return dbv1.MutationType_MUTATION_TYPE_CREATE
	case commonmodels.MutationTypeUpdate:
		return dbv1.MutationType_MUTATION_TYPE_UPDATE
	case commonmodels.MutationTypeDelete:
		return dbv1.MutationType_MUTATION_TYPE_DELETE
	default:
		return dbv1.MutationType_MUTATION_TYPE_UNSPECIFIED
	}
}

func entityTypeToGRPC(typ walmodels.EntityType) dbv1.EntityType {
	switch typ {
	case walmodels.EntityTypeDatabase:
		return dbv1.EntityType_ENTITY_TYPE_DATABASE
	case walmodels.EntityTypeCollection:
		return dbv1.EntityType_ENTITY_TYPE_COLLECTION
	case walmodels.EntityTypeDocument:
		return dbv1.EntityType_ENTITY_TYPE_DOCUMENT
	default:
		return dbv1.EntityType_ENTITY_TYPE_UNSPECIFIED
	}
}
