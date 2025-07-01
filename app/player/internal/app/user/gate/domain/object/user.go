package userobj

import (
	"time"

	"github.com/go-pantheon/fabrica-kit/profile"
	"github.com/go-pantheon/fabrica-kit/version"
	"github.com/go-pantheon/fabrica-util/xid"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/userregister"
	message "github.com/go-pantheon/roma/gen/api/client/message"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"github.com/pkg/errors"
)

type User struct {
	ID  int64
	SID int64

	newborn bool

	Version       int64
	ServerVersion string

	modules map[life.ModuleKey]life.Module
}

func NewUser(id int64) *User {
	u := &User{
		ID:            id,
		ServerVersion: profile.Version(),
		modules:       make(map[life.ModuleKey]life.Module, 16),
	}

	userregister.ForEach(func(key life.ModuleKey, newFunc life.NewModuleFunc) {
		u.modules[key] = newFunc()
	})

	return u
}

func (o *User) DecodeServer(p *dbv1.UserProto) error {
	userVersion, err := checkServerVersion(p.ServerVersion)
	if err != nil {
		return err
	}

	o.ID = p.Id
	o.SID = p.Sid
	o.Version = p.Version
	o.ServerVersion = userVersion

	for key, mp := range p.Modules {
		if mod := o.modules[life.ModuleKey(key)]; mod != nil {
			if err = dbv1.DecodeUserModuleProto(mp, mod); err != nil {
				return err
			}
		}
	}

	return nil
}

func (o *User) EncodeServer(p *dbv1.UserProto, modules []life.ModuleKey) (err error) {
	p.Id = o.ID
	p.Sid = o.SID
	p.Version = o.Version
	p.ServerVersion = o.ServerVersion

	p.Modules = make(map[string]*dbv1.UserModuleProto, 16)

	for _, key := range modules {
		if mod := o.modules[key]; mod != nil {
			p.Modules[string(key)] = dbv1.EncodeUserModuleProto(mod)
		}
	}

	return nil
}

func checkServerVersion(userVersion string) (validUserVersion string, err error) {
	az1, psv, isProfileRelease := version.GetSubVersion(profile.Version())
	if !isProfileRelease {
		return userVersion, nil
	}

	az2, usv, isRelease := version.GetSubVersion(userVersion)
	if !isRelease {
		return profile.Version(), nil
	}

	if az1 != az2 {
		return "", errors.Errorf("userVersion is not equal to profileVersion. userVersion=%s profileVersion=%s", userVersion, profile.Version())
	}

	for i := range psv {
		if psv[i] < usv[i] {
			return "", errors.Errorf("userVersion is greater than profileVersion. userVersion=%s profileVersion=%s", userVersion, profile.Version())
		}

		if psv[i] > usv[i] {
			return profile.Version(), nil
		}
	}

	return userVersion, nil
}

func (o *User) EncodeClient() (*message.UserProto, error) {
	basic, err := o.EncodeClientBasic()
	if err != nil {
		return nil, err
	}

	p := &message.UserProto{
		Basic:    basic,
		Storage:  o.Storage().EncodeClient(),
		HeroList: o.HeroList().EncodeClient(),
		Room:     o.Room().EncodeClient(),
	}

	return p, nil
}

func (o *User) EncodeClientBasic() (p *message.UserBasicProto, err error) {
	p = &message.UserBasicProto{
		Name:            o.Basic().Name,
		Gender:          o.Basic().Gender,
		RechargeAmounts: o.Recharge().EncodeClient(),
	}

	p.Id, err = xid.EncodeID(o.ID)
	if err != nil {
		return nil, err
	}

	p.RechargeAmounts = o.Recharge().EncodeClient()

	return p, nil
}

func (o *User) SetNewborn(b bool) {
	o.newborn = b
}

func (o *User) GetAndResetNewborn() bool {
	b := o.newborn
	o.newborn = false

	return b
}

func (o *User) Now() time.Time {
	if profile.IsDev() {
		return time.Now().Add(o.Dev().TimeOffset())
	}

	return time.Now()
}
