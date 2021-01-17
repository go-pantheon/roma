package dev

import (
	"context"

	"github.com/go-kratos/kratos/middleware"
	"github.com/go-kratos/kratos/transport"
	"github.com/vulcan-frame/vulcan-kit/profile"
	"github.com/vulcan-frame/vulcan-kit/xcontext"
	"google.golang.org/grpc/metadata"
)

func Server() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			if !profile.IsDev() {
				return handler(ctx, req)
			}
			ctx = TransformContext(ctx)
			reply, err := handler(ctx, req)
			return reply, err
		}
	}
}

func TransformContext(ctx context.Context) context.Context {
	if !profile.IsDev() {
		return ctx
	}

	if info, ok := transport.FromServerContext(ctx); ok {
		pairs := make([]string, 0, len(xcontext.Keys))
		for _, k := range xcontext.Keys {
			v := info.RequestHeader().Get(k)
			pairs = append(pairs, k, v)
		}
		ctx = metadata.NewIncomingContext(ctx, metadata.Pairs(pairs...))
	}

	return ctx
}
