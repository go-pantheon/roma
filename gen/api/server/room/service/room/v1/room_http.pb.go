// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.4
// - protoc             (unknown)
// source: room/service/room/v1/room.proto

package servicev1

import (
	context "context"
	http "github.com/go-kratos/kratos/transport/http"
	binding "github.com/go-kratos/kratos/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationRoomServiceGetById = "/room.service.room.v1.RoomService/GetById"
const OperationRoomServiceListById = "/room.service.room.v1.RoomService/ListById"

type RoomServiceHTTPServer interface {
	// GetById Get room data cache by id
	GetById(context.Context, *GetByIdRequest) (*GetByIdResponse, error)
	// ListById Get room data cache list by id list
	ListById(context.Context, *ListByIdRequest) (*ListByIdResponse, error)
}

func RegisterRoomServiceHTTPServer(s *http.Server, srv RoomServiceHTTPServer) {
	r := s.Route("/")
	r.GET("/service/room/id", _RoomService_GetById0_HTTP_Handler(srv))
	r.GET("/service/room/list/id", _RoomService_ListById0_HTTP_Handler(srv))
}

func _RoomService_GetById0_HTTP_Handler(srv RoomServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in GetByIdRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoomServiceGetById)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.GetById(ctx, req.(*GetByIdRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*GetByIdResponse)
		return ctx.Result(200, reply)
	}
}

func _RoomService_ListById0_HTTP_Handler(srv RoomServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in ListByIdRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationRoomServiceListById)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.ListById(ctx, req.(*ListByIdRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*ListByIdResponse)
		return ctx.Result(200, reply)
	}
}

type RoomServiceHTTPClient interface {
	GetById(ctx context.Context, req *GetByIdRequest, opts ...http.CallOption) (rsp *GetByIdResponse, err error)
	ListById(ctx context.Context, req *ListByIdRequest, opts ...http.CallOption) (rsp *ListByIdResponse, err error)
}

type RoomServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewRoomServiceHTTPClient(client *http.Client) RoomServiceHTTPClient {
	return &RoomServiceHTTPClientImpl{client}
}

func (c *RoomServiceHTTPClientImpl) GetById(ctx context.Context, in *GetByIdRequest, opts ...http.CallOption) (*GetByIdResponse, error) {
	var out GetByIdResponse
	pattern := "/service/room/id"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationRoomServiceGetById))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *RoomServiceHTTPClientImpl) ListById(ctx context.Context, in *ListByIdRequest, opts ...http.CallOption) (*ListByIdResponse, error) {
	var out ListByIdResponse
	pattern := "/service/room/list/id"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationRoomServiceListById))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
