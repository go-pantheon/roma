package data

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/fabrica-util/data/db/pg"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/admin/domain"
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
		log:  log.NewHelper(log.With(logger, "module", "player/recharge/admin/data/order")),
	}

	return r, nil
}

func (r *orderPgRepo) GetByID(ctx context.Context, store pkg.Store, transId string) (*dbv1.OrderProto, error) {
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

func (r *orderPgRepo) GetList(ctx context.Context, index, limit int64, cond *dbv1.OrderProto) (orders []*dbv1.OrderProto, total int64, err error) {
	where, args := r.buildWhere(cond)
	countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM "%s" %s`, _tableName, where)

	if err := r.data.QueryRowContext(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, errors.Wrapf(err, "failed to count orders")
	}

	if total == 0 {
		return nil, 0, nil
	}

	query := fmt.Sprintf(`SELECT info, uid, store, trans_id, ack, ack_at FROM "%s" %s ORDER BY ack_at DESC LIMIT %d OFFSET %d`, _tableName, where, limit, index)

	rows, err := r.data.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, 0, errors.Wrapf(err, "failed to get order list")
	}

	defer rows.Close()

	for rows.Next() {
		var (
			o        dbv1.OrderProto
			infoJson []byte
		)

		if err := rows.Scan(&infoJson, &o.Uid, &o.Store, &o.TransId, &o.Ack, &o.AckAt); err != nil {
			return nil, 0, errors.Wrapf(err, "failed to scan order")
		}

		err = json.Unmarshal(infoJson, &o.Info)
		if err != nil {
			return nil, 0, errors.Wrapf(err, "failed to unmarshal order info: %s", string(infoJson))
		}

		orders = append(orders, &o)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, errors.Wrapf(err, "error reading order rows")
	}

	return orders, total, nil
}

func (r *orderPgRepo) UpdateAckStateByID(ctx context.Context, store pkg.Store, transId string, state dbv1.OrderAckState) error {
	query := fmt.Sprintf(`UPDATE "%s" SET ack = $1, ack_at = $2 WHERE trans_id = $3 AND store = $4`, _tableName)

	tag, err := r.data.ExecContext(ctx, query, int32(state), time.Now().Unix(), transId, store)
	if err != nil {
		return errors.Wrapf(err, "failed to update ack state, store: %s, transId: %s, state: %s", store, transId, state.String())
	}

	if res, err := tag.RowsAffected(); err != nil {
		return errors.Wrapf(err, "failed to get rows affected")
	} else if res == 0 {
		return xerrors.ErrDBRecordNotAffected
	}

	return nil
}

func (r *orderPgRepo) buildWhere(cond *dbv1.OrderProto) (string, []any) {
	if cond == nil {
		return "", nil
	}

	var (
		where []string
		args  []any
	)

	argId := 1

	if cond.Uid > 0 {
		where = append(where, fmt.Sprintf("uid = $%d", argId))
		args = append(args, cond.Uid)
		argId++
	}

	if cond.Store != "" {
		where = append(where, fmt.Sprintf("store = $%d", argId))
		args = append(args, cond.Store)
		argId++
	}

	if cond.TransId != "" {
		where = append(where, fmt.Sprintf("trans_id = $%d", argId))
		args = append(args, cond.TransId)
		argId++
	}

	if cond.Ack != int32(dbv1.OrderAckState_ORDER_ACK_STATE_UNSPECIFIED) {
		where = append(where, fmt.Sprintf("ack = $%d", argId))
		args = append(args, cond.Ack)
		argId++
	}

	if cond.Info != nil {
		if cond.Info.Token != "" {
			where = append(where, fmt.Sprintf("info->>'token' = $%d", argId))
			args = append(args, cond.Info.Token)
			argId++
		}

		if cond.Info.Env != "" {
			where = append(where, fmt.Sprintf("info->>'env' = $%d", argId))
			args = append(args, cond.Info.Env)
			argId++
		}

		if cond.Info.ProductId != "" {
			where = append(where, fmt.Sprintf("info->>'product_id' = $%d", argId))
			args = append(args, cond.Info.ProductId)
			argId++
		}

		if cond.Info.PurchasedAt > 0 {
			where = append(where, fmt.Sprintf("(info->>'purchased_at')::bigint = $%d", argId))
			args = append(args, cond.Info.PurchasedAt)
		}
	}

	if len(where) == 0 {
		return "", nil
	}

	return "WHERE " + strings.Join(where, " AND "), args
}
