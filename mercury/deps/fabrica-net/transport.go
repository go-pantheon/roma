package net

import (
	"github.com/go-kratos/kratos/v2/transport"
	"google.golang.org/grpc/metadata"
)

type NetKind string

const (
	NetKindTCP       NetKind = "tcp"
	NetKindKCP       NetKind = "kcp"
	NetKindWebSocket NetKind = "ws"
)

var (
	_ transport.Transporter = (*Transport)(nil)
)

type Transport struct {
	endpoint      string
	operation     string
	requestHeader HeaderCarrier
	replyHeader   HeaderCarrier
}

func NewTransport(endpoint string,
	operation string,
	requestHeader HeaderCarrier,
	replyHeader HeaderCarrier,
) *Transport {
	return &Transport{
		endpoint:      endpoint,
		operation:     operation,
		requestHeader: requestHeader,
		replyHeader:   replyHeader,
	}
}

func (tr *Transport) Kind() transport.Kind {
	return transport.KindGRPC
}

func (tr *Transport) Endpoint() string {
	return tr.endpoint
}

func (tr *Transport) Operation() string {
	return tr.operation
}

func (tr *Transport) RequestHeader() transport.Header {
	return tr.requestHeader
}

func (tr *Transport) ReplyHeader() transport.Header {
	return tr.replyHeader
}

type HeaderCarrier metadata.MD

func (mc HeaderCarrier) Get(key string) string {
	vals := metadata.MD(mc).Get(key)
	if len(vals) > 0 {
		return vals[0]
	}
	return ""
}

func (mc HeaderCarrier) Set(key string, value string) {
	metadata.MD(mc).Set(key, value)
}

func (mc HeaderCarrier) Keys() []string {
	keys := make([]string, 0, len(mc))
	for k := range metadata.MD(mc) {
		keys = append(keys, k)
	}
	return keys
}

func (mc HeaderCarrier) Add(key string, value string) {
	metadata.MD(mc).Append(key, value)
}

func (mc HeaderCarrier) Values(key string) []string {
	return metadata.MD(mc).Get(key)
}
