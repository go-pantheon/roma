package data

import (
	"context"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/fabrica-util/data/db"
	"github.com/go-pantheon/roma/app/room/internal/app/room/gate/domain"
	"github.com/go-pantheon/roma/app/room/internal/data"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/room/v1"
	"github.com/go-pantheon/roma/pkg/util/maths/i64"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.mongodb.org/mongo-driver/v2/mongo/writeconcern"
)

const (
	_collectionName   = "room"
	_idCollectionName = "room_ids"
	batchInc          = 100
)

var _ domain.RoomRepo = (*mongoRepo)(nil)

type mongoRepo struct {
	log          *log.Helper
	data         *data.Data
	collection   *mongo.Collection
	idCollection *mongo.Collection
	idGen        atomic.Int64
	maxId        int64

	once sync.Once
}

func NewMongoRepo(data *data.Data, logger log.Logger) (domain.RoomRepo, error) {
	r := &mongoRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "room/room/gate/data/room")),
	}
	r.collection = data.Mdb.Collection(_collectionName)
	r.idCollection = data.Mdb.Collection(_idCollectionName, options.Collection().SetWriteConcern(writeconcern.Majority()))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	if err := db.InitIncrementIDDoc(ctx, r.idCollection, _collectionName); err != nil {
		return nil, err
	}
	return r, nil
}

func (r *mongoRepo) Create(ctx context.Context, p *dbv1.RoomProto, ctime time.Time) (err error) {
	id, err := r.IncrementID(ctx)
	if err != nil {
		return
	}

	p.Id = id
	if _, err = r.collection.InsertOne(ctx, p); err != nil {
		err = errors.Wrapf(err, "mongo insert room failed. id=%d", p.Id)
		return
	}
	return
}

func (r *mongoRepo) QueryByID(ctx context.Context, id int64, p *dbv1.RoomProto) error {
	filter := bson.D{
		bson.E{Key: "_id", Value: id},
	}
	result := r.collection.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.Wrapf(xerrors.ErrDBRecordNotFound, "mongo query room failed. id=%d", id)
		}
		return errors.Wrapf(err, "mongo query room failed. id=%d", id)
	}

	if err := result.Decode(p); err != nil {
		return errors.Wrapf(err, "mongo decode room failed. id=%d", id)
	}
	return nil
}

func (r *mongoRepo) UpdateByID(ctx context.Context, id int64, proto *dbv1.RoomProto) error {
	lastVersion := proto.Version - 1

	filter := bson.D{
		bson.E{Key: "_id", Value: id},
		bson.E{Key: "version", Value: lastVersion},
	}
	replace, err := bson.Marshal(proto)
	if err != nil {
		return errors.Wrapf(err, "mongo bson encode room failed. id=%d", id)
	}

	result, err := r.collection.ReplaceOne(ctx, filter, replace)
	if err != nil {
		return errors.Wrapf(err, "mongo update room failed. id=%d", id)
	}
	if result.ModifiedCount < 1 {
		return errors.Wrapf(xerrors.ErrDBRecordNotAffected, "mongo update room failed. id=%d version=%d", id, lastVersion)
	}
	return nil
}

func (r *mongoRepo) Exist(ctx context.Context, id int64) (bool, error) {
	filter := bson.D{
		bson.E{Key: "_id", Value: id},
	}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, errors.Wrapf(err, "mongo query room count failed. id=%d", id)
	}
	return count > 0, nil
}

func (r *mongoRepo) IncVersion(ctx context.Context, id int64, newVersion int64) error {
	lastVersion := newVersion - 1

	filter := bson.D{
		bson.E{Key: "_id", Value: id},
		bson.E{Key: "version", Value: lastVersion},
	}
	up := bson.D{
		bson.E{
			Key:   "$set",
			Value: bson.M{"version": newVersion},
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, up)
	if err != nil {
		return errors.Wrapf(err, "mongo increment room version failed. id=%d", id)
	}
	if result.ModifiedCount < 1 {
		return errors.Wrapf(xerrors.ErrDBRecordNotAffected, "mongo increment room version failed. id=%d version=%d", id, lastVersion)
	}
	return nil
}

func (r *mongoRepo) IncrementID(ctx context.Context) (int64, error) {
	var err error
	r.once.Do(func() {
		err = r.resetIncrementIdGen(ctx)
	})
	if err != nil {
		return 0, err
	}

	id := r.idGen.Add(1)
	if id > r.maxId {
		if err := r.resetIncrementIdGen(ctx); err != nil {
			return 0, errors.WithMessagef(err, "mongo increment id failed. id=%d maxId=%d", id, r.maxId)
		}
		id = r.idGen.Add(1)
	}
	return id, nil
}

func (r *mongoRepo) resetIncrementIdGen(ctx context.Context) (err error) {
	r.maxId, err = db.IncrementBatchID(ctx, r.idCollection, _collectionName, batchInc)
	if err != nil {
		return
	}
	r.idGen.Store(i64.Max(0, r.maxId-batchInc))
	return
}
