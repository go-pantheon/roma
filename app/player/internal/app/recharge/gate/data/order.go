package data

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/fabrica-util/data/db/pg"
	"github.com/go-pantheon/fabrica-util/data/db/pg/migrate"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/gate/domain"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/pkg"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/jackc/pgx/v5"
)

const (
	_tableName = "orders"
)

var _ domain.OrderRepo = (*orderPgRepo)(nil)

type orderPgRepo struct {
	log  *log.Helper
	data *pg.DB
}

func NewOrderPgRepo(data *pg.DB, logger log.Logger) (domain.OrderRepo, error) {
	r := &orderPgRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "player/recharge/gate/data")),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := migrate.Migrate(ctx, r.data, _tableName, &dbv1.OrderProto{}, nil); err != nil {
		return nil, err
	}

	if err := r.updateDB(ctx); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *orderPgRepo) updateDB(ctx context.Context) error {
	{
		idxTransIdStoreSQL := `CREATE INDEX IF NOT EXISTS idx_orders_trans_id_store ON orders (trans_id, store);`
		if _, err := r.data.ExecContext(ctx, idxTransIdStoreSQL); err != nil {
			return errors.Wrapf(err, "failed to create index idx_orders_trans_id_store")
		}
	}

	{
		idxUIDSQL := `CREATE INDEX IF NOT EXISTS idx_uid ON orders (uid);`
		if _, err := r.data.ExecContext(ctx, idxUIDSQL); err != nil {
			return errors.Wrapf(err, "failed to create index idx_uid")
		}
	}

	{
		idxInfoProductIdSQL := `CREATE INDEX IF NOT EXISTS idx_orders_info_product_id ON orders ((info->>'product_id'));`
		if _, err := r.data.ExecContext(ctx, idxInfoProductIdSQL); err != nil {
			return errors.Wrapf(err, "failed to create index idx_orders_info_product_id")
		}
	}

	return nil
}

func (r *orderPgRepo) Create(ctx context.Context, order *dbv1.OrderProto) error {
	infoJson, err := json.Marshal(order.Info)
	if err != nil {
		return errors.Wrapf(err, "failed to marshal order info: %+v", order.Info)
	}

	query := fmt.Sprintf(`INSERT INTO "%s" (info, uid, store, trans_id, ack, ack_at)
		VALUES ($1, $2, $3, $4, $5, $6)`, _tableName)

	_, err = r.data.ExecContext(ctx, query, infoJson, order.Uid, order.Store, order.TransId, order.Ack, order.AckAt)
	if err != nil {
		return errors.Wrapf(err, "failed to create order, order: %+v", order)
	}

	return nil
}

func (r *orderPgRepo) GetByTransId(ctx context.Context, store pkg.Store, transId string) (*dbv1.OrderProto, error) {
	query := fmt.Sprintf(`SELECT info, uid, store, trans_id, ack, ack_at
		FROM "%s" WHERE trans_id = $1 AND store = $2 LIMIT 1`, _tableName)

	row := r.data.QueryRowContext(ctx, query, transId, store)

	var (
		o        dbv1.OrderProto
		infoJson []byte
	)

	err := row.Scan(&infoJson, &o.Uid, &o.Store, &o.TransId, &o.Ack, &o.AckAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.Errorf("order not found, store: %s, transId: %s", store, transId)
		}

		return nil, errors.Wrapf(err, "failed to get order by transId, store: %s, transId: %s", store, transId)
	}

	if err = json.Unmarshal(infoJson, &o.Info); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal order info: %s", string(infoJson))
	}

	return &o, nil
}

func (r *orderPgRepo) UpdateAckState(ctx context.Context, store pkg.Store, transId string, ackState dbv1.OrderAckState) error {
	query := `UPDATE orders SET ack = $1, ack_at = $2 WHERE trans_id = $3 AND store = $4`

	tag, err := r.data.ExecContext(ctx, query, int32(ackState), time.Now().Unix(), transId, store)
	if err != nil {
		return errors.Wrapf(err, "failed to update ack state, store: %s, transId: %s, ackState: %s", store, transId, ackState.String())
	}

	ret, err := tag.RowsAffected()
	if err != nil {
		return errors.Wrapf(err, "failed to get rows affected")
	}

	if ret == 0 {
		return xerrors.ErrDBRecordNotAffected
	}

	return nil
}
