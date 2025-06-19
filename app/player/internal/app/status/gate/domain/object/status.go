package object

import (
	"time"

	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/userregister"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"google.golang.org/protobuf/proto"
)

const (
	ModuleKey = "status"
)

func init() {
	userregister.Register(ModuleKey, NewStatus)
}

var _ life.Module = (*Status)(nil)

type Status struct {
	LoginAt        time.Time
	LogoutAt       time.Time
	LatestOnlineAt time.Time
	ClientIP       string

	NextDailyResetAt    time.Time
	DailyOnlineDuration time.Duration
	TotalOnlineDuration time.Duration
}

func NewStatus() life.Module {
	return &Status{}
}

func (o *Status) EncodeServer() proto.Message {
	p := dbv1.UserStatusProtoPool.Get()

	p.OnlineIp = o.ClientIP
	p.LoginAt = o.LoginAt.Unix()
	p.LogoutAt = o.LogoutAt.Unix()
	p.LatestOnlineAt = o.LatestOnlineAt.Unix()

	p.NextDailyResetAt = o.NextDailyResetAt.Unix()
	p.DailyOnlineSeconds = int64(o.DailyOnlineDuration.Seconds())
	p.TotalOnlineSeconds = int64(o.TotalOnlineDuration.Seconds())

	return p
}

func (o *Status) DecodeServer(p proto.Message) error {
	if p == nil {
		return errors.New("status decode server nil")
	}

	op, ok := p.(*dbv1.UserStatusProto)
	if !ok {
		return errors.Errorf("status decode server invalid type: %T", p)
	}

	o.ClientIP = op.OnlineIp
	o.LoginAt = time.Unix(op.LoginAt, 0)
	o.LogoutAt = time.Unix(op.LogoutAt, 0)
	o.LatestOnlineAt = time.Unix(op.LatestOnlineAt, 0)
	o.NextDailyResetAt = time.Unix(op.NextDailyResetAt, 0)
	o.DailyOnlineDuration = time.Duration(op.DailyOnlineSeconds) * time.Second
	o.TotalOnlineDuration = time.Duration(op.TotalOnlineSeconds) * time.Second

	return nil
}
