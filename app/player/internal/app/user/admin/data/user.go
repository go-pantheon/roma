package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/roma/app/player/internal/app/user/admin/domain"
	"github.com/go-pantheon/roma/app/player/internal/data"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
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
		log:  log.NewHelper(log.With(logger, "module", "player/user/admin/data")),
	}
	r.collection = data.Mdb.Collection("user")
	return r, nil
}

func (r *userMongoRepo) GetByID(ctx context.Context, uid int64) (*dbv1.UserProto, error) {
	filter := bson.D{
		bson.E{Key: "_id", Value: uid},
	}
	result := r.collection.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, xerrors.APINotFound("query user failed. uid=%d", uid)
		}
		return nil, xerrors.APIDBFailed("query user failed. uid=%d", uid).WithCause(err)
	}

	bo := &dbv1.UserProto{}
	if err := result.Decode(bo); err != nil {
		return nil, xerrors.APICodecFailed("decode user failed. uid=%d", uid).WithCause(err)
	}
	return bo, nil
}

var userListFields = bson.D{
	bson.E{Key: "name", Value: 1},
	bson.E{Key: "createdAt", Value: 1},
	bson.E{Key: "loginAt", Value: 1},
	bson.E{Key: "lastOnlineAt", Value: 1},
	bson.E{Key: "lastOnlineIp", Value: 1},
}

func (r *userMongoRepo) GetList(ctx context.Context, start, limit int64, cond *dbv1.UserProto) (result []*dbv1.UserProto, count int64, err error) {
	filter := r.buildFilter(ctx, cond)

	count, err = r.collection.CountDocuments(ctx, filter)
	if err != nil {
		err = xerrors.APIDBFailed("query user count failed").WithCause(err)
		return
	}

	opts := options.Find().SetSort(bson.D{bson.E{Key: "lastOnlineAt", Value: -1}}).SetSkip(start).SetLimit(limit)
	opts = opts.SetProjection(userListFields)
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			result = make([]*dbv1.UserProto, 0)
			return
		}
		err = xerrors.APIDBFailed("query user list failed").WithCause(err)
		return
	}

	users := make([]dbv1.UserProto, 0, limit)
	if err = cursor.All(ctx, &users); err != nil {
		err = xerrors.APIDBFailed("create user list failed").WithCause(err)
		return
	}

	result = make([]*dbv1.UserProto, 0, len(users))
	for i := 0; i < len(users); i++ {
		result = append(result, &users[i])
	}
	return
}

func (r *userMongoRepo) UpdateByID(ctx context.Context, uid int64, user *dbv1.UserProto) error {
	filter := bson.D{
		bson.E{Key: "_id", Value: uid},
	}
	replace, err := bson.Marshal(user)
	if err != nil {
		return xerrors.APICodecFailed("encode user bson failed. uid=%d", uid).WithCause(err)
	}

	result, err := r.collection.ReplaceOne(ctx, filter, replace)
	if err != nil {
		return xerrors.APIDBFailed("update user failed. uid=%d", uid).WithCause(err)
	}
	if result.ModifiedCount < 1 {
		return xerrors.APIDBNoAffected("update user failed. uid=%d version=%d", uid, user.Version)
	}
	return nil
}

func (r *userMongoRepo) buildFilter(_ context.Context, cond *dbv1.UserProto) bson.D {
	filter := make(bson.D, 0, 8)
	filter = append(filter, bson.E{Key: "_id", Value: bson.D{bson.E{Key: "$gt", Value: 0}}})
	if len(cond.Name) > 0 {
		filter = append(filter, bson.E{Key: "name", Value: bson.D{bson.E{Key: "$regex", Value: bson.Regex{
			Pattern: cond.Name,
		}}}})
	}

	return filter
}
