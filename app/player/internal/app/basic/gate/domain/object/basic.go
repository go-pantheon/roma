package object

import (
	"time"

	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/userregister"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"google.golang.org/protobuf/proto"
)

const (
	ModuleKey = "basic"
)

const (
	GenderUnset  = 0
	GenderMale   = 1
	GenderFemale = 2
)

var _ life.Module = (*Basic)(nil)

type Basic struct {
	Name      string
	Gender    int32
	CreatedAt time.Time
}

func NewBasic() *Basic {
	o := &Basic{}
	o.Register()
	return o
}

func (o *Basic) Register() {
	userregister.Register(ModuleKey, o)
}

func (o *Basic) Marshal() ([]byte, error) {
	p := dbv1.UserBasicProtoPool.Get()
	defer dbv1.UserBasicProtoPool.Put(p)

	p.Name = o.Name
	p.Gender = o.Gender
	p.CreatedAt = o.CreatedAt.Unix()

	return proto.Marshal(p)
}

func (o *Basic) Unmarshal(bytes []byte) (err error) {
	p := dbv1.UserBasicProtoPool.Get()
	defer dbv1.UserBasicProtoPool.Put(p)

	if err = proto.Unmarshal(bytes, p); err != nil {
		return err
	}

	o.Name = p.Name
	o.Gender = p.Gender
	o.CreatedAt = time.Unix(p.CreatedAt, 0)

	return
}
