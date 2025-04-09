package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/roma/app/room/internal/app/room/admin/domain"
	"github.com/go-pantheon/roma/app/room/internal/data"
	dbv1 "github.com/go-pantheon/roma/gen/api/server/room/admin/room/v1"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

var _ domain.RoomRepo = (*mongoRepo)(nil)

const (
	_collectionName = "room"
)

type mongoRepo struct {
	log        *log.Helper
	data       *data.Data
	collection *mongo.Collection
}

func NewMongoRepo(data *data.Data, logger log.Logger) (domain.RoomRepo, error) {
	r := &mongoRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "room/tcp/data/room")),
	}
	r.collection = data.Mdb.Collection(_collectionName)
	return r, nil
}

func (r *mongoRepo) GetByID(ctx context.Context, id int64) (*dbv1.RoomProto, error) {
	filter := bson.D{
		bson.E{Key: "_id", Value: id},
	}
	result := r.collection.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.Wrapf(xerrors.ErrDBRecordNotFound, "[mongoRepo.GetByID] 查询数据失败。room<%d>", id)
		}
		return nil, errors.Wrapf(err, "[mongoRepo.GetByID] 查询数据失败。room<%d>", id)
	}

	bo := &dbv1.RoomProto{}
	if err := result.Decode(bo); err != nil {
		return nil, errors.Wrapf(err, "[mongoRepo.GetByID] 解码失败。room<%d>", id)
	}
	return bo, nil
}
