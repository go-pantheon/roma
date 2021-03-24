package data

import (
	"context"

	"github.com/go-kratos/kratos/log"
	"github.com/pkg/errors"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/recharge/admin/domain"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/app/recharge/pkg"
	"github.com/vulcan-frame/vulcan-game/app/player/internal/data"
	dbv1 "github.com/vulcan-frame/vulcan-game/gen/api/db/player/v1"
	"github.com/vulcan-frame/vulcan-kit/xerrors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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
			return nil, xerrors.ErrDBRecordNotFound
		}
		return nil, errors.Wrapf(err, "get order failed. store=%s transId=%s", store, transId)
	}

	bo := &dbv1.OrderProto{}
	if err := result.Decode(bo); err != nil {
		return nil, errors.Wrapf(err, "decode order failed. store=%s transId=%s", store, transId)
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

func (r *orderMongoRepo) GetList(ctx context.Context, index, limit int32, cond *dbv1.OrderProto) (result []*dbv1.OrderProto, count int64, err error) {
	filter := r.buildFilter(ctx, cond)

	count, err = r.collection.CountDocuments(ctx, filter)
	if err != nil {
		err = errors.Wrapf(err, "get order list failed. count documents failed.")
		return
	}

	opts := options.Find().SetSort(bson.D{bson.E{Key: "purchaseAt", Value: -1}}).SetSkip(int64(index)).SetLimit(int64(limit))
	opts = opts.SetProjection(orderListFields)
	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			err = xerrors.ErrDBRecordNotFound
			return
		}
		err = errors.Wrapf(err, "get order list failed.")
		return
	}

	orders := make([]dbv1.OrderProto, 0, limit)
	if err = cursor.All(ctx, &orders); err != nil {
		err = errors.Wrapf(err, "create order list failed.")
		return
	}

	result = make([]*dbv1.OrderProto, 0, len(orders))
	for i := 0; i < len(orders); i++ {
		result = append(result, &orders[i])
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
		return errors.Wrapf(err, "update order failed. store=%s transId=%s ack=%d", store, transId, state)
	}
	if result.ModifiedCount < 1 {
		return errors.Wrapf(xerrors.ErrDBRecordNotAffected, "update order failed. store=%s transId=%s ack=%d", store, transId, state)
	}
	return nil
}
