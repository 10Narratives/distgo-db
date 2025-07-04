// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: master/database/v1/database_service.proto

package mdbv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	DatabaseService_ListDatabases_FullMethodName  = "/master.database.v1.DatabaseService/ListDatabases"
	DatabaseService_GetDatabase_FullMethodName    = "/master.database.v1.DatabaseService/GetDatabase"
	DatabaseService_CreateDatabase_FullMethodName = "/master.database.v1.DatabaseService/CreateDatabase"
	DatabaseService_UpdateDatabase_FullMethodName = "/master.database.v1.DatabaseService/UpdateDatabase"
	DatabaseService_DeleteDatabase_FullMethodName = "/master.database.v1.DatabaseService/DeleteDatabase"
)

// DatabaseServiceClient is the client API for DatabaseService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DatabaseServiceClient interface {
	ListDatabases(ctx context.Context, in *ListDatabasesRequest, opts ...grpc.CallOption) (*ListDatabasesResponse, error)
	GetDatabase(ctx context.Context, in *GetDatabaseRequest, opts ...grpc.CallOption) (*Database, error)
	CreateDatabase(ctx context.Context, in *CreateDatabaseRequest, opts ...grpc.CallOption) (*Database, error)
	UpdateDatabase(ctx context.Context, in *UpdateDatabaseRequest, opts ...grpc.CallOption) (*Database, error)
	DeleteDatabase(ctx context.Context, in *DeleteDatabaseRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type databaseServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewDatabaseServiceClient(cc grpc.ClientConnInterface) DatabaseServiceClient {
	return &databaseServiceClient{cc}
}

func (c *databaseServiceClient) ListDatabases(ctx context.Context, in *ListDatabasesRequest, opts ...grpc.CallOption) (*ListDatabasesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListDatabasesResponse)
	err := c.cc.Invoke(ctx, DatabaseService_ListDatabases_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseServiceClient) GetDatabase(ctx context.Context, in *GetDatabaseRequest, opts ...grpc.CallOption) (*Database, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Database)
	err := c.cc.Invoke(ctx, DatabaseService_GetDatabase_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseServiceClient) CreateDatabase(ctx context.Context, in *CreateDatabaseRequest, opts ...grpc.CallOption) (*Database, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Database)
	err := c.cc.Invoke(ctx, DatabaseService_CreateDatabase_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseServiceClient) UpdateDatabase(ctx context.Context, in *UpdateDatabaseRequest, opts ...grpc.CallOption) (*Database, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Database)
	err := c.cc.Invoke(ctx, DatabaseService_UpdateDatabase_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *databaseServiceClient) DeleteDatabase(ctx context.Context, in *DeleteDatabaseRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, DatabaseService_DeleteDatabase_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DatabaseServiceServer is the server API for DatabaseService service.
// All implementations must embed UnimplementedDatabaseServiceServer
// for forward compatibility.
type DatabaseServiceServer interface {
	ListDatabases(context.Context, *ListDatabasesRequest) (*ListDatabasesResponse, error)
	GetDatabase(context.Context, *GetDatabaseRequest) (*Database, error)
	CreateDatabase(context.Context, *CreateDatabaseRequest) (*Database, error)
	UpdateDatabase(context.Context, *UpdateDatabaseRequest) (*Database, error)
	DeleteDatabase(context.Context, *DeleteDatabaseRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedDatabaseServiceServer()
}

// UnimplementedDatabaseServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedDatabaseServiceServer struct{}

func (UnimplementedDatabaseServiceServer) ListDatabases(context.Context, *ListDatabasesRequest) (*ListDatabasesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListDatabases not implemented")
}
func (UnimplementedDatabaseServiceServer) GetDatabase(context.Context, *GetDatabaseRequest) (*Database, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetDatabase not implemented")
}
func (UnimplementedDatabaseServiceServer) CreateDatabase(context.Context, *CreateDatabaseRequest) (*Database, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateDatabase not implemented")
}
func (UnimplementedDatabaseServiceServer) UpdateDatabase(context.Context, *UpdateDatabaseRequest) (*Database, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDatabase not implemented")
}
func (UnimplementedDatabaseServiceServer) DeleteDatabase(context.Context, *DeleteDatabaseRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteDatabase not implemented")
}
func (UnimplementedDatabaseServiceServer) mustEmbedUnimplementedDatabaseServiceServer() {}
func (UnimplementedDatabaseServiceServer) testEmbeddedByValue()                         {}

// UnsafeDatabaseServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DatabaseServiceServer will
// result in compilation errors.
type UnsafeDatabaseServiceServer interface {
	mustEmbedUnimplementedDatabaseServiceServer()
}

func RegisterDatabaseServiceServer(s grpc.ServiceRegistrar, srv DatabaseServiceServer) {
	// If the following call pancis, it indicates UnimplementedDatabaseServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&DatabaseService_ServiceDesc, srv)
}

func _DatabaseService_ListDatabases_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListDatabasesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).ListDatabases(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatabaseService_ListDatabases_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).ListDatabases(ctx, req.(*ListDatabasesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseService_GetDatabase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetDatabaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).GetDatabase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatabaseService_GetDatabase_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).GetDatabase(ctx, req.(*GetDatabaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseService_CreateDatabase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateDatabaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).CreateDatabase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatabaseService_CreateDatabase_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).CreateDatabase(ctx, req.(*CreateDatabaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseService_UpdateDatabase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateDatabaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).UpdateDatabase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatabaseService_UpdateDatabase_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).UpdateDatabase(ctx, req.(*UpdateDatabaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _DatabaseService_DeleteDatabase_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteDatabaseRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DatabaseServiceServer).DeleteDatabase(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: DatabaseService_DeleteDatabase_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DatabaseServiceServer).DeleteDatabase(ctx, req.(*DeleteDatabaseRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// DatabaseService_ServiceDesc is the grpc.ServiceDesc for DatabaseService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var DatabaseService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "master.database.v1.DatabaseService",
	HandlerType: (*DatabaseServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListDatabases",
			Handler:    _DatabaseService_ListDatabases_Handler,
		},
		{
			MethodName: "GetDatabase",
			Handler:    _DatabaseService_GetDatabase_Handler,
		},
		{
			MethodName: "CreateDatabase",
			Handler:    _DatabaseService_CreateDatabase_Handler,
		},
		{
			MethodName: "UpdateDatabase",
			Handler:    _DatabaseService_UpdateDatabase_Handler,
		},
		{
			MethodName: "DeleteDatabase",
			Handler:    _DatabaseService_DeleteDatabase_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "master/database/v1/database_service.proto",
}
