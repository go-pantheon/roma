package object

import (
	"time"

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
	o := &Status{}
	return o
}

func (o *Status) Marshal() ([]byte, error) {
	p := dbv1.UserStatusProtoPool.Get()
	defer dbv1.UserStatusProtoPool.Put(p)

	p.OnlineIp = o.ClientIP
	p.LoginAt = o.LoginAt.Unix()
	p.LogoutAt = o.LogoutAt.Unix()
	p.LatestOnlineAt = o.LatestOnlineAt.Unix()

	p.NextDailyResetAt = o.NextDailyResetAt.Unix()
	p.DailyOnlineSeconds = int64(o.DailyOnlineDuration.Seconds())
	p.TotalOnlineSeconds = int64(o.TotalOnlineDuration.Seconds())

	return proto.Marshal(p)
}

func (o *Status) Unmarshal(bytes []byte) (err error) {
	p := dbv1.UserStatusProtoPool.Get()
	defer dbv1.UserStatusProtoPool.Put(p)

	if err = proto.Unmarshal(bytes, p); err != nil {
		return err
	}

	o.ClientIP = p.OnlineIp
	o.LoginAt = time.Unix(p.LoginAt, 0)
	o.LogoutAt = time.Unix(p.LogoutAt, 0)
	o.LatestOnlineAt = time.Unix(p.LatestOnlineAt, 0)
	o.NextDailyResetAt = time.Unix(p.NextDailyResetAt, 0)
	o.DailyOnlineDuration = time.Duration(p.DailyOnlineSeconds) * time.Second
	o.TotalOnlineDuration = time.Duration(p.TotalOnlineSeconds) * time.Second

	return
}
