// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.4
// - protoc             (unknown)
// source: message/equip_service.proto

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

const OperationEquipServiceEquipTakeOff = "/message.EquipService/EquipTakeOff"
const OperationEquipServiceEquipUpgrade = "/message.EquipService/EquipUpgrade"
const OperationEquipServiceEquipWear = "/message.EquipService/EquipWear"

type EquipServiceHTTPServer interface {
	// EquipTakeOff Take off equipment
	EquipTakeOff(context.Context, *CSEquipTakeOff) (*SCEquipTakeOff, error)
	// EquipUpgrade Equipment upgrade
	EquipUpgrade(context.Context, *CSEquipUpgrade) (*SCEquipUpgrade, error)
	// EquipWear Wear equipment
	EquipWear(context.Context, *CSEquipWear) (*SCEquipWear, error)
}

func RegisterEquipServiceHTTPServer(s *http.Server, srv EquipServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/equip/wear", _EquipService_EquipWear0_HTTP_Handler(srv))
	r.POST("/equip/takeoff", _EquipService_EquipTakeOff0_HTTP_Handler(srv))
	r.POST("/equip/upgrade", _EquipService_EquipUpgrade0_HTTP_Handler(srv))
}

func _EquipService_EquipWear0_HTTP_Handler(srv EquipServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CSEquipWear
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationEquipServiceEquipWear)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.EquipWear(ctx, req.(*CSEquipWear))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SCEquipWear)
		return ctx.Result(200, reply)
	}
}

func _EquipService_EquipTakeOff0_HTTP_Handler(srv EquipServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CSEquipTakeOff
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationEquipServiceEquipTakeOff)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.EquipTakeOff(ctx, req.(*CSEquipTakeOff))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SCEquipTakeOff)
		return ctx.Result(200, reply)
	}
}

func _EquipService_EquipUpgrade0_HTTP_Handler(srv EquipServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CSEquipUpgrade
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationEquipServiceEquipUpgrade)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.EquipUpgrade(ctx, req.(*CSEquipUpgrade))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SCEquipUpgrade)
		return ctx.Result(200, reply)
	}
}

type EquipServiceHTTPClient interface {
	EquipTakeOff(ctx context.Context, req *CSEquipTakeOff, opts ...http.CallOption) (rsp *SCEquipTakeOff, err error)
	EquipUpgrade(ctx context.Context, req *CSEquipUpgrade, opts ...http.CallOption) (rsp *SCEquipUpgrade, err error)
	EquipWear(ctx context.Context, req *CSEquipWear, opts ...http.CallOption) (rsp *SCEquipWear, err error)
}

type EquipServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewEquipServiceHTTPClient(client *http.Client) EquipServiceHTTPClient {
	return &EquipServiceHTTPClientImpl{client}
}

func (c *EquipServiceHTTPClientImpl) EquipTakeOff(ctx context.Context, in *CSEquipTakeOff, opts ...http.CallOption) (*SCEquipTakeOff, error) {
	var out SCEquipTakeOff
	pattern := "/equip/takeoff"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationEquipServiceEquipTakeOff))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *EquipServiceHTTPClientImpl) EquipUpgrade(ctx context.Context, in *CSEquipUpgrade, opts ...http.CallOption) (*SCEquipUpgrade, error) {
	var out SCEquipUpgrade
	pattern := "/equip/upgrade"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationEquipServiceEquipUpgrade))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *EquipServiceHTTPClientImpl) EquipWear(ctx context.Context, in *CSEquipWear, opts ...http.CallOption) (*SCEquipWear, error) {
	var out SCEquipWear
	pattern := "/equip/wear"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationEquipServiceEquipWear))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
