package object

import (
	"time"

	"github.com/go-pantheon/fabrica-util/data/db/postgresql"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/fabrica-util/xtime"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/userregister"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"google.golang.org/protobuf/proto"
)

const (
	ModuleKey = life.ModuleKey("basic")
)

const (
	GenderUnset  = 0
	GenderMale   = 1
	GenderFemale = 2
)

func init() {
	userregister.Register(ModuleKey, NewBasic, userregister.WithPGColumnType(postgresql.JSONB))
}

var _ life.Module = (*Basic)(nil)

type Basic struct {
	Name      string
	Gender    int32
	CreatedAt time.Time
}

func NewBasic() life.Module {
	return &Basic{
		Name:      "test",
		CreatedAt: time.Now(),
	}
}

func (o *Basic) IsLifeModule() {}

func (o *Basic) EncodeServer() proto.Message {
	p := dbv1.UserBasicProtoPool.Get()

	p.Name = o.Name
	p.Gender = o.Gender
	p.CreatedAt = o.CreatedAt.Unix()

	return p
}

func (o *Basic) DecodeServer(p proto.Message) error {
	if p == nil {
		return errors.New("basic decode server nil")
	}

	op, ok := p.(*dbv1.UserBasicProto)
	if !ok {
		return errors.Errorf("basic decode server invalid type: %T", p)
	}

	o.Name = op.Name
	o.Gender = op.Gender
	o.CreatedAt = xtime.Time(op.CreatedAt)

	return nil
}
