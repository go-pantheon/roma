package biz

import (
	"strings"
	"time"
	"unicode/utf8"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-util/xtime"
	"github.com/go-pantheon/roma/app/player/internal/app/basic/gate/domain/object"
	statusobj "github.com/go-pantheon/roma/app/player/internal/app/status/gate/domain/object"
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

func NewUserUseCase(mgr *core.Manager, do *domain.UserDomain, storageDo *storagedo.StorageDomain, logger log.Logger) *UserUseCase {
	uc := &UserUseCase{
		log:       log.NewHelper(log.With(logger, "module", "player/user/gate/biz")),
		mgr:       mgr,
		do:        do,
		storageDo: storageDo,
	}

	mgr.RegisterOnCreatedEvent(uc.onCreated)
	mgr.RegisterOnLoadEvent(uc.onLoad)
	mgr.RegisterSecondTick(uc.secondTick)

	return uc
}

func (uc *UserUseCase) onCreated(ctx core.Context) error {
	ctime := ctx.Now()

	user := ctx.User()

	user.Basic().CreatedAt = ctime
	user.Status().LoginAt = ctime
	user.Status().LatestOnlineAt = ctime

	ctx.Changed(object.ModuleKey, statusobj.ModuleKey)

	user.SID = ctx.SID()

	if err := uc.do.UpdateSID(ctx, user.ID, user.SID, user.Version); err != nil {
		return err
	}

	return nil
}

func (uc *UserUseCase) onLoad(ctx core.Context) error {
	ctx.User().Status().ClientIP = ctx.ClientIP()

	return nil
}

func (uc *UserUseCase) dailyReset(ctx core.Context) (changed bool) {
	var (
		user  = ctx.User()
		ctime = ctx.Now()
	)
	if ctime.Before(user.Status().NextDailyResetAt) {
		return
	}
	user.Status().DailyOnlineDuration = 0
	user.Status().NextDailyResetAt = xtime.NextDailyTime(ctime, 0)
	changed = true
	return
}

func (uc *UserUseCase) secondTick(ctx core.Context) error {
	user := ctx.User()
	ctime := ctx.Now()

	if dur := ctime.Sub(user.Status().LatestOnlineAt); dur > 0 && dur < minOfflineDuration {
		user.Status().DailyOnlineDuration += dur
		user.Status().TotalOnlineDuration += dur
		user.Status().LatestOnlineAt = ctime
	} else {
		// set re-login if offline for a long time
		user.Status().LoginAt = ctime
		user.Status().LogoutAt = user.Status().LatestOnlineAt
		user.Status().LatestOnlineAt = ctime
	}

	ctx.Changed(statusobj.ModuleKey)

	return nil
}

func (uc *UserUseCase) Login(ctx core.Context, cs *climsg.CSLogin) (*climsg.SCLogin, error) {
	sc := &climsg.SCLogin{}
	user := ctx.User()

	uc.dailyReset(ctx)

	userProto, err := user.EncodeClient()
	if err != nil {
		return nil, err
	}

	sc.Code = climsg.SCLogin_Succeeded
	sc.User = userProto
	sc.ServerTime = time.Now().Unix()
	sc.Newborn = user.GetAndResetNewborn()

	return sc, nil
}

func (uc *UserUseCase) UpdateName(ctx core.Context, cs *climsg.CSUpdateName) (*climsg.SCUpdateName, error) {
	sc := &climsg.SCUpdateName{}

	name := strings.TrimSpace(cs.Name)

	if size := utf8.RuneCountInString(name); size < nameMinLength || size > nameMaxLength {
		sc.Code = climsg.SCUpdateName_ErrNameIllegal
		return sc, nil
	}

	basic := ctx.User().Basic()

	if name == basic.Name {
		sc.Code = climsg.SCUpdateName_ErrNameNotChanged
		return sc, nil
	}

	basic.Name = name

	ctx.Changed(object.ModuleKey)

	sc.Code = climsg.SCUpdateName_Succeeded

	return sc, nil
}

func (uc *UserUseCase) SetGender(ctx core.Context, cs *climsg.CSSetGender) (*climsg.SCSetGender, error) {
	sc := &climsg.SCSetGender{}

	if cs.Gender != object.GenderMale && cs.Gender != object.GenderFemale {
		sc.Code = climsg.SCSetGender_ErrGenderIllegal
		return sc, nil
	}

	basic := ctx.User().Basic()
	if basic.Gender != object.GenderUnset {
		sc.Code = climsg.SCSetGender_ErrGenderSet
		return sc, nil
	}

	basic.Gender = cs.Gender

	ctx.Changed(object.ModuleKey)

	sc.Code = climsg.SCSetGender_Succeeded

	return sc, nil
}
