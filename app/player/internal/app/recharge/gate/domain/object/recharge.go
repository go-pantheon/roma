package object

import (
	"fmt"
	"math/big"

	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/userregister"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"google.golang.org/protobuf/proto"
)

const (
	ModuleKey = "recharge"
)

const (
	RechargePrecision = int64(2)
)

func init() {
	userregister.Register(ModuleKey, NewRecharge)
}

var _ life.Module = (*Recharge)(nil)

type Recharge struct {
	amount big.Rat
}

func NewRecharge() life.Module {
	o := buildRecharge(0)
	return o
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

func NewRechargeProto() *dbv1.UserRechargeProto {
	p := &dbv1.UserRechargeProto{}
	return p
}

func (o *Recharge) EncodeServer() proto.Message {
	p := dbv1.UserRechargeProtoPool.Get()
	p.Amount = o.amount.String()

	return p
}

func (o *Recharge) DecodeServer(p proto.Message) error {
	if p == nil {
		return errors.New("recharge decode server nil")
	}

	op, ok := p.(*dbv1.UserRechargeProto)
	if !ok {
		return errors.Errorf("recharge decode server invalid type: %T", p)
	}

	amount, err := amountFromString(op.Amount)
	if err != nil {
		return errors.Wrap(err, "failed to parse recharge amount")
	}

	o.amount = amount

	return nil
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
