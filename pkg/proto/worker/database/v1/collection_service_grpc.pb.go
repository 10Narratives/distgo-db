// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.29.3
// source: worker/database/v1/collection_service.proto

package dbv1

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
	CollectionService_ListCollections_FullMethodName  = "/worker.database.v1.CollectionService/ListCollections"
	CollectionService_GetCollection_FullMethodName    = "/worker.database.v1.CollectionService/GetCollection"
	CollectionService_CreateCollection_FullMethodName = "/worker.database.v1.CollectionService/CreateCollection"
	CollectionService_UpdateCollection_FullMethodName = "/worker.database.v1.CollectionService/UpdateCollection"
	CollectionService_DeleteCollection_FullMethodName = "/worker.database.v1.CollectionService/DeleteCollection"
)

// CollectionServiceClient is the client API for CollectionService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CollectionServiceClient interface {
	ListCollections(ctx context.Context, in *ListCollectionsRequest, opts ...grpc.CallOption) (*ListCollectionsResponse, error)
	GetCollection(ctx context.Context, in *GetCollectionRequest, opts ...grpc.CallOption) (*Collection, error)
	CreateCollection(ctx context.Context, in *CreateCollectionRequest, opts ...grpc.CallOption) (*Collection, error)
	UpdateCollection(ctx context.Context, in *UpdateCollectionRequest, opts ...grpc.CallOption) (*Collection, error)
	DeleteCollection(ctx context.Context, in *DeleteCollectionRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
}

type collectionServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCollectionServiceClient(cc grpc.ClientConnInterface) CollectionServiceClient {
	return &collectionServiceClient{cc}
}

func (c *collectionServiceClient) ListCollections(ctx context.Context, in *ListCollectionsRequest, opts ...grpc.CallOption) (*ListCollectionsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListCollectionsResponse)
	err := c.cc.Invoke(ctx, CollectionService_ListCollections_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *collectionServiceClient) GetCollection(ctx context.Context, in *GetCollectionRequest, opts ...grpc.CallOption) (*Collection, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Collection)
	err := c.cc.Invoke(ctx, CollectionService_GetCollection_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *collectionServiceClient) CreateCollection(ctx context.Context, in *CreateCollectionRequest, opts ...grpc.CallOption) (*Collection, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Collection)
	err := c.cc.Invoke(ctx, CollectionService_CreateCollection_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *collectionServiceClient) UpdateCollection(ctx context.Context, in *UpdateCollectionRequest, opts ...grpc.CallOption) (*Collection, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(Collection)
	err := c.cc.Invoke(ctx, CollectionService_UpdateCollection_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *collectionServiceClient) DeleteCollection(ctx context.Context, in *DeleteCollectionRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, CollectionService_DeleteCollection_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CollectionServiceServer is the server API for CollectionService service.
// All implementations must embed UnimplementedCollectionServiceServer
// for forward compatibility.
type CollectionServiceServer interface {
	ListCollections(context.Context, *ListCollectionsRequest) (*ListCollectionsResponse, error)
	GetCollection(context.Context, *GetCollectionRequest) (*Collection, error)
	CreateCollection(context.Context, *CreateCollectionRequest) (*Collection, error)
	UpdateCollection(context.Context, *UpdateCollectionRequest) (*Collection, error)
	DeleteCollection(context.Context, *DeleteCollectionRequest) (*emptypb.Empty, error)
	mustEmbedUnimplementedCollectionServiceServer()
}

// UnimplementedCollectionServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCollectionServiceServer struct{}

func (UnimplementedCollectionServiceServer) ListCollections(context.Context, *ListCollectionsRequest) (*ListCollectionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCollections not implemented")
}
func (UnimplementedCollectionServiceServer) GetCollection(context.Context, *GetCollectionRequest) (*Collection, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCollection not implemented")
}
func (UnimplementedCollectionServiceServer) CreateCollection(context.Context, *CreateCollectionRequest) (*Collection, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCollection not implemented")
}
func (UnimplementedCollectionServiceServer) UpdateCollection(context.Context, *UpdateCollectionRequest) (*Collection, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateCollection not implemented")
}
func (UnimplementedCollectionServiceServer) DeleteCollection(context.Context, *DeleteCollectionRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteCollection not implemented")
}
func (UnimplementedCollectionServiceServer) mustEmbedUnimplementedCollectionServiceServer() {}
func (UnimplementedCollectionServiceServer) testEmbeddedByValue()                           {}

// UnsafeCollectionServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CollectionServiceServer will
// result in compilation errors.
type UnsafeCollectionServiceServer interface {
	mustEmbedUnimplementedCollectionServiceServer()
}

func RegisterCollectionServiceServer(s grpc.ServiceRegistrar, srv CollectionServiceServer) {
	// If the following call pancis, it indicates UnimplementedCollectionServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&CollectionService_ServiceDesc, srv)
}

func _CollectionService_ListCollections_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCollectionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CollectionServiceServer).ListCollections(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CollectionService_ListCollections_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CollectionServiceServer).ListCollections(ctx, req.(*ListCollectionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CollectionService_GetCollection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCollectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CollectionServiceServer).GetCollection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CollectionService_GetCollection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CollectionServiceServer).GetCollection(ctx, req.(*GetCollectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CollectionService_CreateCollection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCollectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CollectionServiceServer).CreateCollection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CollectionService_CreateCollection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CollectionServiceServer).CreateCollection(ctx, req.(*CreateCollectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CollectionService_UpdateCollection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateCollectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CollectionServiceServer).UpdateCollection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CollectionService_UpdateCollection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CollectionServiceServer).UpdateCollection(ctx, req.(*UpdateCollectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CollectionService_DeleteCollection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteCollectionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CollectionServiceServer).DeleteCollection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: CollectionService_DeleteCollection_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CollectionServiceServer).DeleteCollection(ctx, req.(*DeleteCollectionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CollectionService_ServiceDesc is the grpc.ServiceDesc for CollectionService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CollectionService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "worker.database.v1.CollectionService",
	HandlerType: (*CollectionServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ListCollections",
			Handler:    _CollectionService_ListCollections_Handler,
		},
		{
			MethodName: "GetCollection",
			Handler:    _CollectionService_GetCollection_Handler,
		},
		{
			MethodName: "CreateCollection",
			Handler:    _CollectionService_CreateCollection_Handler,
		},
		{
			MethodName: "UpdateCollection",
			Handler:    _CollectionService_UpdateCollection_Handler,
		},
		{
			MethodName: "DeleteCollection",
			Handler:    _CollectionService_DeleteCollection_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "worker/database/v1/collection_service.proto",
}
