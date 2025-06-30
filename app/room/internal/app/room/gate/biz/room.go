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

	mgr.RegisterOnCreatedEvent(uc.onCreated)
	mgr.RegisterMinuteTick(uc.onMinuteTick)
	mgr.RegisterSecondTick(uc.onSecondTick)

	return uc
}

func (uc *RoomUseCase) onCreated(ctx core.Context) error {
	uc.log.WithContext(ctx).Debugf("room is created. id=%d", ctx.OID())
	return nil
}

func (uc *RoomUseCase) onMinuteTick(ctx core.Context) error {
	uc.log.WithContext(ctx).Debugf("room on minute tick. id=%d", ctx.OID())
	return nil
}

func (uc *RoomUseCase) onSecondTick(ctx core.Context) error {
	uc.log.WithContext(ctx).Debugf("room on second tick. id=%d", ctx.OID())
	return nil
}

func (uc *RoomUseCase) RoomList(ctx core.Context, cs *climsg.CSRoomList) (*climsg.SCRoomList, error) {
	sc := &climsg.SCRoomList{}

	return sc, nil
}

func (uc *RoomUseCase) CreateRoom(ctx core.Context, cs *climsg.CSCreateRoom) (*climsg.SCCreateRoom, error) {
	sc := &climsg.SCCreateRoom{}

	room := ctx.Room()

	sc.Room = room.EncodeClient()
	sc.Code = climsg.SCCreateRoom_Succeeded

	uc.updateRoomInfo(ctx, room)

	ctx.Changed()

	return sc, nil
}

func (uc *RoomUseCase) updateRoomInfo(c core.Context, room *object.Room) {
	//	TODO sync data with player service
}
