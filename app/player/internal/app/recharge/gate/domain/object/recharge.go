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

var _ life.Module = (*Recharge)(nil)

type Recharge struct {
	amount big.Rat
}

func NewRecharge() *Recharge {
	o := buildRecharge(0)
	o.Register()
	return o
}

func (o *Recharge) Register() {
	userregister.Register(ModuleKey, o)
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

func (o *Recharge) Marshal() ([]byte, error) {
	p := dbv1.RechargeProtoPool.Get()
	defer dbv1.RechargeProtoPool.Put(p)

	p.Amount = o.amount.String()
	return proto.Marshal(p)
}

func (o *Recharge) Unmarshal(data []byte) error {
	p := dbv1.RechargeProtoPool.Get()
	defer dbv1.RechargeProtoPool.Put(p)

	if err := proto.Unmarshal(data, p); err != nil {
		return errors.Wrap(err, "failed to unmarshal recharge")
	}

	amount, err := amountFromString(p.Amount)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal recharge")
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
