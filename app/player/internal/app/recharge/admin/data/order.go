package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/admin/domain"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/pkg"
	"github.com/go-pantheon/roma/app/player/internal/data"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var _ domain.OrderRepo = (*orderMongoRepo)(nil)

type orderMongoRepo struct {
	log        *log.Helper
	data       *data.Data
	collection *mongo.Collection
}

func NewOrderMongoRepo(data *data.Data, logger log.Logger) (domain.OrderRepo, error) {
	r := &orderMongoRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "player/recharge/admin/data/order")),
	}

	r.collection = data.Mdb.Collection("order")

	return r, nil
}

func (r *orderMongoRepo) GetByID(ctx context.Context, store pkg.Store, transId string) (*dbv1.OrderProto, error) {
	filter := bson.D{
		bson.E{Key: "store", Value: store},
		bson.E{Key: "transId", Value: transId},
	}
	result := r.collection.FindOne(ctx, filter)
	if err := result.Err(); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, xerrors.APINotFound("transId=%s store=%s", transId, store)
		}
		return nil, xerrors.APIDBFailed("store=%s transId=%s", store, transId).WithCause(err)
	}

	bo := &dbv1.OrderProto{}
	if err := result.Decode(bo); err != nil {
		return nil, xerrors.APICodecFailed("store=%s transId=%s", store, transId).WithCause(err)
	}
	return bo, nil
}

var orderListFields = bson.D{
	bson.E{Key: "store", Value: 1},
	bson.E{Key: "transId", Value: 1},
	bson.E{Key: "uid", Value: 1},
	bson.E{Key: "ack", Value: 1},
	bson.E{Key: "productId", Value: 1},
	bson.E{Key: "purchaseAt", Value: 1},
	bson.E{Key: "ackAt", Value: 1},
}

func (r *orderMongoRepo) GetList(ctx context.Context, index, limit int64, cond *dbv1.OrderProto) (ret []*dbv1.OrderProto, count int64, err error) {
	filter := r.buildFilter(ctx, cond)

	count, err = r.collection.CountDocuments(ctx, filter)
	if err != nil {
		err = xerrors.APIDBFailed("count documents failed").WithCause(err)
		return
	}

	opts := options.Find().SetSort(bson.D{bson.E{Key: "purchaseAt", Value: -1}}).SetSkip(int64(index)).SetLimit(int64(limit))
	opts = opts.SetProjection(orderListFields)
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			ret = make([]*dbv1.OrderProto, 0)
			return
		}
		err = xerrors.APIDBFailed("get order list failed").WithCause(err)
		return
	}

	orders := make([]dbv1.OrderProto, 0, limit)
	if err = cursor.All(ctx, &orders); err != nil {
		err = xerrors.APIDBFailed("create order list failed").WithCause(err)
		return
	}

	ret = make([]*dbv1.OrderProto, 0, len(orders))
	for i := 0; i < len(orders); i++ {
		ret = append(ret, &orders[i])
	}
	return
}

func (r *orderMongoRepo) buildFilter(_ context.Context, cond *dbv1.OrderProto) bson.D {
	filter := make(bson.D, 0, 8)
	if len(cond.Store) > 0 {
		filter = append(filter, bson.E{Key: "store", Value: cond.Store})
	}
	if len(cond.TransId) > 0 {
		filter = append(filter, bson.E{Key: "transId", Value: cond.TransId})
	}
	if cond.Ack > 0 {
		filter = append(filter, bson.E{Key: "ack", Value: cond.Ack})
	}
	if cond.Uid > 0 {
		filter = append(filter, bson.E{Key: "uid", Value: cond.Uid})
	}
	return filter
}

func (r *orderMongoRepo) UpdateAckStateByID(ctx context.Context, store pkg.Store, transId string, state dbv1.OrderAckState) error {
	filter := bson.D{
		bson.E{Key: "store", Value: store},
		bson.E{Key: "transId", Value: transId},
	}
	result, err := r.collection.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: bson.D{
			bson.E{Key: "ack", Value: state},
		}},
	})
	if err != nil {
		return xerrors.APIDBFailed("store=%s transId=%s ack=%d", store, transId, state).WithCause(err)
	}
	if result.ModifiedCount < 1 {
		return xerrors.APIDBNoAffected("store=%s transId=%s ack=%d", store, transId, state)
	}
	return nil
}
