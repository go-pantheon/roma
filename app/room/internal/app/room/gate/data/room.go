package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/room/internal/app/room/gate/domain"
	"github.com/go-pantheon/roma/app/room/internal/data"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/room/v1"
	pkmongo "github.com/go-pantheon/roma/pkg/data/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	_mongoRoomTableName = "room"
	_roomIDField        = "id"
)

var _ domain.RoomRepo = (*roomMongoRepo)(nil)

type roomMongoRepo struct {
	log  *log.Helper
	repo *pkmongo.BaseRepo
}

func NewRoomMongoRepo(data *data.Data, logger log.Logger) (domain.RoomRepo, error) {
	coll := data.Mdb.Collection(_mongoRoomTableName)
	r := &roomMongoRepo{
		repo: pkmongo.NewBaseRepo(coll),
		log:  log.NewHelper(log.With(logger, "module", "room/room/gate/data")),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: _roomIDField, Value: 1}},
		Options: options.Index().SetUnique(true),
	})

	if err != nil {
		return nil, errors.Wrap(err, "create mongo index for room failed")
	}

	return r, nil
}

func (r *roomMongoRepo) Create(ctx context.Context, p *dbv1.RoomProto, ctime time.Time) error {
	if p == nil {
		return errors.New("room proto is nil")
	}
	if p.Id == 0 {
		return errors.New("room id is required for creation")
	}
	p.CreatedAt = ctime.Unix()

	if err := r.repo.Create(ctx, p); err != nil {
		return errors.Wrapf(err, "creating room %d", p.Id)
	}
	return nil
}

func (r *roomMongoRepo) QueryByID(ctx context.Context, id int64, p *dbv1.RoomProto) error {
	if id == 0 {
		return errors.New("room id is required")
	}
	if err := r.repo.FindByID(ctx, _roomIDField, id, p, nil); err != nil {
		return errors.Wrapf(err, "querying room %d", id)
	}
	return nil
}

func (r *roomMongoRepo) UpdateByID(ctx context.Context, id int64, room *dbv1.RoomProto) error {
	if room == nil {
		return errors.New("room proto is nil")
	}
	if id == 0 {
		return errors.New("room id is required for update")
	}

	update := bson.M{
		"$set": bson.M{
			"sid":       room.Sid,
			"room_type": room.RoomType,
			"members":   room.Members,
			"version":   room.Version,
		},
	}

	if err := r.repo.UpdateOne(ctx, _roomIDField, id, room.Version, update); err != nil {
		return errors.Wrapf(err, "updating room %d", id)
	}
	return nil
}

func (r *roomMongoRepo) Exist(ctx context.Context, id int64) (bool, error) {
	if id == 0 {
		return false, errors.New("room id is required")
	}
	return r.repo.Exists(ctx, _roomIDField, id)
}

func (r *roomMongoRepo) IncVersion(ctx context.Context, id int64, newVersion int64) error {
	if id == 0 {
		return errors.New("room id is required")
	}
	if err := r.repo.IncrementVersion(ctx, _roomIDField, id, newVersion); err != nil {
		return errors.Wrapf(err, "incrementing version for room %d", id)
	}
	return nil
}
