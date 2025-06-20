package data

import (
	"context"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/gate/domain"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/pkg"
	"github.com/go-pantheon/roma/app/player/internal/data"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
)

const (
	_tableName = "order"
)

var _ domain.OrderRepo = (*orderPostgresRepo)(nil)

type orderPostgresRepo struct {
	log  *log.Helper
	data *data.Data
}

func NewOrderPostgresRepo(data *data.Data, logger log.Logger) (domain.OrderRepo, error) {
	r := &orderPostgresRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "player/recharge/gate/data")),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := r.initDB(ctx); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *orderPostgresRepo) initDB(ctx context.Context) error {
	if err := r.createTable(ctx); err != nil {
		return err
	}

	if err := r.updateTable(ctx); err != nil {
		return err
	}

	return nil
}

func (r *orderPostgresRepo) createTable(ctx context.Context) error {
	const createTableSQL = `
	CREATE TABLE IF NOT EXISTS "order" (
		"id" BIGINT PRIMARY KEY,
		"uid" BIGINT NOT NULL,
		"store" VARCHAR(255) NOT NULL,
		"trans_id" VARCHAR(255) NOT NULL,
		"ack" INT NOT NULL,
		"ack_at" BIGINT NOT NULL,
		"info" JSONB NOT NULL
	);
	`
	_, err := r.data.Pdb.Exec(ctx, createTableSQL)
	if err != nil {
		return errors.Wrap(err, "failed to create user table")
	}

	r.log.Infof("[PostgreSQL] created user table")

	return nil
}

func (r *orderPostgresRepo) updateTable(ctx context.Context) error {
	return errors.New("not implemented")
}

func (r *orderPostgresRepo) Create(ctx context.Context, order *dbv1.OrderProto) error {
	return errors.New("not implemented")
}

func (r *orderPostgresRepo) GetByTransId(ctx context.Context, store pkg.Store, transId string) (*dbv1.OrderProto, error) {
	return nil, errors.New("not implemented")
}

func (r *orderPostgresRepo) UpdateAckState(ctx context.Context, store pkg.Store, transId string, ackState dbv1.OrderAckState) error {
	return errors.New("not implemented")
}
