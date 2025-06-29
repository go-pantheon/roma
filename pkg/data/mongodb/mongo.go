package mongodb

import (
	"context"

	"github.com/go-pantheon/fabrica-util/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// BaseRepo provides a generic base for MongoDB repositories.
type BaseRepo struct {
	Coll *mongo.Collection
}

// NewBaseRepo creates a new BaseRepo.
func NewBaseRepo(coll *mongo.Collection) *BaseRepo {
	return &BaseRepo{Coll: coll}
}

// Create inserts a new document into the collection.
func (r *BaseRepo) Create(ctx context.Context, doc any) error {
	if doc == nil {
		return errors.New("document is nil")
	}

	_, err := r.Coll.InsertOne(ctx, doc)
	if err != nil {
		return errors.Wrap(err, "creating document failed")
	}

	return nil
}

// FindByID finds a single document by its ID field, with optional projection.
// `idKey` is the name of the ID field (e.g., "_id" or "id").
// `id` is the value of the ID.
// `result` is a pointer to the struct to decode the document into.
// `projection` is an optional BSON document specifying which fields to return.
func (r *BaseRepo) FindByID(ctx context.Context, idKey string, id any, result any, projection bson.D) error {
	if result == nil {
		return errors.New("result pointer is nil")
	}

	filter := bson.M{idKey: id}
	opts := options.FindOne()

	if len(projection) > 0 {
		opts.SetProjection(projection)
	}

	err := r.Coll.FindOne(ctx, filter, opts).Decode(result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return errors.Wrapf(err, "document not found. id=%d", id)
		}

		return errors.Wrapf(err, "document querying failed. id=%d", id)
	}

	return nil
}

// UpdateOne provides a generic way to update a single document using optimistic locking.
// `idKey` is the name of the ID field.
// `id` is the value of the ID.
// `version` is the *new* version of the document. The filter will check for `version - 1`.
// `update` is the update payload (e.g., `bson.M{"$set": ...}`).
func (r *BaseRepo) UpdateOne(ctx context.Context, idKey string, id any, version int64, update bson.M) error {
	filter := bson.M{idKey: id, "version": version - 1}

	res, err := r.Coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrapf(err, "document updating failed. id=%d version=%d", id, version)
	}

	if res.MatchedCount == 0 {
		return errors.Errorf("document not found or version mismatch. id=%d version=%d", id, version)
	}

	return nil
}

// Exists checks if a document with the given ID exists.
func (r *BaseRepo) Exists(ctx context.Context, idKey string, id any) (bool, error) {
	filter := bson.M{idKey: id}
	opts := options.Count().SetLimit(1)

	count, err := r.Coll.CountDocuments(ctx, filter, opts)
	if err != nil {
		return false, errors.Wrapf(err, "counting document failed. id=%d", id)
	}

	return count > 0, nil
}

// IncrementVersion increments the version of a document using optimistic locking.
func (r *BaseRepo) IncrementVersion(ctx context.Context, idKey string, id any, newVersion int64) error {
	filter := bson.M{idKey: id, "version": newVersion - 1}
	update := bson.M{"$set": bson.M{"version": newVersion}}

	res, err := r.Coll.UpdateOne(ctx, filter, update)
	if err != nil {
		return errors.Wrapf(err, "incrementing version failed. id=%d", id)
	}

	if res.MatchedCount == 0 {
		return errors.Errorf("document not found or version mismatch. id=%d newVersion=%d", id, newVersion)
	}

	return nil
}
