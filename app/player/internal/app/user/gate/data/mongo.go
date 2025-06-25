package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/data/mongodb"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	_mongoUserTableName = "user"
	_userIDField        = "id"
	_userSIDField       = "sid"
)

var _ domain.UserRepo = (*userMongoRepo)(nil)

type userMongoRepo struct {
	log  *log.Helper
	repo *mongodb.BaseRepo
}

func NewUserMongoRepo(data *mongodb.DB, logger log.Logger) (r domain.UserRepo, err error) {
	coll := data.DB.Collection(_mongoUserTableName)
	repo := &userMongoRepo{
		repo: mongodb.NewBaseRepo(coll),
		log:  log.NewHelper(log.With(logger, "module", "player/user/gate/data/mongo")),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err = coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: _userIDField, Value: 1}},
		Options: options.Index().SetUnique(true),
	})
	if err != nil {
		return nil, errors.Wrap(err, "create mongo uid index for user failed")
	}

	_, err = coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: _userSIDField, Value: 1}},
	})
	if err != nil {
		return nil, errors.Wrap(err, "create mongo sid index for user failed")
	}

	return repo, nil
}

func (r *userMongoRepo) Create(ctx context.Context, uid int64, defaultUser *dbv1.UserProto, ctime time.Time) error {
	if defaultUser == nil {
		return errors.New("user proto is nil")
	}
	defaultUser.Id = uid

	if err := r.repo.Create(ctx, defaultUser); err != nil {
		return errors.Wrapf(err, "inserting user %d", uid)
	}
	return nil
}

func (r *userMongoRepo) QueryByID(ctx context.Context, uid int64, p *dbv1.UserProto, mods []life.ModuleKey) error {
	if p == nil {
		return errors.New("user proto is nil")
	}

	// Specific projection logic for user repo
	projection := bson.D{}
	if len(mods) > 0 {
		projection = append(projection, bson.E{Key: "id", Value: 1}, bson.E{Key: "version", Value: 1}, bson.E{Key: "server_version", Value: 1})
		for _, mod := range mods {
			projection = append(projection, bson.E{Key: "modules." + string(mod), Value: 1})
		}
	}

	if err := r.repo.FindByID(ctx, _userIDField, uid, p, projection); err != nil {
		return errors.Wrapf(err, "querying user %d", uid)
	}
	return nil
}

func (r *userMongoRepo) UpdateByID(ctx context.Context, uid int64, user *dbv1.UserProto) error {
	if user == nil {
		return errors.New("user proto is nil")
	}

	updateFields := bson.M{
		"version": user.Version,
	}
	if user.Modules != nil {
		for k, v := range user.Modules {
			updateFields["modules."+k] = v
		}
	}
	update := bson.M{"$set": updateFields}

	if err := r.repo.UpdateOne(ctx, _userIDField, uid, user.Version, update); err != nil {
		return errors.Wrapf(err, "updating user %d", uid)
	}

	return nil
}

func (r *userMongoRepo) IsExist(ctx context.Context, uid int64) (bool, error) {
	return r.repo.Exists(ctx, _userIDField, uid)
}

func (r *userMongoRepo) IncVersion(ctx context.Context, uid int64, newVersion int64) error {
	if err := r.repo.IncrementVersion(ctx, _userIDField, uid, newVersion); err != nil {
		return errors.Wrapf(err, "incrementing version for user %d", uid)
	}
	return nil
}
