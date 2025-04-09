package xcontext

import (
	"context"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-pantheon/fabrica-util/xsync"
	"github.com/pkg/errors"
)

// Use the custom type for your constants
const (
	CtxSID         = "x-md-global-sid" // Server ID
	CtxUID         = "x-md-global-uid" // User ID
	CtxOID         = "x-md-global-oid" // Object ID
	CtxColor       = "x-md-global-color"
	CtxStatus      = "x-md-global-status"
	CtxReferer     = "x-md-global-referer" // example: gate:10.0.1.31 or player:10.0.2.31
	CtxClientIP    = "x-md-global-client-ip"
	CtxGateReferer = "x-md-global-gate-referer" // example: 10.0.1.31:9100#10001
)

var Keys = []string{CtxSID, CtxUID, CtxOID, CtxStatus, CtxColor, CtxReferer, CtxClientIP, CtxGateReferer}

func SetColor(ctx context.Context, color string) context.Context {
	return metadata.AppendToClientContext(ctx, string(CtxColor), color)
}

func Color(ctx context.Context) string {
	if md, ok := metadata.FromServerContext(ctx); ok {
		return md.Get(CtxColor)
	}
	return ""
}

func SetUID(ctx context.Context, id int64) context.Context {
	return metadata.AppendToClientContext(ctx, CtxUID, strconv.FormatInt(id, 10))
}

func UID(ctx context.Context) (int64, error) {
	md, ok := metadata.FromServerContext(ctx)
	if !ok {
		return 0, errors.New("metadata not in context")
	}

	str := md.Get(CtxUID)
	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "uid must be int64, uid=%s", str)
	}
	return id, nil
}

func SetOID(ctx context.Context, id int64) context.Context {
	return metadata.AppendToClientContext(ctx, CtxOID, strconv.FormatInt(id, 10))
}

func OID(ctx context.Context) (int64, error) {
	md, ok := metadata.FromServerContext(ctx)
	if !ok {
		return 0, errors.New("metadata not in context")
	}

	str := md.Get(CtxOID)
	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "oid must be int64, oid=%s", str)
	}
	return id, nil
}

func SetSID(ctx context.Context, id int64) context.Context {
	return metadata.AppendToClientContext(ctx, CtxSID, strconv.FormatInt(id, 10))
}

func SID(ctx context.Context) (int64, error) {
	md, ok := metadata.FromServerContext(ctx)
	if !ok {
		return 0, errors.New("metadata not in context")
	}

	str := md.Get(CtxSID)
	id, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0, errors.Wrapf(err, "sid must be int64, sid=%s", str)
	}
	return id, nil
}

func SetStatus(ctx context.Context, status int64) context.Context {
	if status == 0 {
		return ctx
	}
	return metadata.AppendToClientContext(ctx, CtxStatus, strconv.FormatInt(status, 10))
}

func Status(ctx context.Context) int64 {
	if md, ok := metadata.FromServerContext(ctx); ok {
		v := md.Get(CtxStatus)
		status, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			log.Errorf("status must be int64, status=%s", v)
			return 0
		}
		return status
	}
	return 0
}

func SetClientIP(ctx context.Context, ip string) context.Context {
	if len(ip) == 0 {
		return ctx
	}
	return metadata.AppendToClientContext(ctx, CtxClientIP, strings.Split(ip, ":")[0])
}

func ClientIP(ctx context.Context) string {
	if md, ok := metadata.FromServerContext(ctx); ok {
		return md.Get(CtxClientIP)
	}
	return ""
}

func SetGateReferer(ctx context.Context, server string, wid uint64) context.Context {
	if len(server) == 0 {
		return ctx
	}
	return metadata.AppendToClientContext(ctx, CtxGateReferer, fmt.Sprintf("%s#%d", server, wid))
}

func GateReferer(ctx context.Context) string {
	if md, ok := metadata.FromServerContext(ctx); ok {
		return md.Get(CtxGateReferer)
	}
	return ""
}

func RemoteAddr(conn net.Conn) string {
	if conn == nil {
		return ""
	}
	addr := conn.RemoteAddr()
	if addr == nil {
		return ""
	}
	return addr.String()
}

func LocalAddr(conn net.Conn) string {
	if conn == nil {
		return ""
	}
	addr := conn.LocalAddr()
	if addr == nil {
		return ""
	}
	return addr.String()
}

// ReadDeadline https://github.com/google/mtail/commit/8dd02e80f9e42eebb59fee10c24c7cc686f9e481
type ReadDeadline interface {
	SetReadDeadline(t time.Time) error
}

type Deadline interface {
	SetDeadline(t time.Time) error
}

type DeadlineSetter interface {
	xsync.WaitStoppable
}

// SetDeadlineWithContext use context to control the deadline of the connection
func SetDeadlineWithContext(ctx context.Context, d Deadline, tag string) {
	go func() {
		<-ctx.Done()
		log.Debugf("[net.SetDeadlineWithContext] %s start to close", tag)
		if err := d.SetDeadline(time.Now()); err != nil {
			log.Errorf("[net.SetDeadlineWithContext] %s close failed. %+v", tag, err)
		}
	}()
}

// CloseOnCancel close the connection when the context is canceled
func CloseOnCancel(ctx context.Context, closer io.Closer, tag string) {
	go func() {
		<-ctx.Done()
		log.Debugf("[net.CloseOnCancel] %s start to close", tag)
		if err := closer.Close(); err != nil {
			log.Errorf("[net.CloseOnCancel] %s close failed. %+v", tag, err)
		}
	}()
}

// SetDeadlineWithTimeout set the deadline to close the connection after the specified timeout
func SetDeadlineWithTimeout(d Deadline, timeout time.Duration, tag string) {
	go func() {
		timer := time.NewTimer(timeout)
		<-timer.C
		log.Debugf("[net.SetDeadlineWithTimeout] %s start to close", tag)
		if err := d.SetDeadline(time.Now()); err != nil {
			log.Errorf("[net.SetDeadlineWithTimeout] %s close failed. %+v", tag, err)
		}
	}()
}
