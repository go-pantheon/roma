package object

import (
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/userregister"
	"github.com/go-pantheon/roma/gamedata"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"github.com/go-pantheon/roma/pkg/zerrors"
	"google.golang.org/protobuf/proto"
)

const (
	ModuleKey = life.ModuleKey("storage")
)

func init() {
	userregister.Register(ModuleKey, NewStorage)
}

var _ life.Module = (*Storage)(nil)

type Storage struct {
	Items         map[int64]*ItemInfo
	Packs         map[int64]*PackInfo
	RecoveryInfos map[int64]*RecoveryInfo
}

func NewStorage() life.Module {
	o := &Storage{
		Items:         make(map[int64]*ItemInfo, 128),
		Packs:         make(map[int64]*PackInfo, 128),
		RecoveryInfos: make(map[int64]*RecoveryInfo, 128),
	}

	return o
}

func (o *Storage) IsLifeModule() {}

func (o *Storage) EncodeServer() proto.Message {
	p := dbv1.UserStorageProtoPool.Get()

	p.Items = make(map[int64]uint64, len(o.Items))
	p.Packs = make(map[int64]uint64, len(o.Items))
	p.RecoveryInfos = make(map[int64]*dbv1.ItemRecoveryInfoProto, len(o.Items))

	for _, item := range o.Items {
		p.Items[item.Data().Id()] = item.Amount()
	}

	for _, pack := range o.Packs {
		p.Packs[pack.Data().Id()] = pack.Amount()
	}

	for _, info := range o.RecoveryInfos {
		pr := dbv1.ItemRecoveryInfoProtoPool.Get()
		p.RecoveryInfos[info.Id] = info.encodeServer(pr)
	}

	return p
}

func (o *Storage) DecodeServer(p proto.Message) error {
	if p == nil {
		return errors.New("storage decode server nil")
	}

	op, ok := p.(*dbv1.UserStorageProto)
	if !ok {
		return errors.Errorf("storage decode server invalid type: %T", p)
	}

	for id, count := range op.Items {
		if prize, err := gamedata.TryNewItemPrize(id, count); err != nil {
			log.Errorf("decode item info error: %+v", err)
		} else {
			o.Items[prize.Data().Id()] = NewItemInfo(prize.Data(), count)
		}
	}

	for id, count := range op.Packs {
		if info, err := gamedata.TryNewPackPrize(id, count); err != nil {
			log.Errorf("decode pack info error: %+v", err)
		} else {
			o.Packs[info.Data().Id()] = NewPackInfo(info.Data(), count)
		}
	}

	for _, pr := range op.RecoveryInfos {
		if info, err := decodeRecoveryInfo(pr, o.Items); err != nil {
			log.Errorf("decode recoverable item error: %+v", err)
		} else {
			o.RecoveryInfos[info.Id] = info
		}
	}

	return nil
}

func (o *Storage) EncodeClient() *climsg.UserStorageProto {
	p := &climsg.UserStorageProto{
		Items:         make(map[int64]uint64, len(o.Items)),
		Packs:         make(map[int64]uint64, len(o.Items)),
		RecoveryInfos: make(map[int64]*climsg.ItemRecoveryInfoProto, len(o.Items)),
	}

	for _, item := range o.Items {
		p.Items[item.Data().Id()] = item.Amount()
	}

	for _, pack := range o.Packs {
		p.Packs[pack.Data().Id()] = pack.Amount()
	}

	for _, info := range o.RecoveryInfos {
		info.EncodeClient(p)
	}

	return p
}

func (o *Storage) AddItem(d *gamedata.ResourceItemData, amount uint64) (err error) {
	if d == nil {
		err = errors.Wrapf(zerrors.ErrGameDataNotFound, "Data=ResourceItemData")
		return
	}

	if amount == 0 {
		err = errors.Wrapf(zerrors.ErrEmptyPrize, "Data=ResourceItemData id=%d, amount=%d", d.Id(), amount)
		return
	}

	if item := o.Items[d.Id()]; item != nil {
		item.Add(amount)
	} else {
		o.Items[d.Id()] = NewItemInfo(d, amount)
		if d.Type == gamedata.ItemTypeRecovery {
			o.RecoveryInfos[d.Id()] = newRecoveryInfo(d, time.Now())
		}
	}

	return nil
}

func (o *Storage) SubItem(d *gamedata.ResourceItemData, amount uint64) (err error) {
	if d == nil {
		return errors.Wrapf(zerrors.ErrGameDataNotFound, "Data=ResourceItemData")
	}

	if amount == 0 {
		return errors.Wrapf(zerrors.ErrEmptyPrize, "Data=ResourceItemData id=%d, amount=%d", d.Id(), amount)
	}

	item := o.Items[d.Id()]
	if item == nil {
		return errors.Wrapf(zerrors.ErrGameDataNotFound, "Data=ResourceItemData id=%d", d.Id())
	}

	if item.Amount() < amount {
		return errors.Wrapf(zerrors.ErrCostInsufficient, "Data=ResourceItemData id=%d, amount=%d", d.Id(), amount)
	}

	item.Sub(amount)

	if item.Amount() == 0 {
		switch item.Data().Type {
		case gamedata.ItemTypeRechargeCurrency, gamedata.ItemTypeGameCurrency, gamedata.ItemTypeRecovery:
		default:
			delete(o.Items, d.Id())
		}
	}

	return nil
}

func (o *Storage) AddPack(d *gamedata.ResourcePackData, amount uint64) (err error) {
	if d == nil {
		err = errors.Wrapf(zerrors.ErrGameDataNotFound, "Data=ResourcePackData")
		return
	}

	if amount == 0 {
		err = errors.Wrapf(zerrors.ErrEmptyPrize, "Data=ResourcePackData id=%d, amount=%d", d.Id(), amount)
		return
	}

	if pack := o.Packs[d.Id()]; pack != nil {
		pack.Add(amount)
	} else {
		o.Packs[d.Id()] = NewPackInfo(d, amount)
	}

	return
}

func (o *Storage) SubPack(d *gamedata.ResourcePackData, amount uint64) (err error) {
	if d == nil {
		err = errors.Wrapf(zerrors.ErrGameDataNotFound, "Data=ResourcePackData")
		return
	}

	if amount == 0 {
		err = errors.Wrapf(zerrors.ErrEmptyPrize, "Data=ResourcePackData id=%d, amount=%d", d.Id(), amount)
		return
	}

	pack := o.Packs[d.Id()]
	if pack == nil {
		err = errors.Wrapf(zerrors.ErrGameDataNotFound, "Data=ResourcePackData id=%d", d.Id())
		return
	}

	if pack.Amount() < amount {
		err = errors.Wrapf(zerrors.ErrCostInsufficient, "Data=ResourcePackData id=%d, amount=%d", d.Id(), amount)
		return
	}

	pack.Sub(amount)

	if pack.Amount() == 0 {
		delete(o.Packs, d.Id())
	}

	return
}
