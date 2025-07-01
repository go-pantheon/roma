package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/room/internal/app/room/admin/domain"
	adminv1 "github.com/go-pantheon/roma/gen/api/server/room/admin/room/v1"
	"github.com/go-pantheon/roma/pkg/data/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

const (
	_mongoRoomTableName = "room"
)

var _ domain.RoomRepo = (*roomMongoRepo)(nil)

type roomMongoRepo struct {
	log  *log.Helper
	data *mongodb.DB
	coll *mongo.Collection
}

func NewRoomMongoRepo(data *mongodb.DB, logger log.Logger) (domain.RoomRepo, error) {
	r := &roomMongoRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "room/admin/data")),
		coll: data.DB.Collection(_mongoRoomTableName),
	}

	return r, nil
}

func (r *roomMongoRepo) GetByID(ctx context.Context, id int64) (*adminv1.RoomProto, error) {
	if id == 0 {
		return nil, errors.New("room id is required")
	}

	filter := bson.M{"_id": id}

	var room adminv1.RoomProto

	err := r.coll.FindOne(ctx, filter).Decode(&room)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errors.Wrapf(err, "room %d not found", id)
		}

		return nil, errors.Wrapf(err, "querying room %d", id)
	}

	return &room, nil
}
