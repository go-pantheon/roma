package data

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"regexp"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/admin/domain"
	"github.com/go-pantheon/roma/app/player/internal/app/recharge/pkg"
	"github.com/go-pantheon/roma/app/player/internal/data"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type AdminOrderRepoSuite struct {
	suite.Suite
	mock pgxmock.PgxPoolIface
	repo domain.OrderRepo
}

func (s *AdminOrderRepoSuite) SetupTest() {
	var err error
	s.mock, err = pgxmock.NewPool()
	require.NoError(s.T(), err)

	repo := &orderPgRepo{
		data: &data.Data{Pdb: s.mock},
		log:  log.NewHelper(log.NewStdLogger(io.Discard)),
	}
	s.repo = repo
}

func TestAdminOrderRepoSuite(t *testing.T) {
	suite.Run(t, new(AdminOrderRepoSuite))
}

func (s *AdminOrderRepoSuite) TestGetByID() {
	store := pkg.Store("test-store")
	transId := "test-trans-id"
	query := `SELECT info, uid, store, trans_id, ack, ack_at FROM "orders" WHERE trans_id = $1 AND store = $2 LIMIT 1`

	s.T().Run("GetSuccess", func(t *testing.T) {
		s.SetupTest()
		info := &dbv1.OrderInfoProto{ProductId: "test_product"}
		infoJson, err := json.Marshal(info)
		require.NoError(t, err)

		rows := pgxmock.NewRows([]string{"info", "uid", "store", "trans_id", "ack", "ack_at"}).
			AddRow(infoJson, int64(123), "test-store", transId, int32(dbv1.OrderAckState_ORDER_ACK_STATE_PENDING), time.Now().Unix())

		s.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(transId, store).WillReturnRows(rows)

		order, err := s.repo.GetByID(context.Background(), store, transId)
		require.NoError(t, err)
		require.NotNil(t, order)
		require.NoError(t, s.mock.ExpectationsWereMet())
	})

	s.T().Run("NotFound", func(t *testing.T) {
		s.SetupTest()
		s.mock.ExpectQuery(regexp.QuoteMeta(query)).WithArgs(transId, store).WillReturnError(pgx.ErrNoRows)

		_, err := s.repo.GetByID(context.Background(), store, transId)
		require.Error(t, err)
		require.NoError(t, s.mock.ExpectationsWereMet())
	})
}

func (s *AdminOrderRepoSuite) TestUpdateAckStateByID() {
	store := pkg.Store("test-store")
	transId := "test-trans-id"
	state := dbv1.OrderAckState_ORDER_ACK_STATE_SUCCEEDED
	query := `UPDATE "orders" SET ack = $1, ack_at = $2 WHERE trans_id = $3 AND store = $4`

	s.T().Run("UpdateSuccess", func(t *testing.T) {
		s.SetupTest()
		s.mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(int32(state), pgxmock.AnyArg(), transId, store).
			WillReturnResult(pgxmock.NewResult("UPDATE", 1))

		err := s.repo.UpdateAckStateByID(context.Background(), store, transId, state)
		require.NoError(t, err)
		require.NoError(t, s.mock.ExpectationsWereMet())
	})

	s.T().Run("UpdateNotFound", func(t *testing.T) {
		s.SetupTest()
		s.mock.ExpectExec(regexp.QuoteMeta(query)).
			WithArgs(int32(state), pgxmock.AnyArg(), transId, store).
			WillReturnResult(pgxmock.NewResult("UPDATE", 0))

		err := s.repo.UpdateAckStateByID(context.Background(), store, transId, state)
		require.Error(t, err)
		require.NoError(t, s.mock.ExpectationsWereMet())
	})
}

func (s *AdminOrderRepoSuite) TestGetList() {
	info := &dbv1.OrderInfoProto{ProductId: "test_product"}
	infoJson, err := json.Marshal(info)
	require.NoError(s.T(), err)

	rows := pgxmock.NewRows([]string{"info", "uid", "store", "trans_id", "ack", "ack_at"}).
		AddRow(infoJson, int64(123), "test-store", "test-trans-id-1", int32(dbv1.OrderAckState_ORDER_ACK_STATE_PENDING), time.Now().Unix()).
		AddRow(infoJson, int64(456), "test-store", "test-trans-id-2", int32(dbv1.OrderAckState_ORDER_ACK_STATE_SUCCEEDED), time.Now().Unix())

	s.T().Run("GetListNoConditions", func(t *testing.T) {
		s.SetupTest()
		countQuery := `SELECT COUNT(*) FROM "orders"`
		listQuery := `SELECT info, uid, store, trans_id, ack, ack_at FROM "orders"  ORDER BY ack_at DESC LIMIT 10 OFFSET 0`

		s.mock.ExpectQuery(regexp.QuoteMeta(countQuery)).WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(int64(2)))
		s.mock.ExpectQuery(regexp.QuoteMeta(listQuery)).WillReturnRows(rows)

		orders, total, err := s.repo.GetList(context.Background(), 0, 10, nil)
		require.NoError(t, err)
		require.Equal(t, int64(2), total)
		require.Len(t, orders, 2)
		require.NoError(t, s.mock.ExpectationsWereMet())
	})

	s.T().Run("GetListWithUidCondition", func(t *testing.T) {
		s.SetupTest()
		cond := &dbv1.OrderProto{Uid: 123}
		where := "WHERE uid = $1"
		countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM "orders" %s`, where)
		listQuery := fmt.Sprintf(`SELECT info, uid, store, trans_id, ack, ack_at FROM "orders" %s ORDER BY ack_at DESC LIMIT 10 OFFSET 0`, where)

		uidRows := pgxmock.NewRows([]string{"info", "uid", "store", "trans_id", "ack", "ack_at"}).
			AddRow(infoJson, int64(123), "test-store", "test-trans-id-1", int32(dbv1.OrderAckState_ORDER_ACK_STATE_PENDING), time.Now().Unix())

		s.mock.ExpectQuery(regexp.QuoteMeta(countQuery)).WithArgs(cond.Uid).WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(int64(1)))
		s.mock.ExpectQuery(regexp.QuoteMeta(listQuery)).WithArgs(cond.Uid).WillReturnRows(uidRows)

		orders, total, err := s.repo.GetList(context.Background(), 0, 10, cond)
		require.NoError(t, err)
		require.Equal(t, int64(1), total)
		require.Len(t, orders, 1)
		require.Equal(t, cond.Uid, orders[0].Uid)
		require.NoError(t, s.mock.ExpectationsWereMet())
	})

	s.T().Run("GetListWithInfoCondition", func(t *testing.T) {
		s.SetupTest()
		cond := &dbv1.OrderProto{
			Info: &dbv1.OrderInfoProto{
				ProductId: "test_product_id",
			},
		}

		infoWithCondition := &dbv1.OrderInfoProto{ProductId: "test_product_id"}
		infoJsonWithCondition, err := json.Marshal(infoWithCondition)
		require.NoError(t, err)

		where := "WHERE info->>'product_id' = $1"
		countQuery := fmt.Sprintf(`SELECT COUNT(*) FROM "orders" %s`, where)
		listQuery := fmt.Sprintf(`SELECT info, uid, store, trans_id, ack, ack_at FROM "orders" %s ORDER BY ack_at DESC LIMIT 10 OFFSET 0`, where)

		infoRows := pgxmock.NewRows([]string{"info", "uid", "store", "trans_id", "ack", "ack_at"}).
			AddRow(infoJsonWithCondition, int64(789), "another-store", "test-trans-id-3", int32(dbv1.OrderAckState_ORDER_ACK_STATE_SUCCEEDED), time.Now().Unix())

		s.mock.ExpectQuery(regexp.QuoteMeta(countQuery)).WithArgs(cond.Info.ProductId).WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(int64(1)))
		s.mock.ExpectQuery(regexp.QuoteMeta(listQuery)).WithArgs(cond.Info.ProductId).WillReturnRows(infoRows)

		orders, total, err := s.repo.GetList(context.Background(), 0, 10, cond)
		require.NoError(t, err)
		require.Equal(t, int64(1), total)
		require.Len(t, orders, 1)
		require.Equal(t, cond.Info.ProductId, orders[0].Info.ProductId)
		require.NoError(t, s.mock.ExpectationsWereMet())
	})
}
