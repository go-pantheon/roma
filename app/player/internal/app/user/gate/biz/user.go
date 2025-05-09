package biz

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-util/xtime"
	"github.com/go-pantheon/roma/app/player/internal/app/basic/gate/domain/object"
	storagedo "github.com/go-pantheon/roma/app/player/internal/app/storage/gate/domain"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain"
	"github.com/go-pantheon/roma/app/player/internal/core"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
)

const (
	minOfflineDuration = time.Minute * 5
	nameMinLength      = 1
	nameMaxLength      = 12
)

type UserUseCase struct {
	log       *log.Helper
	mgr       *core.Manager
	do        *domain.UserDomain
	storageDo *storagedo.StorageDomain
}

func NewUserUseCase(mgr *core.Manager,
	do *domain.UserDomain,
	storageDo *storagedo.StorageDomain,
	logger log.Logger) *UserUseCase {
	uc := &UserUseCase{
		log:       log.NewHelper(log.With(logger, "module", "player/user/gate/biz")),
		mgr:       mgr,
		do:        do,
		storageDo: storageDo,
	}

	mgr.OnCreatedEventRegister(uc.onCreated)
	mgr.OnLoadEventRegister(uc.onLoad)
	mgr.SecondTickRegister(uc.secondTick)

	return uc
}

func (uc *UserUseCase) onCreated(ctx core.Context) {
	ctime := ctx.Now()
	ctx.User().CreatedAt = ctime
	ctx.User().LoginAt = ctime
	ctx.User().LastOnlineAt = ctime

	ctx.Changed()
}

func (uc *UserUseCase) onLoad(ctx core.Context) {
	ctx.User().LastOnlineIP = ctx.ClientIP()
}

func (uc *UserUseCase) dailyReset(ctx core.Context) (changed bool) {
	var (
		user  = ctx.User()
		ctime = ctx.Now()
	)
	if ctime.Before(user.NextDailyResetAt) {
		return
	}
	user.DailyOnlineDuration = 0
	user.NextDailyResetAt = xtime.NextDailyTime(ctime, 0)
	changed = true
	return
}

func (uc *UserUseCase) secondTick(ctx core.Context) {
	user := ctx.User()
	ctime := ctx.Now()

	if dur := ctime.Sub(user.LastOnlineAt); dur > 0 && dur < minOfflineDuration {
		user.DailyOnlineDuration += dur
		user.TotalOnlineDuration += dur
		user.LastOnlineAt = ctime
	} else {
		// set re-login if offline for a long time
		user.LoginAt = ctime
		user.LogoutAt = user.LastOnlineAt
		user.LastOnlineAt = ctime
		ctx.Changed()
	}
}

func (uc *UserUseCase) Login(ctx core.Context, cs *climsg.CSLogin) (sc *climsg.SCLogin, err error) {
	sc = &climsg.SCLogin{}
	user := ctx.User()

	uc.dailyReset(ctx)

	sc.Code = climsg.SCLogin_Succeeded
	sc.User = user.EncodeClient()
	sc.ServerTime = time.Now().Unix()
	sc.Newborn = user.GetAndResetNewborn()

	return
}

func (uc *UserUseCase) UpdateName(ctx core.Context, cs *climsg.CSUpdateName) (sc *climsg.SCUpdateName, err error) {
	sc = &climsg.SCUpdateName{}

	name := strings.TrimSpace(cs.Name)

	if size := utf8.RuneCountInString(name); size < nameMinLength || size > nameMaxLength {
		sc.Code = climsg.SCUpdateName_ErrNameIllegal
		return
	}

	if name == ctx.User().Name {
		sc.Code = climsg.SCUpdateName_ErrNameNotChanged
		return
	}

	ctx.User().Name = name
	ctx.Changed()

	sc.Code = climsg.SCUpdateName_Succeeded
	return
}

func (uc *UserUseCase) SetGender(ctx core.Context, cs *climsg.CSSetGender) (sc *climsg.SCSetGender, err error) {
	sc = &climsg.SCSetGender{}

	if cs.Gender != object.GenderMale && cs.Gender != object.GenderFemale {
		sc.Code = climsg.SCSetGender_ErrGenderIllegal
		return
	}

	o := ctx.User().Basic
	if o.Gender != object.GenderUnset {
		sc.Code = climsg.SCSetGender_ErrGenderSet
		return
	}

	o.Gender = cs.Gender
	ctx.Changed()

	sc.Code = climsg.SCSetGender_Succeeded
	return
}
