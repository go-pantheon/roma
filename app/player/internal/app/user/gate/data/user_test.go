package data_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	basicobj "github.com/go-pantheon/roma/app/player/internal/app/basic/gate/domain/object"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/data"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain"
	idata "github.com/go-pantheon/roma/app/player/internal/data"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

var testModuleKeys = []life.ModuleKey{
	basicobj.ModuleKey,
	"testmodule",
}

func newTestRepoWithMock(t *testing.T) (domain.UserRepo, pgxmock.PgxPoolIface) {
	t.Helper()

	mock, err := pgxmock.NewPool()
	require.NoError(t, err)

	d := &idata.Data{Pdb: mock}
	logger := log.DefaultLogger

	// Mock DB initialization
	mock.ExpectExec(regexp.QuoteMeta(`CREATE TABLE IF NOT EXISTS "user"`)).
		WillReturnResult(pgxmock.NewResult("CREATE", 1))

	queryColsSQL := `SELECT column_name FROM information_schema.columns`
	mock.ExpectQuery(regexp.QuoteMeta(queryColsSQL)).
		WithArgs("user").
		WillReturnRows(pgxmock.NewRows([]string{"column_name"}).AddRow("id"))

	// Mock column additions
	alterSQLBasic := `ALTER TABLE "user" ADD COLUMN IF NOT EXISTS "basic" JSONB;`
	mock.ExpectExec(regexp.QuoteMeta(alterSQLBasic)).WillReturnResult(pgxmock.NewResult("ALTER", 1))
	alterSQLTest := `ALTER TABLE "user" ADD COLUMN IF NOT EXISTS "testmodule" BYTEA;`
	mock.ExpectExec(regexp.QuoteMeta(alterSQLTest)).WillReturnResult(pgxmock.NewResult("ALTER", 1))

	repo, err := data.TestNewUserPgRepo(d, logger, testModuleKeys)
	require.NoError(t, err)
	require.NotNil(t, repo)

	return repo, mock
}

func TestUserPgRepo_Create(t *testing.T) {
	repo, mock := newTestRepoWithMock(t)
	uid := int64(1001)
	user := &dbv1.UserProto{
		Id:            uid,
		Version:       1,
		ServerVersion: "v1.0.0",
		Modules: map[string]*dbv1.UserModuleProto{
			string(basicobj.ModuleKey): {
				Module: &dbv1.UserModuleProto_Basic{
					Basic: &dbv1.UserBasicProto{Name: "test_user"},
				},
			},
		},
	}
	insertSQL := `INSERT INTO "user" (id, version, server_version, "basic") VALUES ($1, $2, $3, $4)`

	basicModule := user.Modules[string(basicobj.ModuleKey)]
	basicJSON, err := protojson.Marshal(basicModule)
	require.NoError(t, err)

	mock.ExpectExec(regexp.QuoteMeta(insertSQL)).
		WithArgs(uid, user.Version, user.ServerVersion, basicJSON).
		WillReturnResult(pgxmock.NewResult("INSERT", 1))

	err = repo.Create(context.Background(), uid, user, time.Now())
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPgRepo_QueryByID(t *testing.T) {
	repo, mock := newTestRepoWithMock(t)
	uid := int64(1002)

	expectedBasicModule := &dbv1.UserModuleProto{
		Module: &dbv1.UserModuleProto_Basic{
			Basic: &dbv1.UserBasicProto{Name: "queried_user"},
		},
	}
	basicJSON, err := protojson.Marshal(expectedBasicModule)
	require.NoError(t, err)

	expectedTestModule := &dbv1.UserModuleProto{}
	testModuleBytes, err := proto.Marshal(expectedTestModule)
	require.NoError(t, err)

	querySQL := `SELECT id, version, server_version, "basic", "testmodule" FROM "user" WHERE id = $1`
	rows := pgxmock.NewRows([]string{"id", "version", "server_version", "basic", "testmodule"}).
		AddRow(uid, int64(2), "v1.0.1", basicJSON, testModuleBytes)

	mock.ExpectQuery(regexp.QuoteMeta(querySQL)).WithArgs(uid).WillReturnRows(rows)

	user := &dbv1.UserProto{
		Modules: map[string]*dbv1.UserModuleProto{
			string(basicobj.ModuleKey): {},
			"testmodule":               {},
		},
	}
	err = repo.QueryByID(context.Background(), uid, user, testModuleKeys)
	require.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())

	assert.Equal(t, int64(2), user.Version)
	assert.Equal(t, "v1.0.1", user.ServerVersion)
	queriedBasic := user.Modules[string(basicobj.ModuleKey)].GetBasic()
	require.NotNil(t, queriedBasic)
	assert.Equal(t, "queried_user", queriedBasic.Name)
}

func TestUserPgRepo_UpdateByID(t *testing.T) {
	repo, mock := newTestRepoWithMock(t)
	uid := int64(1003)
	user := &dbv1.UserProto{
		Id:            uid,
		Version:       3,
		ServerVersion: "v1.0.2",
		Modules: map[string]*dbv1.UserModuleProto{
			string(basicobj.ModuleKey): {
				Module: &dbv1.UserModuleProto_Basic{
					Basic: &dbv1.UserBasicProto{Name: "updated_user"},
				},
			},
		},
	}

	updateSQL := `UPDATE "user" SET version = $1 , "basic" = $4 WHERE id = $2 AND version = $3`

	basicModule := user.Modules[string(basicobj.ModuleKey)]
	basicJSON, err := protojson.Marshal(basicModule)
	require.NoError(t, err)

	mock.ExpectExec(regexp.QuoteMeta(updateSQL)).
		WithArgs(user.Version, uid, user.Version-1, basicJSON).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err = repo.UpdateByID(context.Background(), uid, user)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPgRepo_IsExist(t *testing.T) {
	uid := int64(1004)
	querySQL := `SELECT EXISTS(SELECT 1 FROM "user" WHERE id = $1);`

	t.Run("Exists", func(t *testing.T) {
		repo, mock := newTestRepoWithMock(t)
		rows := pgxmock.NewRows([]string{"exists"}).AddRow(true)
		mock.ExpectQuery(regexp.QuoteMeta(querySQL)).WithArgs(uid).WillReturnRows(rows)
		exists, err := repo.IsExist(context.Background(), uid)
		assert.NoError(t, err)
		assert.True(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("Not Exists", func(t *testing.T) {
		repo, mock := newTestRepoWithMock(t)
		rows := pgxmock.NewRows([]string{"exists"}).AddRow(false)
		mock.ExpectQuery(regexp.QuoteMeta(querySQL)).WithArgs(uid).WillReturnRows(rows)
		exists, err := repo.IsExist(context.Background(), uid)
		assert.NoError(t, err)
		assert.False(t, exists)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestUserPgRepo_IncVersion(t *testing.T) {
	repo, mock := newTestRepoWithMock(t)
	uid := int64(1005)
	newVersion := int64(5)

	querySQL := `UPDATE "user" SET version = $1 WHERE id = $2 AND version = $3;`
	mock.ExpectExec(regexp.QuoteMeta(querySQL)).
		WithArgs(newVersion, uid, newVersion-1).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := repo.IncVersion(context.Background(), uid, newVersion)
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
