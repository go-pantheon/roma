package object

import (
	climsg "github.com/vulcan-frame/vulcan-game/gen/api/client/message"
	dbv1 "github.com/vulcan-frame/vulcan-game/gen/api/db/player/v1"
)

const (
	GenderUnset  = 0
	GenderMale   = 1
	GenderFemale = 2
)

type Basic struct {
	Gender   int32
	Recharge *Recharge
}

func NewBasic() *Basic {
	return &Basic{
		Recharge: NewRecharge(),
	}
}

func NewBasicProto() *dbv1.UserBasicProto {
	p := &dbv1.UserBasicProto{}
	p.Recharge = NewRechargeProto()
	return p
}

func (o *Basic) EncodeServer(p *dbv1.UserBasicProto) {
	p.Gender = o.Gender
	p.Recharge = o.Recharge.EncodeServer()
}

func (o *Basic) DecodeServer(p *dbv1.UserBasicProto) (err error) {
	if p == nil {
		return
	}
	o.Gender = p.Gender
	o.Recharge.DecodeServer(p.Recharge)
	return
}

func (o *Basic) EncodeClient() *climsg.UserBasicProto {
	p := &climsg.UserBasicProto{}
	p.Gender = o.Gender
	p.RechargeAmounts = o.Recharge.EncodeClient()
	return p
}
