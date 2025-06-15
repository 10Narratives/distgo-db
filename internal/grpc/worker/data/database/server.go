package databasegrpc

import (
	"context"

	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate mockery --name DatabaseService --output ./mocks/
type DatabaseService interface {
	CreateDatabase(ctx context.Context, databaseID, displayName string) (databasemodels.Database, error)
	DeleteDatabase(ctx context.Context, name string) error
	UpdateDatabase(ctx context.Context, database databasemodels.Database, paths []string) (databasemodels.Database, error)
	Database(ctx context.Context, name string) (databasemodels.Database, error)
	Databases(ctx context.Context, size int32, token string) ([]databasemodels.Database, string, error)
}

type ServerAPI struct {
	dbv1.UnimplementedDatabaseServiceServer
	databaseSrv DatabaseService
}

func New(databaseSrv DatabaseService) *ServerAPI {
	return &ServerAPI{
		databaseSrv: databaseSrv,
	}
}

func Register(server *grpc.Server, service DatabaseService) {
	dbv1.RegisterDatabaseServiceServer(server, New(service))
}

var _ dbv1.DatabaseServiceServer = &ServerAPI{}

func (s *ServerAPI) CreateDatabase(ctx context.Context, req *dbv1.CreateDatabaseRequest) (*dbv1.Database, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	databaseID := req.GetDatabaseId()
	database := req.GetDatabase()

	created, err := s.databaseSrv.CreateDatabase(ctx, databaseID, database.GetDisplayName())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return convertDatabaseToGRPC(created), nil
}

func (s *ServerAPI) DeleteDatabase(ctx context.Context, req *dbv1.DeleteDatabaseRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.databaseSrv.DeleteDatabase(ctx, req.GetName())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *ServerAPI) GetDatabase(ctx context.Context, req *dbv1.GetDatabaseRequest) (*dbv1.Database, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	database, err := s.databaseSrv.Database(ctx, req.GetName())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return convertDatabaseToGRPC(database), nil
}

func (s *ServerAPI) ListDatabases(ctx context.Context, req *dbv1.ListDatabasesRequest) (*dbv1.ListDatabasesResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	databases, next, err := s.databaseSrv.Databases(ctx, req.GetPageSize(), req.GetPageToken())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	listed := make([]*dbv1.Database, 0, len(databases))
	for _, db := range databases {
		listed = append(listed, convertDatabaseToGRPC(db))
	}

	return &dbv1.ListDatabasesResponse{
		Databases:     listed,
		NextPageToken: next,
	}, nil
}

func (s *ServerAPI) UpdateDatabase(ctx context.Context, req *dbv1.UpdateDatabaseRequest) (*dbv1.Database, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	database := databasemodels.Database{
		Name:        req.GetDatabase().GetName(),
		DisplayName: req.GetDatabase().GetDisplayName(),
	}

	updated, err := s.databaseSrv.UpdateDatabase(ctx, database, req.GetUpdateMask().GetPaths())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return convertDatabaseToGRPC(updated), nil
}

func convertDatabaseToGRPC(src databasemodels.Database) *dbv1.Database {
	return &dbv1.Database{
		Name:        src.Name,
		DisplayName: src.DisplayName,
		CreatedAt:   timestamppb.New(src.CreatedAt),
		UpdatedAt:   timestamppb.New(src.UpdatedAt),
	}
}
