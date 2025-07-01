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
	_mongoUserTableName     = "user"
	_userIDField            = "id"
	_userSIDField           = "sid"
	_userVersionField       = "version"
	_userServerVersionField = "server_version"
	_userModulesField       = "modules"
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

	// create uid unique index
	if _, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys:    bson.D{{Key: _userIDField, Value: 1}},
		Options: options.Index().SetUnique(true),
	}); err != nil {
		return nil, errors.Wrap(err, "create mongo uid index for user failed")
	}

	// create sid index
	if _, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: _userSIDField, Value: 1}},
	}); err != nil {
		return nil, errors.Wrap(err, "create mongo sid index for user failed")
	}

	// create id+version index for optimistic locking
	if _, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{
		Keys: bson.D{{Key: _userIDField, Value: 1}, {Key: _userVersionField, Value: 1}},
	}); err != nil {
		return nil, errors.Wrap(err, "create mongo version index for user failed")
	}

	return repo, nil
}

func (r *userMongoRepo) Create(ctx context.Context, user *dbv1.UserProto, ctime time.Time) error {
	if user == nil {
		return errors.New("user proto is nil")
	}

	if err := r.repo.Create(ctx, user); err != nil {
		return errors.Wrapf(err, "creating user uid=%d", user.Id)
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
		projection = append(projection,
			bson.E{Key: _userIDField, Value: 1},
			bson.E{Key: _userSIDField, Value: 1},
			bson.E{Key: _userVersionField, Value: 1},
			bson.E{Key: _userServerVersionField, Value: 1},
		)
		for _, mod := range mods {
			projection = append(projection, bson.E{Key: _userModulesField + "." + string(mod), Value: 1})
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
		_userSIDField:           user.Sid,
		_userVersionField:       user.Version,
		_userServerVersionField: user.ServerVersion,
	}

	if user.Modules != nil {
		for k, v := range user.Modules {
			updateFields[_userModulesField+"."+k] = v
		}
	}

	if err := r.repo.UpdateOne(ctx, _userIDField, uid, user.Version, updateFields); err != nil {
		return errors.Wrapf(err, "updating user %d", uid)
	}

	return nil
}

func (r *userMongoRepo) IsExist(ctx context.Context, uid int64) (bool, error) {
	return r.repo.Exists(ctx, _userIDField, uid)
}

func (r *userMongoRepo) UpdateSID(ctx context.Context, uid int64, sid int64, version int64) error {
	updateFields := bson.M{
		_userSIDField:     sid,
		_userVersionField: version,
	}

	if err := r.repo.UpdateOne(ctx, _userIDField, uid, version, updateFields); err != nil {
		return errors.Wrapf(err, "updating sid for user %d", uid)
	}

	return nil
}

func (r *userMongoRepo) IncVersion(ctx context.Context, uid int64, newVersion int64) error {
	if err := r.repo.IncrementVersion(ctx, _userIDField, uid, newVersion); err != nil {
		return errors.Wrapf(err, "incrementing version for user %d", uid)
	}

	return nil
}
