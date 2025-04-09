// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.4
// - protoc             (unknown)
// source: player/admin/storage/v1/storage.proto

package adminv1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationStorageAdminAddItem = "/player.admin.storage.v1.StorageAdmin/AddItem"

type StorageAdminHTTPServer interface {
	AddItem(context.Context, *AddItemRequest) (*AddItemResponse, error)
}

func RegisterStorageAdminHTTPServer(s *http.Server, srv StorageAdminHTTPServer) {
	r := s.Route("/")
	r.POST("/admin/storage/item/add", _StorageAdmin_AddItem0_HTTP_Handler(srv))
}

func _StorageAdmin_AddItem0_HTTP_Handler(srv StorageAdminHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in AddItemRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationStorageAdminAddItem)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.AddItem(ctx, req.(*AddItemRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*AddItemResponse)
		return ctx.Result(200, reply)
	}
}

type StorageAdminHTTPClient interface {
	AddItem(ctx context.Context, req *AddItemRequest, opts ...http.CallOption) (rsp *AddItemResponse, err error)
}

type StorageAdminHTTPClientImpl struct {
	cc *http.Client
}

func NewStorageAdminHTTPClient(client *http.Client) StorageAdminHTTPClient {
	return &StorageAdminHTTPClientImpl{client}
}

func (c *StorageAdminHTTPClientImpl) AddItem(ctx context.Context, in *AddItemRequest, opts ...http.CallOption) (*AddItemResponse, error) {
	var out AddItemResponse
	pattern := "/admin/storage/item/add"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationStorageAdminAddItem))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
