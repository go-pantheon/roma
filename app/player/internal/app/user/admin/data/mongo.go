package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/player/internal/app/user/admin/domain"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	_mongoUserTableName = "user"
)

var _ domain.UserRepo = (*userMongoRepo)(nil)

type userMongoRepo struct {
	log  *log.Helper
	data *mongo.Database
	coll *mongo.Collection
}

func NewUserMongoRepo(data *mongo.Database, logger log.Logger) domain.UserRepo {
	return &userMongoRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "player/user/admin/data")),
		coll: data.Collection(_mongoUserTableName),
	}
}

func (r *userMongoRepo) GetByID(ctx context.Context, user *dbv1.UserProto, mods []life.ModuleKey) error {
	if user.Id == 0 {
		return errors.New("user id is required")
	}

	filter := bson.M{"id": user.Id}
	projection := bson.D{}

	if len(mods) == 0 {
		projection = append(projection, bson.E{Key: "_id", Value: 0})
	} else {
		projection = append(projection, bson.E{Key: "id", Value: 1})
		for _, mod := range mods {
			projection = append(projection, bson.E{Key: "modules." + string(mod), Value: 1})
		}
	}

	opts := options.FindOne().SetProjection(projection)

	err := r.coll.FindOne(ctx, filter, opts).Decode(user)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return xerrors.ErrDBRecordNotFound
		}

		return errors.Wrapf(err, "querying user %d", user.Id)
	}

	return nil
}

func (r *userMongoRepo) GetList(ctx context.Context, start, limit int64, conds map[life.ModuleKey]*dbv1.UserModuleProto, mods []life.ModuleKey) (ret []*dbv1.UserProto, total int64, err error) {
	filter := bson.M{}
	for modKey, modProto := range conds {
		// This logic assumes we are filtering based on fields within a module.
		// For simplicity, this example just filters by the whole module object.
		// A real implementation might need more complex logic to build filters from inside the proto.
		filter["modules."+string(modKey)] = modProto
	}

	total, err = r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, errors.Wrap(err, "counting users failed")
	}

	if total == 0 {
		return []*dbv1.UserProto{}, 0, nil
	}

	// Use the builder pattern for options
	findOptions := options.Find().SetSkip(start).SetLimit(limit)

	if len(mods) > 0 {
		projection := bson.D{bson.E{Key: "id", Value: 1}} // Always include ID

		for _, mod := range mods {
			projection = append(projection, bson.E{Key: "modules." + string(mod), Value: 1})
		}

		findOptions.SetProjection(projection)
	}

	cursor, err := r.coll.Find(ctx, filter, findOptions)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return []*dbv1.UserProto{}, 0, nil
		}

		return nil, 0, errors.Wrap(err, "finding users failed")
	}

	defer func() {
		if closeErr := cursor.Close(ctx); closeErr != nil {
			err = errors.Join(err, closeErr)
		}
	}()

	var users []*dbv1.UserProto
	if err = cursor.All(ctx, &users); err != nil {
		return nil, 0, errors.Wrap(err, "decoding users failed")
	}

	return users, total, nil
}
