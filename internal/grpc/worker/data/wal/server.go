package walgrpc

import (
	"context"
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
	WALEntries(ctx context.Context, size int32, token string, from, to time.Time) ([]walmodels.WALEntry, string, error)
	TruncateWAL(ctx context.Context, before time.Time) error
}

type ServerAPI struct {
	dbv1.UnimplementedWALServiceServer
	service WALService
}

var _ dbv1.WALServiceServer = ServerAPI{}

func New(service WALService) *ServerAPI {
	return &ServerAPI{
		service: service,
	}
}

func Register(server *grpc.Server, service WALService) {
	dbv1.RegisterWALServiceServer(server, New(service))
}

func (s ServerAPI) ListWALEntries(ctx context.Context, req *dbv1.ListWALEntriesRequest) (*dbv1.ListWALEntriesResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	entries, token, err := s.service.WALEntries(ctx, req.GetPageSize(), req.GetPageToken(), req.GetStartTime().AsTime(), req.GetEndTime().AsTime())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	listed := make([]*dbv1.WALEntry, 0, len(entries))
	for _, entry := range entries {
		listed = append(listed, &dbv1.WALEntry{
			Id:            entry.ID,
			Target:        entry.Target,
			OperationType: mutationTypeToGRPC(entry.Type),
			NewValue:      []byte(entry.NewValue),
			OldValue:      []byte(entry.OldValue),
			Timestamp:     timestamppb.New(entry.Timestamp),
		})
	}

	return &dbv1.ListWALEntriesResponse{
		Entries:       listed,
		NextPageToken: token,
	}, nil
}

func (s ServerAPI) TruncateWAL(ctx context.Context, req *dbv1.TruncateWALRequest) (*emptypb.Empty, error) {
	if req.GetBefore() == nil {
		return nil, status.Error(codes.InvalidArgument, "missing 'before' timestamp")
	}

	before := req.GetBefore().AsTime()
	if before.IsZero() {
		return nil, status.Error(codes.InvalidArgument, "invalid 'before' timestamp")
	}

	err := s.service.TruncateWAL(ctx, before)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
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
