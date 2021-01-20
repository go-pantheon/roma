package userobj

import (
	"fmt"
	"math/big"

	dbv1 "github.com/vulcan-frame/vulcan-game/gen/api/db/player/v1"
)

const RechargePrecision = int64(2)

type Recharge struct {
	amount big.Rat
}

func NewRecharge() *Recharge {
	return buildRecharge(0)
}

func NewRechargeFromCents(cents int64) *Recharge {
	return buildRecharge(cents)
}

func buildRecharge(cents int64) *Recharge {
	r := &Recharge{}
	num := big.NewInt(cents)
	denom := new(big.Int).Exp(big.NewInt(10), big.NewInt(RechargePrecision), nil)
	r.amount.SetFrac(num, denom)
	return r
}

func NewRechargeProto() *dbv1.RechargeProto {
	p := &dbv1.RechargeProto{}
	return p
}

func (o *Recharge) EncodeServer() *dbv1.RechargeProto {
	p := &dbv1.RechargeProto{
		Amount: o.amount.String(),
	}
	return p
}

func (o *Recharge) DecodeServer(p *dbv1.RechargeProto) (err error) {
	if p == nil {
		return
	}
	o.amount, err = amountFromString(p.Amount)
	return
}

func (o *Recharge) EncodeClient() float32 {
	amount, _ := o.amount.Float64()
	return float32(amount)
}

func (o *Recharge) AddRechargeFromString(format string) (err error) {
	var toAdd big.Rat
	if toAdd, err = amountFromString(format); err != nil {
		return err
	}
	o.amount = *new(big.Rat).Add(&o.amount, &toAdd)
	return nil
}

func (o *Recharge) AddRecharge(cents int64) (err error) {
	toAdd := buildRecharge(cents)
	o.amount = *new(big.Rat).Add(&o.amount, &toAdd.amount)
	return nil
}

func amountFromString(str string) (ret big.Rat, err error) {
	ret = big.Rat{}
	if _, ok := ret.SetString(str); !ok {
		err = fmt.Errorf("invalid amount format: %s", str)
		return
	}
	return
}
