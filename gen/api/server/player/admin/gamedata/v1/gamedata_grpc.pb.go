// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: player/admin/gamedata/v1/gamedata.proto

package adminv1

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	GamedataAdmin_GetItemList_FullMethodName = "/player.admin.gamedata.v1.GamedataAdmin/GetItemList"
)

// GamedataAdminClient is the client API for GamedataAdmin service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// Gamedata admin service
// Open to the server cluster
// Provide HTTP and gRPC interfaces
type GamedataAdminClient interface {
	// Query all configuration items
	GetItemList(ctx context.Context, in *GetItemListRequest, opts ...grpc.CallOption) (*GetItemListResponse, error)
}

type gamedataAdminClient struct {
	cc grpc.ClientConnInterface
}

func NewGamedataAdminClient(cc grpc.ClientConnInterface) GamedataAdminClient {
	return &gamedataAdminClient{cc}
}

func (c *gamedataAdminClient) GetItemList(ctx context.Context, in *GetItemListRequest, opts ...grpc.CallOption) (*GetItemListResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetItemListResponse)
	err := c.cc.Invoke(ctx, GamedataAdmin_GetItemList_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GamedataAdminServer is the server API for GamedataAdmin service.
// All implementations must embed UnimplementedGamedataAdminServer
// for forward compatibility.
//
// Gamedata admin service
// Open to the server cluster
// Provide HTTP and gRPC interfaces
type GamedataAdminServer interface {
	// Query all configuration items
	GetItemList(context.Context, *GetItemListRequest) (*GetItemListResponse, error)
	mustEmbedUnimplementedGamedataAdminServer()
}

// UnimplementedGamedataAdminServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedGamedataAdminServer struct{}

func (UnimplementedGamedataAdminServer) GetItemList(context.Context, *GetItemListRequest) (*GetItemListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetItemList not implemented")
}
func (UnimplementedGamedataAdminServer) mustEmbedUnimplementedGamedataAdminServer() {}
func (UnimplementedGamedataAdminServer) testEmbeddedByValue()                       {}

// UnsafeGamedataAdminServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to GamedataAdminServer will
// result in compilation errors.
type UnsafeGamedataAdminServer interface {
	mustEmbedUnimplementedGamedataAdminServer()
}

func RegisterGamedataAdminServer(s grpc.ServiceRegistrar, srv GamedataAdminServer) {
	// If the following call pancis, it indicates UnimplementedGamedataAdminServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&GamedataAdmin_ServiceDesc, srv)
}

func _GamedataAdmin_GetItemList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetItemListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(GamedataAdminServer).GetItemList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: GamedataAdmin_GetItemList_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(GamedataAdminServer).GetItemList(ctx, req.(*GetItemListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// GamedataAdmin_ServiceDesc is the grpc.ServiceDesc for GamedataAdmin service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var GamedataAdmin_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "player.admin.gamedata.v1.GamedataAdmin",
	HandlerType: (*GamedataAdminServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetItemList",
			Handler:    _GamedataAdmin_GetItemList_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "player/admin/gamedata/v1/gamedata.proto",
}
