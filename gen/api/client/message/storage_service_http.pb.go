// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.4
// - protoc             (unknown)
// source: message/storage_service.proto

package climsg

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

const OperationStorageServiceUsePack = "/message.StorageService/UsePack"

type StorageServiceHTTPServer interface {
	// UsePack Use pack
	UsePack(context.Context, *CSUsePack) (*SCUsePack, error)
}

func RegisterStorageServiceHTTPServer(s *http.Server, srv StorageServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/storage/use/pack", _StorageService_UsePack0_HTTP_Handler(srv))
}

func _StorageService_UsePack0_HTTP_Handler(srv StorageServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CSUsePack
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationStorageServiceUsePack)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UsePack(ctx, req.(*CSUsePack))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SCUsePack)
		return ctx.Result(200, reply)
	}
}

type StorageServiceHTTPClient interface {
	UsePack(ctx context.Context, req *CSUsePack, opts ...http.CallOption) (rsp *SCUsePack, err error)
}

type StorageServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewStorageServiceHTTPClient(client *http.Client) StorageServiceHTTPClient {
	return &StorageServiceHTTPClientImpl{client}
}

func (c *StorageServiceHTTPClientImpl) UsePack(ctx context.Context, in *CSUsePack, opts ...http.CallOption) (*SCUsePack, error) {
	var out SCUsePack
	pattern := "/storage/use/pack"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationStorageServiceUsePack))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
