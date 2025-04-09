// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.4
// - protoc             (unknown)
// source: message/user_service.proto

package climsg

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

const OperationUserServiceLogin = "/message.UserService/Login"
const OperationUserServiceSetGender = "/message.UserService/SetGender"
const OperationUserServiceUpdateName = "/message.UserService/UpdateName"

type UserServiceHTTPServer interface {
	// Login Login
	Login(context.Context, *CSLogin) (*SCLogin, error)
	// SetGender Set gender
	SetGender(context.Context, *CSSetGender) (*SCSetGender, error)
	// UpdateName Update name
	UpdateName(context.Context, *CSUpdateName) (*SCUpdateName, error)
}

func RegisterUserServiceHTTPServer(s *http.Server, srv UserServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/user/login", _UserService_Login0_HTTP_Handler(srv))
	r.POST("/user/name/update", _UserService_UpdateName0_HTTP_Handler(srv))
	r.POST("/user/gender/set", _UserService_SetGender0_HTTP_Handler(srv))
}

func _UserService_Login0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CSLogin
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceLogin)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Login(ctx, req.(*CSLogin))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SCLogin)
		return ctx.Result(200, reply)
	}
}

func _UserService_UpdateName0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CSUpdateName
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceUpdateName)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.UpdateName(ctx, req.(*CSUpdateName))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SCUpdateName)
		return ctx.Result(200, reply)
	}
}

func _UserService_SetGender0_HTTP_Handler(srv UserServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in CSSetGender
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationUserServiceSetGender)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.SetGender(ctx, req.(*CSSetGender))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*SCSetGender)
		return ctx.Result(200, reply)
	}
}

type UserServiceHTTPClient interface {
	Login(ctx context.Context, req *CSLogin, opts ...http.CallOption) (rsp *SCLogin, err error)
	SetGender(ctx context.Context, req *CSSetGender, opts ...http.CallOption) (rsp *SCSetGender, err error)
	UpdateName(ctx context.Context, req *CSUpdateName, opts ...http.CallOption) (rsp *SCUpdateName, err error)
}

type UserServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewUserServiceHTTPClient(client *http.Client) UserServiceHTTPClient {
	return &UserServiceHTTPClientImpl{client}
}

func (c *UserServiceHTTPClientImpl) Login(ctx context.Context, in *CSLogin, opts ...http.CallOption) (*SCLogin, error) {
	var out SCLogin
	pattern := "/user/login"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationUserServiceLogin))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *UserServiceHTTPClientImpl) SetGender(ctx context.Context, in *CSSetGender, opts ...http.CallOption) (*SCSetGender, error) {
	var out SCSetGender
	pattern := "/user/gender/set"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationUserServiceSetGender))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *UserServiceHTTPClientImpl) UpdateName(ctx context.Context, in *CSUpdateName, opts ...http.CallOption) (*SCUpdateName, error) {
	var out SCUpdateName
	pattern := "/user/name/update"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationUserServiceUpdateName))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
