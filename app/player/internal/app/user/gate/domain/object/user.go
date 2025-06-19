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

	Modules map[life.ModuleKey]life.Module
}

func NewUser(id int64, sid int64, serverVersion string) *User {
	u := &User{
		ID:            id,
		SID:           sid,
		ServerVersion: serverVersion,
		Modules:       make(map[life.ModuleKey]life.Module, 16),
	}

	userregister.ForEach(func(key life.ModuleKey, module life.Module) {
		u.Modules[key] = module
	})

	return u
}

func (o *User) Unmarshal(p *dbv1.UserProto) error {
	userVersion, err := checkServerVersion(o.ServerVersion, p.ServerVersion)
	if err != nil {
		return err
	}

	o.Version = p.Version
	o.ServerVersion = userVersion

	for key, data := range p.Modules {
		if mod := o.Modules[life.ModuleKey(key)]; mod != nil {
			if err := mod.Unmarshal(data); err != nil {
				return errors.WithMessagef(err, "module unmarshal failed. uid=%d, mod=%s", o.ID, key)
			}
		}
	}

	return nil
}

func (o *User) EncodeServer(p *dbv1.UserProto, modules []life.ModuleKey, all bool) (err error) {
	p.Version = o.Version
	p.ServerVersion = o.ServerVersion

	p.Modules = make(map[string][]byte, 16)

	if all {
		for key, mod := range o.Modules {
			data, err := mod.Marshal()
			if err != nil {
				return err
			}

			p.Modules[string(key)] = data
		}
	} else {
		for _, key := range modules {
			if mod := o.Modules[key]; mod != nil {
				data, err := mod.Marshal()
				if err != nil {
					return err
				}

				p.Modules[string(key)] = data
			}
		}
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
