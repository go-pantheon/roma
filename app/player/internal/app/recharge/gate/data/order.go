package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/gate/domain"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/pkg"
	"github.com/go-pantheon/roma/app/player/internal/data"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var _ domain.OrderRepo = (*orderMongoRepo)(nil)

// TODO: base on PostgreSQL
type orderMongoRepo struct {
	log        *log.Helper
	data       *data.Data
	collection *mongo.Collection
}

func NewOrderMongoRepo(data *data.Data, logger log.Logger) (domain.OrderRepo, error) {
	r := &orderMongoRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "player/recharge/gate/data/order")),
	}

	r.collection = data.Mdb.Collection("order")

	if _, err := r.collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.M{"uid": 1},
		Options: options.Index().SetName("uid_idx"),
	}); err != nil {
		return nil, errors.Wrapf(err, "create index[uid_idx] failed.")
	}
	if _, err := r.collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.M{"store": 1},
		Options: options.Index().SetName("store_idx"),
	}); err != nil {
		return nil, errors.Wrapf(err, "create index[store_idx] failed.")
	}
	if _, err := r.collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.M{"transId": 1},
		Options: options.Index().SetName("trans_id_idx"),
	}); err != nil {
		return nil, errors.Wrapf(err, "create index[trans_id_idx] failed.")
	}
	if _, err := r.collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{
			bson.E{Key: "store", Value: 1},
			bson.E{Key: "transId", Value: 1},
		},
		Options: options.Index().SetName("store_trans_id_idx").SetUnique(true),
	}); err != nil {
		return nil, errors.Wrapf(err, "create index[store_trans_id_idx] failed.")
	}
	return r, nil
}

func (r *orderMongoRepo) Create(ctx context.Context, order *dbv1.OrderProto) (err error) {
	if order == nil {
		return errors.Wrapf(err, "create order failed. order is nil")
	}

	filter := bson.D{
		bson.E{Key: "store", Value: order.Store},
		bson.E{Key: "transId", Value: order.TransId},
	}

	if _, err = r.collection.ReplaceOne(ctx, filter, order, options.Replace().SetUpsert(true)); err != nil {
		return errors.Wrapf(err, "create order failed. store=%s transId=%s", order.Store, order.TransId)
	}
	return nil
}

func (r *orderMongoRepo) GetByTransId(ctx context.Context, store pkg.Store, transId string) (*dbv1.OrderProto, error) {
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

func (r *orderMongoRepo) UpdateAckState(ctx context.Context, store pkg.Store, transId string, ack dbv1.OrderAckState) error {
	filter := bson.D{
		bson.E{Key: "store", Value: store},
		bson.E{Key: "transId", Value: transId},
	}
	result, err := r.collection.UpdateOne(ctx, filter, bson.D{
		{Key: "$set", Value: bson.D{
			bson.E{Key: "ack", Value: ack},
		}},
	})
	if err != nil {
		return errors.Wrapf(err, "update order failed. store=%s transId=%s ack=%d", store, transId, ack)
	}
	if result.ModifiedCount < 1 {
		return errors.Wrapf(xerrors.ErrDBRecordNotAffected, "update order failed. store=%s transId=%s ack=%d", store, transId, ack)
	}
	return nil
}
