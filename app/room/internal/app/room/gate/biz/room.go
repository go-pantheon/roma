package biz

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/roma/app/room/internal/app/room/gate/domain"
	"github.com/go-pantheon/roma/app/room/internal/app/room/gate/domain/object"
	"github.com/go-pantheon/roma/app/room/internal/core"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
)

type RoomUseCase struct {
	log *log.Helper
	mgr *core.Manager
	do  *domain.RoomDomain
}

func NewRoomUseCase(mgr *core.Manager,
	do *domain.RoomDomain,
	logger log.Logger,
) *RoomUseCase {
	uc := &RoomUseCase{
		log: log.NewHelper(log.With(logger, "module", "room/room/gate/biz")),
		mgr: mgr,
		do:  do,
	}

	mgr.OnCreatedEventRegister(uc.onCreated)
	mgr.MinuteTickRegister(uc.onMinuteTick)
	mgr.SecondTickRegister(uc.onSecondTick)
	return uc
}

func (uc *RoomUseCase) onCreated(ctx core.Context) {
	uc.log.WithContext(ctx).Debugf("room is created. id=%d", ctx.OID())
}

func (uc *RoomUseCase) onMinuteTick(ctx core.Context) {
	uc.log.WithContext(ctx).Debugf("room on minute tick. id=%d", ctx.OID())
}

func (uc *RoomUseCase) onSecondTick(ctx core.Context) {
	uc.log.WithContext(ctx).Debugf("room on second tick. id=%d", ctx.OID())
}

func (uc *RoomUseCase) RoomList(ctx core.Context, cs *climsg.CSRoomList) (sc *climsg.SCRoomList, err error) {
	sc = &climsg.SCRoomList{}
	return
}

func (uc *RoomUseCase) CreateRoom(ctx core.Context, cs *climsg.CSCreateRoom) (sc *climsg.SCCreateRoom, err error) {
	sc = &climsg.SCCreateRoom{}

	room := ctx.Room()

	sc.Room = room.EncodeClient()
	sc.Code = climsg.SCCreateRoom_Succeeded

	uc.updateRoomInfo(ctx, room)
	return
}

func (uc *RoomUseCase) updateRoomInfo(c core.Context, room *object.Room) {
	//	TODO sync data with player service
}
