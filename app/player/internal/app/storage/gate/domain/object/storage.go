package object

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/roma/gamedata"
	climsg "github.com/go-pantheon/roma/gen/api/client/message"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/errs"
	"github.com/pkg/errors"
)

type Storage struct {
	Items         map[int64]*ItemInfo
	Packs         map[int64]*PackInfo
	RecoveryInfos map[int64]*RecoveryInfo
}

func NewStorage() *Storage {
	o := &Storage{
		Items: make(map[int64]*ItemInfo, 128),
		Packs: make(map[int64]*PackInfo, 128),
	}
	return o
}

func NewStorageProto() *dbv1.UserStorageProto {
	p := &dbv1.UserStorageProto{
		Items:         make(map[int64]uint64, 128),
		Packs:         make(map[int64]uint64, 128),
		RecoveryInfos: make(map[int64]*dbv1.ItemRecoveryInfoProto, 128),
	}
	return p
}

func (o *Storage) EncodeServer(p *dbv1.UserStorageProto) {
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
		p.RecoveryInfos[info.Id] = NewRecoveryInfoProto()
		info.EncodeServer(p.RecoveryInfos[info.Id])
	}
}

func (o *Storage) DecodeServer(ctx context.Context, p *dbv1.UserStorageProto) {
	if p == nil {
		return
	}

	for id, count := range p.Items {
		if prize, err := gamedata.TryNewItemPrize(id, count); err != nil {
			log.Errorf("decode item info error: %+v", err)
		} else {
			o.Items[prize.Data().Id()] = NewItemInfo(prize.Data(), count)
		}
	}

	for id, count := range p.Packs {
		if info, err := gamedata.TryNewPackPrize(id, count); err != nil {
			log.Errorf("decode pack info error: %+v", err)
		} else {
			o.Packs[info.Data().Id()] = NewPackInfo(info.Data(), count)
		}
	}

	for _, pr := range p.RecoveryInfos {
		if info, err := DecodeRecoveryInfo(pr, o.Items); err != nil {
			log.Errorf("decode recoverable item error: %+v", err)
		} else {
			o.RecoveryInfos[info.Id] = info
		}
	}
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
		err = errors.Wrapf(errs.ErrGameDataNotFound, "Data=ResourceItemData")
		return
	}

	if amount == 0 {
		err = errors.Wrapf(errs.ErrEmptyPrize, "Data=ResourceItemData id=%d, amount=%d", d.Id(), amount)
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
	return
}

func (o *Storage) SubItem(d *gamedata.ResourceItemData, amount uint64) (err error) {
	if d == nil {
		err = errors.Wrapf(errs.ErrGameDataNotFound, "Data=ResourceItemData")
		return
	}

	if amount == 0 {
		err = errors.Wrapf(errs.ErrEmptyPrize, "Data=ResourceItemData id=%d, amount=%d", d.Id(), amount)
		return
	}

	item := o.Items[d.Id()]
	if item == nil {
		err = errors.Wrapf(errs.ErrGameDataNotFound, "Data=ResourceItemData id=%d", d.Id())
		return
	}
	if item.Amount() < amount {
		err = errors.Wrapf(errs.ErrCostInsufficient, "Data=ResourceItemData id=%d, amount=%d", d.Id(), amount)
		return
	}

	item.Sub(amount)

	if item.Amount() == 0 {
		switch item.Data().Type {
		case gamedata.ItemTypeRechargeCurrency, gamedata.ItemTypeGameCurrency, gamedata.ItemTypeRecovery:
		default:
			delete(o.Items, d.Id())
		}
	}
	return
}

func (o *Storage) AddPack(d *gamedata.ResourcePackData, amount uint64) (err error) {
	if d == nil {
		err = errors.Wrapf(errs.ErrGameDataNotFound, "Data=ResourcePackData")
		return
	}

	if amount == 0 {
		err = errors.Wrapf(errs.ErrEmptyPrize, "Data=ResourcePackData id=%d, amount=%d", d.Id(), amount)
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
		err = errors.Wrapf(errs.ErrGameDataNotFound, "Data=ResourcePackData")
		return
	}

	if amount == 0 {
		err = errors.Wrapf(errs.ErrEmptyPrize, "Data=ResourcePackData id=%d, amount=%d", d.Id(), amount)
		return
	}

	pack := o.Packs[d.Id()]
	if pack == nil {
		err = errors.Wrapf(errs.ErrGameDataNotFound, "Data=ResourcePackData id=%d", d.Id())
		return
	}
	if pack.Amount() < amount {
		err = errors.Wrapf(errs.ErrCostInsufficient, "Data=ResourcePackData id=%d, amount=%d", d.Id(), amount)
		return
	}

	pack.Sub(amount)

	if pack.Amount() == 0 {
		delete(o.Packs, d.Id())
	}
	return
}
