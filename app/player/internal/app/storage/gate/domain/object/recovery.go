package object

import (
	"time"

	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/fabrica-util/xtime"
	"github.com/go-pantheon/roma/gamedata"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/util/maths/u64"
	"github.com/pkg/errors"
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

func decodeRecoveryInfo(p *dbv1.ItemRecoveryInfoProto, items map[int64]*ItemInfo) (*RecoveryInfo, error) {
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

func (o *RecoveryInfo) encodeServer(p *dbv1.ItemRecoveryInfoProto) *dbv1.ItemRecoveryInfoProto {
	p.DataId = o.Id
	p.UpdatedAt = o.updatedAt.Unix()
	return p
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
