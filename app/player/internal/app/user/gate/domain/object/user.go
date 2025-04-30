package userobj

import (
	"context"
	"time"

	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/version"
	"github.com/go-pantheon/fabrica-util/xtime"
	basicobj "github.com/go-pantheon/roma/app/player/internal/app/basic/gate/domain/object"
	heroobj "github.com/go-pantheon/roma/app/player/internal/app/hero/gate/domain/object"
	plunderobj "github.com/go-pantheon/roma/app/player/internal/app/plunder/gate/domain/object"
	roomobj "github.com/go-pantheon/roma/app/player/internal/app/room/gate/domain/object"
	storageobj "github.com/go-pantheon/roma/app/player/internal/app/storage/gate/domain/object"
	message "github.com/go-pantheon/roma/gen/api/client/message"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
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

	Dev      *Dev
	System   *System
	Plunders *plunderobj.Plunders
	Storage  *storageobj.Storage
	HeroList *heroobj.HeroList
	Room     *roomobj.Room
}

func NewUser(id int64, name string) *User {
	u := &User{
		Id:   id,
		Name: name,
	}

	u.Basic = basicobj.NewBasic()
	u.Dev = NewDev()
	u.System = NewSystem()
	u.Plunders = plunderobj.NewPlunders()
	u.Storage = storageobj.NewStorage()
	u.HeroList = heroobj.NewHeroList()
	u.Room = roomobj.NewRoom()

	return u
}

func NewUserProto() *dbv1.UserProto {
	p := &dbv1.UserProto{}
	p.Basic = basicobj.NewBasicProto()
	p.Dev = NewDevProto()
	p.System = NewSystemProto()
	p.Plunders = plunderobj.NewPlundersProto()
	p.Storage = storageobj.NewStorageProto()
	p.HeroList = heroobj.NewHeroListProto()
	p.Room = roomobj.NewRoomProto()
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
	o.Plunders.EncodeServer(p.Plunders)
	o.Storage.EncodeServer(p.Storage)
	o.HeroList.EncodeServer(p.HeroList)
	o.Room.EncodeServer(p.Room)

	return p
}

// UserProtoSize returns the estimated memory size of UserProto
// TODO: optimize the size calculation
func UserProtoSize(p *dbv1.UserProto) int {
	return proto.Size(p) * 10
}

// UserSize returns the estimated memory size of User
// TODO: optimize the size calculation
func UserSize(o *User) int {
	return 512
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
	o.Plunders.DecodeServer(p.Plunders)
	o.Storage.DecodeServer(ctx, p.Storage)
	o.Room.DecodeServer(p.Room)
	if err = o.HeroList.DecodeServer(ctx, p.HeroList); err != nil {
		return errors.WithMessagef(err, "heroList unmarshal failed. uid=%d", o.Id)
	}
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
		Basic:    o.Basic.EncodeClient(),
		Storage:  o.Storage.EncodeClient(),
		HeroList: o.HeroList.EncodeClient(),
		Room:     o.Room.EncodeClient(),
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
