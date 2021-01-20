package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/log"
	"github.com/pkg/errors"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/user/gate/domain"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/data"
	dbv1 "github.com/vulcan-frame/vulcan-game/gen/api/db/player/v1"
	"github.com/vulcan-frame/vulcan-kit/xerrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	_collectionName = "user"
)

var _ domain.UserRepo = (*userMongoRepo)(nil)

type userMongoRepo struct {
	log        *log.Helper
	data       *data.Data
	collection *mongo.Collection
}

func NewUserMongoRepo(data *data.Data, logger log.Logger) (domain.UserRepo, error) {
	r := &userMongoRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "player/user/gate/data")),
	}
	r.collection = data.Mdb.Collection(_collectionName)
	if _, err := r.collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.M{"name": 1},
		Options: options.Index().SetName("name_idx"),
	}); err != nil {
		return nil, errors.Wrapf(err, "mongo create index[name_idx] failed")
	}
	return r, nil
}

func (r *userMongoRepo) Create(ctx context.Context, uid int64, p *dbv1.UserProto, ctime time.Time) error {
	if _, err := r.collection.InsertOne(ctx, p); err != nil {
		return errors.Wrapf(err, "mongo insert user failed. uid=%d", uid)
	}
	return nil
}

func (r *userMongoRepo) QueryByID(ctx context.Context, uid int64, p *dbv1.UserProto) error {
	filter := bson.D{
		bson.E{Key: "_id", Value: uid},
	}
	result := r.collection.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.Wrapf(xerrors.ErrDBRecordNotFound, "mongo query user failed. uid=%d", uid)
		}
		return errors.Wrapf(err, "mongo query user failed. uid=%d", uid)
	}

	if err := result.Decode(p); err != nil {
		return errors.Wrapf(err, "mongo decode user failed. uid=%d", uid)
	}
	return nil
}

func (r *userMongoRepo) UpdateLoginTime(ctx context.Context, uid int64, loginAt, logoutAt time.Time) error {
	up := bson.D{bson.E{
		Key: "$set",
		Value: bson.M{
			"loginAt":      loginAt.Unix(),
			"lastOnlineAt": loginAt.Unix(),
			"lastLogoutAt": logoutAt.Unix(),
		}}}
	result, err := r.collection.UpdateByID(ctx, uid, up)
	if err != nil {
		return errors.Wrapf(err, "mongo update user login time failed. uid=%d", uid)
	}

	if result.ModifiedCount < 1 {
		return errors.Wrapf(xerrors.ErrDBRecordNotAffected, "mongo update user login time failed. uid=%d", uid)
	}
	return nil
}

func (r *userMongoRepo) UpdateLastOnlineTime(ctx context.Context, uid int64, t time.Time) error {
	up := bson.D{bson.E{
		Key:   "$set",
		Value: bson.M{"lastOnlineAt": t.Unix()},
	}}
	result, err := r.collection.UpdateByID(ctx, uid, up)
	if err != nil {
		return errors.Wrapf(err, "mongo update user last online time failed. uid=%d", uid)
	}

	if result.ModifiedCount < 1 {
		return errors.Wrapf(xerrors.ErrDBRecordNotAffected, "mongo update user last online time failed. uid=%d", uid)
	}
	return nil
}

func (r *userMongoRepo) UpdateByID(ctx context.Context, uid int64, user *dbv1.UserProto) error {
	filter := bson.D{
		bson.E{Key: "_id", Value: uid},
		bson.E{Key: "version", Value: user.Version - 1},
	}
	replace, err := bson.Marshal(user)
	if err != nil {
		return errors.Wrapf(err, "mongo bson encode user failed. uid=%d", uid)
	}

	result, err := r.collection.ReplaceOne(ctx, filter, replace)
	if err != nil {
		return errors.Wrapf(err, "mongo update user failed. uid=%d", uid)
	}
	if result.ModifiedCount < 1 {
		return errors.Wrapf(xerrors.ErrDBRecordNotAffected, "mongo update user failed. uid=%d", uid)
	}
	return nil
}

func (r *userMongoRepo) Exist(ctx context.Context, uid int64) (bool, error) {
	filter := bson.D{
		bson.E{Key: "_id", Value: uid},
	}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, errors.Wrapf(err, "mongo query user count failed. uid=%d", uid)
	}
	return count > 0, nil
}

func (r *userMongoRepo) IncVersion(ctx context.Context, uid int64, newVersion int64) error {
	filter := bson.D{
		bson.E{Key: "_id", Value: uid},
		bson.E{Key: "version", Value: newVersion - 1},
	}
	up := bson.D{
		bson.E{
			Key:   "$set",
			Value: bson.M{"version": newVersion},
		},
	}

	result, err := r.collection.UpdateOne(ctx, filter, up)
	if err != nil {
		return errors.Wrapf(err, "mongo increment user version failed. uid=%d", uid)
	}
	if result.ModifiedCount < 1 {
		return errors.Wrapf(xerrors.ErrDBRecordNotAffected, "mongo increment user version failed. uid=%d", uid)
	}
	return nil
}
