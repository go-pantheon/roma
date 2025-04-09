package dev

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
	"github.com/go-pantheon/fabrica-net/xcontext"
	"google.golang.org/grpc/metadata"
)

func Server(l log.Logger) middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			ctx = TransformContext(ctx)
			return handler(ctx, req)
		}
	}
}

func TransformContext(ctx context.Context) context.Context {
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

func IsAdminPath(ctx context.Context) bool {
	tp, ok := transport.FromServerContext(ctx)
	if !ok {
		return false
	}
	info, ok := tp.(*http.Transport)
	if !ok {
		return false
	}
	return strings.Index(info.Request().RequestURI, "/admin") == 0
}
