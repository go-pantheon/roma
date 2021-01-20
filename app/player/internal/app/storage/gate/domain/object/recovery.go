package object

import (
	"time"

	"github.com/pkg/errors"
	"github.com/vulcan-frame/vulcan-game/gamedata"
	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
	dbv1 "github.com/vulcan-frame/vulcan-game/gen/api/db/player/v1"
	"github.com/vulcan-frame/vulcan-game/pkg/util/maths/u64"
	"github.com/vulcan-frame/vulcan-kit/xerrors"
	"github.com/vulcan-frame/vulcan-util/xtime"
)

// RecoveryInfo Auto Recoverable Item Info
type RecoveryInfo struct {
	Id            int64
	max           uint64
	updatedAt     time.Time // last updated time
	recoverPerSec float64   // recover per second
}

func newRecoveryInfo(d *gamedata.ResourceItemData, ctime time.Time) *RecoveryInfo {
	item := &RecoveryInfo{
		Id:            d.Id(),
		recoverPerSec: d.RecoverPerSec,
		updatedAt:     ctime,
		max:           d.Max(),
	}
	return item
}

func NewRecoveryInfoProto() *dbv1.ItemRecoveryInfoProto {
	p := &dbv1.ItemRecoveryInfoProto{}
	return p
}

func DecodeRecoveryInfo(p *dbv1.ItemRecoveryInfoProto, items map[int64]*ItemInfo) (*RecoveryInfo, error) {
	if p == nil {
		return nil, errors.Wrapf(xerrors.ErrDBProtoDecode, "ItemRecoveryInfoProto is nil")
	}

	item := items[p.DataId]
	if item == nil {
		return nil, errors.Wrapf(xerrors.ErrDBProtoDecode, "item not found. dataId=%d", p.DataId)
	}

	o := newRecoveryInfo(item.Data(), xtime.Time(p.UpdatedAt))
	return o, nil
}

// Recover calculate the new count based on the current time and the recover per second
func (o *RecoveryInfo) Recover(ctime time.Time) (toAdd uint64) {
	if o.updatedAt.IsZero() {
		// something wrong, set the updated time and wait for next recover
		o.updatedAt = ctime
		return
	}

	seconds := int64(ctime.Sub(o.updatedAt).Seconds())
	toAdd = uint64(float64(seconds) * o.recoverPerSec)
	if toAdd <= 0 {
		return
	}
	toAdd = u64.Min(toAdd, o.max)
	o.updatedAt = ctime
	return
}

func (o *RecoveryInfo) EncodeServer(p *dbv1.ItemRecoveryInfoProto) {
	p.DataId = o.Id
	p.UpdatedAt = o.updatedAt.Unix()
}

func (o *RecoveryInfo) EncodeClient(storage *climsg.UserStorageProto) {
	p := &climsg.ItemRecoveryInfoProto{
		DataId:        o.Id,
		Max:           int64(o.max),
		RecoverPerSec: float32(o.recoverPerSec),
		UpdatedAt:     o.updatedAt.Unix(),
	}
	storage.RecoveryInfos[o.Id] = p
}
