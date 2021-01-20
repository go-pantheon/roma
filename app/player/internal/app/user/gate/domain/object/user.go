package userobj

import (
	"context"
	"time"

	"github.com/pkg/errors"
	basicobj "github.com/vulcan-frame/vulcan-game/app/player/internal/app/basic/gate/domain/object"
	storageobj "github.com/vulcan-frame/vulcan-game/app/player/internal/app/storage/gate/domain/object"
	message "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
	dbv1 "github.com/vulcan-frame/vulcan-game/gen/api/db/player/v1"
	"github.com/vulcan-frame/vulcan-kit/profile"
	"github.com/vulcan-frame/vulcan-kit/version"
	"github.com/vulcan-frame/vulcan-util/xtime"
)

type User struct {
	Id           int64
	Name         string
	LoginAt      time.Time
	LastOnlineAt time.Time
	LogoutAt     time.Time
	CreatedAt    time.Time
	LastOnlineIP string
	Version      int64
	newborn      bool

	NextDailyResetAt    time.Time
	DailyOnlineDuration time.Duration
	TotalOnlineDuration time.Duration
	ServerVersion       string

	Basic *basicobj.Basic

	Dev     *Dev
	System  *System
	Storage *storageobj.Storage
}

func NewUser(id int64, name string) *User {
	u := &User{
		Id:   id,
		Name: name,
	}

	u.Basic = basicobj.NewBasic()
	u.Dev = NewDev()
	u.System = NewSystem()
	u.Storage = storageobj.NewStorage()

	return u
}

func NewUserProto() *dbv1.UserProto {
	p := &dbv1.UserProto{}
	p.Basic = basicobj.NewBasicProto()
	p.Dev = NewDevProto()
	p.System = NewSystemProto()
	p.Storage = storageobj.NewStorageProto()
	return p
}

func (o *User) EncodeServer(p *dbv1.UserProto) *dbv1.UserProto {
	p.Id = o.Id
	p.Name = o.Name

	p.LoginAt = o.LoginAt.Unix()
	p.LogoutAt = o.LogoutAt.Unix()
	p.LastOnlineAt = o.LastOnlineAt.Unix()
	p.CreatedAt = o.CreatedAt.Unix()
	p.LastOnlineIp = o.LastOnlineIP
	p.Version = o.Version

	p.NextDailyResetAt = o.NextDailyResetAt.Unix()
	p.DailyOnlineSeconds = int64(o.DailyOnlineDuration.Seconds())
	p.TotalOnlineSeconds = int64(o.TotalOnlineDuration.Seconds())
	p.ServerVersion = o.ServerVersion

	o.Basic.EncodeServer(p.Basic)
	o.Dev.EncodeServer(p.Dev)
	o.System.EncodeServer(p.System)
	o.Storage.EncodeServer(p.Storage)

	return p
}

func (o *User) DecodeServer(ctx context.Context, p *dbv1.UserProto) (err error) {
	if p == nil {
		return
	}

	newServerVersion, err := checkServerVersion(profile.Version(), p.ServerVersion)
	if err != nil {
		return errors.WithMessagef(err, "serverVersion unmarshal failed. uid=%d", o.Id)
	}
	o.ServerVersion = newServerVersion

	o.Id = p.Id
	o.Name = p.Name
	o.LoginAt = xtime.Time(p.LoginAt)
	o.LogoutAt = xtime.Time(p.LogoutAt)
	o.LastOnlineAt = xtime.Time(p.LastOnlineAt)
	o.CreatedAt = xtime.Time(p.CreatedAt)
	o.LastOnlineIP = p.LastOnlineIp
	o.Version = p.Version

	o.NextDailyResetAt = xtime.Time(p.NextDailyResetAt)
	o.DailyOnlineDuration = time.Duration(p.DailyOnlineSeconds) * time.Second
	o.TotalOnlineDuration = time.Duration(p.TotalOnlineSeconds) * time.Second

	if err = o.Basic.DecodeServer(p.Basic); err != nil {
		return errors.WithMessagef(err, "basic unmarshal failed. uid=%d", o.Id)
	}
	o.Dev.DecodeServer(p.Dev)
	o.System.DecodeServer(p.System)
	o.Storage.DecodeServer(ctx, p.Storage)
	return nil
}

func checkServerVersion(serverVersion, userVersion string) (newUserVersion string, err error) {
	az1, ssv, isRelease := version.GetSubVersion(profile.Version())
	if !isRelease {
		newUserVersion = userVersion
		return
	}
	az2, usv, isRelease := version.GetSubVersion(userVersion)
	if !isRelease {
		newUserVersion = serverVersion
		return
	}

	if az1 != az2 {
		err = errors.Errorf("userVersion is not equal to serverVersion. userVersion=%s serverVersion=%s", userVersion, serverVersion)
		return
	}

	for i := 0; i < len(ssv); i++ {
		if ssv[i] < usv[i] {
			err = errors.Errorf("userVersion is greater than serverVersion. userVersion=%s serverVersion=%s", userVersion, serverVersion)
			return
		}
		if ssv[i] > usv[i] {
			newUserVersion = serverVersion
			return
		}
	}
	newUserVersion = userVersion
	return
}

func (o *User) EncodeClient() *message.UserProto {
	p := &message.UserProto{
		Basic:   o.Basic.EncodeClient(),
		Storage: o.Storage.EncodeClient(),
	}
	return p
}

func (o *User) SetNewborn(b bool) {
	o.newborn = b
}

func (o *User) GetAndResetNewborn() bool {
	b := o.newborn
	o.newborn = false
	return b
}
