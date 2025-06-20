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

type ServerAPI struct {
	mdbv1.UnimplementedDatabaseServiceServer
	client dbv1.DatabaseServiceClient
}

func New(client dbv1.DatabaseServiceClient) *ServerAPI {
	return &ServerAPI{
		client: client,
	}
}

func Register(server *grpc.Server, client dbv1.DatabaseServiceClient) {
	mdbv1.RegisterDatabaseServiceServer(server, &ServerAPI{client: client})
}

func (s *ServerAPI) CreateDatabase(ctx context.Context, req *mdbv1.CreateDatabaseRequest) (*mdbv1.Database, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	clientReq := &dbv1.CreateDatabaseRequest{
		DatabaseId: req.DatabaseId,
		Database: &dbv1.Database{
			DisplayName: req.Database.DisplayName,
		},
	}

	resp, err := s.client.CreateDatabase(ctx, clientReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create database: %v", err)
	}

	return &mdbv1.Database{
		Name:        resp.Name,
		DisplayName: resp.DisplayName,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
	}, nil
}

func (s *ServerAPI) DeleteDatabase(ctx context.Context, req *mdbv1.DeleteDatabaseRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	clientReq := &dbv1.DeleteDatabaseRequest{
		Name: req.Name,
	}

	_, err := s.client.DeleteDatabase(ctx, clientReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to delete database: %v", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *ServerAPI) GetDatabase(ctx context.Context, req *mdbv1.GetDatabaseRequest) (*mdbv1.Database, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	clientReq := &dbv1.GetDatabaseRequest{
		Name: req.Name,
	}

	resp, err := s.client.GetDatabase(ctx, clientReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get database: %v", err)
	}

	return &mdbv1.Database{
		Name:        resp.Name,
		DisplayName: resp.DisplayName,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
	}, nil
}

func (s *ServerAPI) ListDatabases(ctx context.Context, req *mdbv1.ListDatabasesRequest) (*mdbv1.ListDatabasesResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	clientReq := &dbv1.ListDatabasesRequest{
		PageSize:  req.PageSize,
		PageToken: req.PageToken,
	}

	resp, err := s.client.ListDatabases(ctx, clientReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list databases: %v", err)
	}

	databases := make([]*mdbv1.Database, len(resp.Databases))
	for i, db := range resp.Databases {
		databases[i] = &mdbv1.Database{
			Name:        db.Name,
			DisplayName: db.DisplayName,
			CreatedAt:   db.CreatedAt,
			UpdatedAt:   db.UpdatedAt,
		}
	}

	return &mdbv1.ListDatabasesResponse{
		Databases:     databases,
		NextPageToken: resp.NextPageToken,
	}, nil
}

func (s *ServerAPI) UpdateDatabase(ctx context.Context, req *mdbv1.UpdateDatabaseRequest) (*mdbv1.Database, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	clientReq := &dbv1.UpdateDatabaseRequest{
		Database: &dbv1.Database{
			Name:        req.Database.Name,
			DisplayName: req.Database.DisplayName,
		},
		UpdateMask: req.UpdateMask,
	}

	resp, err := s.client.UpdateDatabase(ctx, clientReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update database: %v", err)
	}

	return &mdbv1.Database{
		Name:        resp.Name,
		DisplayName: resp.DisplayName,
		CreatedAt:   resp.CreatedAt,
		UpdatedAt:   resp.UpdatedAt,
	}, nil
}
