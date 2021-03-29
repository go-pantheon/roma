// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.8.4
// - protoc             (unknown)
// source: gate/service/push/v1/push.proto

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

const OperationPushServiceBroadcast = "/gate.service.push.v1.PushService/Broadcast"
const OperationPushServiceMulticast = "/gate.service.push.v1.PushService/Multicast"
const OperationPushServicePush = "/gate.service.push.v1.PushService/Push"

type PushServiceHTTPServer interface {
	Broadcast(context.Context, *BroadcastRequest) (*BroadcastResponse, error)
	Multicast(context.Context, *MulticastRequest) (*MulticastResponse, error)
	Push(context.Context, *PushRequest) (*PushResponse, error)
}

func RegisterPushServiceHTTPServer(s *http.Server, srv PushServiceHTTPServer) {
	r := s.Route("/")
	r.POST("/push", _PushService_Push0_HTTP_Handler(srv))
	r.POST("/multicast", _PushService_Multicast0_HTTP_Handler(srv))
	r.POST("/broadcast", _PushService_Broadcast0_HTTP_Handler(srv))
}

func _PushService_Push0_HTTP_Handler(srv PushServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in PushRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPushServicePush)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Push(ctx, req.(*PushRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*PushResponse)
		return ctx.Result(200, reply)
	}
}

func _PushService_Multicast0_HTTP_Handler(srv PushServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in MulticastRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPushServiceMulticast)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Multicast(ctx, req.(*MulticastRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*MulticastResponse)
		return ctx.Result(200, reply)
	}
}

func _PushService_Broadcast0_HTTP_Handler(srv PushServiceHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in BroadcastRequest
		if err := ctx.Bind(&in); err != nil {
			return err
		}
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationPushServiceBroadcast)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.Broadcast(ctx, req.(*BroadcastRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*BroadcastResponse)
		return ctx.Result(200, reply)
	}
}

type PushServiceHTTPClient interface {
	Broadcast(ctx context.Context, req *BroadcastRequest, opts ...http.CallOption) (rsp *BroadcastResponse, err error)
	Multicast(ctx context.Context, req *MulticastRequest, opts ...http.CallOption) (rsp *MulticastResponse, err error)
	Push(ctx context.Context, req *PushRequest, opts ...http.CallOption) (rsp *PushResponse, err error)
}

type PushServiceHTTPClientImpl struct {
	cc *http.Client
}

func NewPushServiceHTTPClient(client *http.Client) PushServiceHTTPClient {
	return &PushServiceHTTPClientImpl{client}
}

func (c *PushServiceHTTPClientImpl) Broadcast(ctx context.Context, in *BroadcastRequest, opts ...http.CallOption) (*BroadcastResponse, error) {
	var out BroadcastResponse
	pattern := "/broadcast"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationPushServiceBroadcast))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *PushServiceHTTPClientImpl) Multicast(ctx context.Context, in *MulticastRequest, opts ...http.CallOption) (*MulticastResponse, error) {
	var out MulticastResponse
	pattern := "/multicast"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationPushServiceMulticast))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}

func (c *PushServiceHTTPClientImpl) Push(ctx context.Context, in *PushRequest, opts ...http.CallOption) (*PushResponse, error) {
	var out PushResponse
	pattern := "/push"
	path := binding.EncodeURL(pattern, in, false)
	opts = append(opts, http.Operation(OperationPushServicePush))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "POST", path, in, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, nil
}
