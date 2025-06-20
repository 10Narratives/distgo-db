package databaseapi

import (
	"context"

	mdbv1 "github.com/10Narratives/distgo-db/pkg/proto/master/database/v1"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DatabaseRedirector interface {
	CreateDatabase(context.Context, *dbv1.CreateDatabaseRequest) (*dbv1.Database, error)
	DeleteDatabase(context.Context, *dbv1.DeleteDatabaseRequest) (*emptypb.Empty, error)
	GetDatabase(context.Context, *dbv1.GetDatabaseRequest) (*dbv1.Database, error)
	ListDatabases(context.Context, *dbv1.ListDatabasesRequest) (*dbv1.ListDatabasesResponse, error)
	UpdateDatabase(context.Context, *dbv1.UpdateDatabaseRequest) (*dbv1.Database, error)
}

type ServerAPI struct {
	mdbv1.UnimplementedDatabaseServiceServer
	redirector DatabaseRedirector
}

func New(redirector DatabaseRedirector) *ServerAPI {
	return &ServerAPI{
		redirector: redirector,
	}
}

func Register(server *grpc.Server, redirector DatabaseRedirector) {
	mdbv1.RegisterDatabaseServiceServer(server, New(redirector))
}

func (s *ServerAPI) CreateDatabase(ctx context.Context, req *mdbv1.CreateDatabaseRequest) (*mdbv1.Database, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := s.redirector.CreateDatabase(ctx, &dbv1.CreateDatabaseRequest{
		DatabaseId: req.DatabaseId,
		Database:   convertDatabaseFromGRPC(req.Database),
	})

	return convertDatabaseToGRPC(resp), err
}

func (s *ServerAPI) DeleteDatabase(ctx context.Context, req *mdbv1.DeleteDatabaseRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return s.redirector.DeleteDatabase(ctx, &dbv1.DeleteDatabaseRequest{
		Name: req.GetName(),
	})
}

func (s *ServerAPI) GetDatabase(ctx context.Context, req *mdbv1.GetDatabaseRequest) (*mdbv1.Database, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := s.redirector.GetDatabase(ctx, &dbv1.GetDatabaseRequest{
		Name: req.GetName(),
	})

	return convertDatabaseToGRPC(resp), err
}

func (s *ServerAPI) ListDatabases(ctx context.Context, req *mdbv1.ListDatabasesRequest) (*mdbv1.ListDatabasesResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := s.redirector.ListDatabases(ctx, &dbv1.ListDatabasesRequest{
		PageSize:  req.GetPageSize(),
		PageToken: req.GetPageToken(),
	})

	listed := make([]*mdbv1.Database, 0, len(resp.Databases))
	for _, database := range resp.Databases {
		listed = append(listed, convertDatabaseToGRPC(database))
	}

	return &mdbv1.ListDatabasesResponse{
		Databases:     listed,
		NextPageToken: resp.NextPageToken,
	}, err
}

func (s *ServerAPI) UpdateDatabase(ctx context.Context, req *mdbv1.UpdateDatabaseRequest) (*mdbv1.Database, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := s.redirector.UpdateDatabase(ctx, &dbv1.UpdateDatabaseRequest{
		Database:   convertDatabaseFromGRPC(req.GetDatabase()),
		UpdateMask: req.GetUpdateMask(),
	})

	return convertDatabaseToGRPC(resp), err
}

func convertDatabaseFromGRPC(db *mdbv1.Database) *dbv1.Database {
	if db == nil {
		return nil
	}
	return &dbv1.Database{
		Name:        db.GetName(),
		DisplayName: db.GetDisplayName(),
	}
}

func convertDatabaseToGRPC(db *dbv1.Database) *mdbv1.Database {
	if db == nil {
		return nil
	}
	return &mdbv1.Database{
		Name:        db.GetName(),
		DisplayName: db.GetDisplayName(),
	}
}
