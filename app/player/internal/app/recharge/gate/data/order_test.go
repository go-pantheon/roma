package data

import (
	"context"
	"encoding/json"
	"io"
	"regexp"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/gate/domain"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/pkg"
	"github.com/go-pantheon/roma/app/player/internal/data"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type OrderRepoSuite struct {
	suite.Suite
	mock pgxmock.PgxPoolIface
	repo domain.OrderRepo
}

func (s *OrderRepoSuite) SetupTest() {
	var err error
	s.mock, err = pgxmock.NewPool()
	require.NoError(s.T(), err)

	repo := &orderPgRepo{
		data: &data.Data{Pdb: s.mock},
		log:  log.NewHelper(log.NewStdLogger(io.Discard)),
	}
	s.repo = repo
}

func TestOrderRepoSuite(t *testing.T) {
	suite.Run(t, new(OrderRepoSuite))
}

func (s *OrderRepoSuite) TestCreate() {
	s.T().Run("CreateOrderSuccess", func(t *testing.T) {
		s.SetupTest()
		order := &dbv1.OrderProto{
			Info:    &dbv1.OrderInfoProto{ProductId: "test_product"},
			Uid:     12345,
			Store:   "test-store",
			TransId: "test-trans-id",
			Ack:     int32(dbv1.OrderAckState_ORDER_ACK_STATE_PENDING),
			AckAt:   time.Now().Unix(),
		}

		query := `INSERT INTO "orders" (info, uid, store, trans_id, ack, ack_at) VALUES ($1, $2, $3, $4, $5, $6)`
		s.mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(pgxmock.AnyArg(), order.Uid, order.Store, order.TransId, order.Ack, order.AckAt).
			WillReturnResult(pgxmock.NewResult("INSERT", 1))

		err := s.repo.Create(context.Background(), order)
		require.NoError(t, err)
		require.NoError(t, s.mock.ExpectationsWereMet())
	})

	s.T().Run("CreateOrderDBFail", func(t *testing.T) {
		s.SetupTest()
		order := &dbv1.OrderProto{
			Info:    &dbv1.OrderInfoProto{ProductId: "test_product"},
			Uid:     12345,
			Store:   "test-store",
			TransId: "test-trans-id",
			Ack:     int32(dbv1.OrderAckState_ORDER_ACK_STATE_PENDING),
			AckAt:   time.Now().Unix(),
		}

		query := `INSERT INTO "orders" (info, uid, store, trans_id, ack, ack_at) VALUES ($1, $2, $3, $4, $5, $6)`
		s.mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(pgxmock.AnyArg(), order.Uid, order.Store, order.TransId, order.Ack, order.AckAt).
			WillReturnError(errors.New("db error"))

		err := s.repo.Create(context.Background(), order)
		require.Error(t, err)
		require.NoError(t, s.mock.ExpectationsWereMet())
	})
}

func (s *OrderRepoSuite) TestGetByTransId() {
	store := pkg.Store("test-store")
	transId := "test-trans-id"
	query := `SELECT info, uid, store, trans_id, ack, ack_at FROM "orders" WHERE trans_id = $1 AND store = $2 LIMIT 1`

	s.T().Run("GetOrderSuccess", func(t *testing.T) {
		s.SetupTest()
		info := &dbv1.OrderInfoProto{ProductId: "test_product"}
		infoJson, err := json.Marshal(info)
		require.NoError(t, err)

		rows := pgxmock.NewRows([]string{"info", "uid", "store", "trans_id", "ack", "ack_at"}).
			AddRow(infoJson, int64(12345), "test-store", transId, int32(dbv1.OrderAckState_ORDER_ACK_STATE_PENDING), time.Now().Unix())

		s.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(transId, store).WillReturnRows(rows)

		order, err := s.repo.GetByTransId(context.Background(), store, transId)
		require.NoError(t, err)
		require.NotNil(t, order)
		require.Equal(t, transId, order.TransId)
		require.NoError(t, s.mock.ExpectationsWereMet())
	})

	s.T().Run("GetOrderNotFound", func(t *testing.T) {
		s.SetupTest()
		s.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(transId, store).WillReturnError(pgx.ErrNoRows)

		_, err := s.repo.GetByTransId(context.Background(), store, transId)
		require.Error(t, err)
		require.True(t, errors.Is(err, pgx.ErrNoRows) || (err != nil && err.Error() == "order not found, store: test-store, transId: test-trans-id"))
		require.NoError(t, s.mock.ExpectationsWereMet())
	})

	s.T().Run("GetOrderScanFail", func(t *testing.T) {
		s.SetupTest()
		rows := pgxmock.NewRows([]string{"info"}).AddRow([]byte("invalid-json"))
		s.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(transId, store).WillReturnRows(rows)

		_, err := s.repo.GetByTransId(context.Background(), store, transId)
		require.Error(t, err)
		require.NoError(t, s.mock.ExpectationsWereMet())
	})
}

func (s *OrderRepoSuite) TestUpdateAckState() {
	store := pkg.Store("test-store")
	transId := "test-trans-id"
	ackState := dbv1.OrderAckState_ORDER_ACK_STATE_SUCCEEDED
	query := `UPDATE orders SET ack = $1, ack_at = $2 WHERE trans_id = $3 AND store = $4`

	s.T().Run("UpdateAckSuccess", func(t *testing.T) {
		s.SetupTest()
		s.mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(int32(ackState), pgxmock.AnyArg(), transId, store).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := s.repo.UpdateAckState(context.Background(), store, transId, ackState)
		require.NoError(t, err)
		require.NoError(t, s.mock.ExpectationsWereMet())
	})

	s.T().Run("UpdateAckNotFound", func(t *testing.T) {
		s.SetupTest()
		s.mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(int32(ackState), pgxmock.AnyArg(), transId, store).
			WillReturnResult(pgxmock.NewResult("UPDATE", 0))

		err := s.repo.UpdateAckState(context.Background(), store, transId, ackState)
		require.Error(t, err)
		require.Equal(t, "order not found, store: test-store, transId: test-trans-id", err.Error())
		require.NoError(t, s.mock.ExpectationsWereMet())
	})
}
