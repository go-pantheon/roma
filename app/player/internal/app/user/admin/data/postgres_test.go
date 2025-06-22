package data_test

import (
	"context"
	"regexp"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	basicobj "github.com/go-pantheon/roma/app/player/internal/app/basic/gate/domain/object"
	"github.com/go-pantheon/roma/app/player/internal/app/user/admin/data"
	"github.com/go-pantheon/roma/app/player/internal/app/user/admin/domain"
	idata "github.com/go-pantheon/roma/app/player/internal/data"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/encoding/protojson"
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

	repo, err := data.NewUserPostgresRepo(d, logger)
	require.NoError(t, err)
	require.NotNil(t, repo)

	return repo, mock
}

func TestUserPgRepo_GetByID(t *testing.T) {
	t.Parallel()
	repo, mock := newTestRepoWithMock(t)
	uid := int64(1001)
	user := &dbv1.UserProto{Id: uid}

	querySQL := `SELECT id, version, server_version, "basic", "testmodule" FROM "user" WHERE id = $1`

	expectedBasicModule := &dbv1.UserModuleProto{
		Module: &dbv1.UserModuleProto_Basic{
			Basic: &dbv1.UserBasicProto{Name: "get_user"},
		},
	}
	basicJSON, err := protojson.Marshal(expectedBasicModule)
	require.NoError(t, err)

	rows := pgxmock.NewRows([]string{"id", "version", "server_version", "basic", "testmodule"}).
		AddRow(uid, int64(1), "v1.0.0", basicJSON, []byte{})

	mock.ExpectQuery(regexp.QuoteMeta(querySQL)).WithArgs(uid).WillReturnRows(rows)

	err = repo.GetByID(context.Background(), user, testModuleKeys)
	require.NoError(t, err)

	assert.Equal(t, int64(1), user.Version)
	assert.Equal(t, "v1.0.0", user.ServerVersion)
	queriedBasic := user.Modules[string(basicobj.ModuleKey)].GetBasic()
	require.NotNil(t, queriedBasic)
	assert.Equal(t, "get_user", queriedBasic.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPgRepo_GetList(t *testing.T) {
	t.Parallel()

	repo, mock := newTestRepoWithMock(t)

	countSQL := `SELECT COUNT(*) FROM "user"`
	mock.ExpectQuery(regexp.QuoteMeta(countSQL)).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(int64(4)))

	matchingUserModule := &dbv1.UserModuleProto{
		Module: &dbv1.UserModuleProto_Basic{
			Basic: &dbv1.UserBasicProto{Name: "matching_user"},
		},
	}
	matchingUserJSON, err := protojson.Marshal(matchingUserModule)
	require.NoError(t, err)

	nonMatchingUserModule := &dbv1.UserModuleProto{
		Module: &dbv1.UserModuleProto_Basic{
			Basic: &dbv1.UserBasicProto{Name: "other_user"},
		},
	}
	nonMatchingUserJSON, err := protojson.Marshal(nonMatchingUserModule)
	require.NoError(t, err)

	listSQL := `SELECT "id", "version", "server_version", "basic", "testmodule" FROM "user"  ORDER BY id ASC LIMIT $1 OFFSET $2`
	rows := pgxmock.NewRows([]string{"id", "version", "server_version", "basic", "testmodule"}).
		AddRow(int64(1001), int64(1), "v1.0.0", matchingUserJSON, []byte{}).
		AddRow(int64(1002), int64(1), "v1.0.0", matchingUserJSON, []byte{}).
		AddRow(int64(1003), int64(1), "v1.0.0", nonMatchingUserJSON, []byte{}).
		AddRow(int64(1004), int64(1), "v1.0.0", nonMatchingUserJSON, []byte{})

	mock.ExpectQuery(regexp.QuoteMeta(listSQL)).WithArgs(int64(10), int64(0)).WillReturnRows(rows)

	result, count, err := repo.GetList(context.Background(), 0, 10, nil, testModuleKeys)
	require.NoError(t, err)

	assert.Equal(t, int64(4), count)
	require.Len(t, result, 4)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPgRepo_GetList_WithConditions(t *testing.T) {
	t.Parallel()

	repo, mock := newTestRepoWithMock(t)
	uid1 := int64(1001)
	uid2 := int64(1002)

	conds := map[life.ModuleKey]*dbv1.UserModuleProto{
		basicobj.ModuleKey: {
			Module: &dbv1.UserModuleProto_Basic{
				Basic: &dbv1.UserBasicProto{Name: "matching_user"},
			},
		},
	}

	basicCondModule := conds[basicobj.ModuleKey]
	basicCondJSON, err := protojson.Marshal(basicCondModule)
	require.NoError(t, err)

	whereSQL := ` WHERE "basic" @> $1::jsonb AND "basic" @> $2::jsonb`

	// mock count query
	countSQL := `SELECT COUNT(*) FROM "user"` + whereSQL
	mock.ExpectQuery(regexp.QuoteMeta(countSQL)).
		WithArgs(basicCondJSON, basicCondJSON).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(int64(2)))

	expectedBasicModule := &dbv1.UserModuleProto{
		Module: &dbv1.UserModuleProto_Basic{
			Basic: &dbv1.UserBasicProto{Name: "matching_user"},
		},
	}
	basicJSON, err := protojson.Marshal(expectedBasicModule)
	require.NoError(t, err)

	// mock list query
	listSQL := `SELECT "id", "version", "server_version", "basic", "testmodule" FROM "user"` + whereSQL + ` ORDER BY id ASC LIMIT $3 OFFSET $4`
	rows := pgxmock.NewRows([]string{"id", "version", "server_version", "basic", "testmodule"}).
		AddRow(uid1, int64(1), "v1.0.0", basicJSON, []byte{}).
		AddRow(uid2, int64(1), "v1.0.0", basicJSON, []byte{})
	mock.ExpectQuery(regexp.QuoteMeta(listSQL)).WithArgs(basicCondJSON, basicCondJSON, int64(10), int64(0)).WillReturnRows(rows)

	result, count, err := repo.GetList(context.Background(), 0, 10, conds, testModuleKeys)
	require.NoError(t, err)

	assert.Equal(t, int64(2), count)
	require.Len(t, result, 2)

	// Check user 1
	user1 := result[0]
	assert.Equal(t, uid1, user1.Id)
	queriedBasic1 := user1.Modules[string(basicobj.ModuleKey)].GetBasic()
	require.NotNil(t, queriedBasic1)
	assert.Equal(t, "matching_user", queriedBasic1.Name)

	// Check user 2
	user2 := result[1]
	assert.Equal(t, uid2, user2.Id)
	queriedBasic2 := user2.Modules[string(basicobj.ModuleKey)].GetBasic()
	require.NotNil(t, queriedBasic2)
	assert.Equal(t, "matching_user", queriedBasic2.Name)

	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUserPgRepo_GetByID_InvalidInput(t *testing.T) {
	t.Parallel()
	repo, _ := newTestRepoWithMock(t)

	err := repo.GetByID(context.Background(), nil, testModuleKeys)
	assert.Error(t, err)

	err = repo.GetByID(context.Background(), &dbv1.UserProto{Id: 0}, testModuleKeys)
	assert.Error(t, err)
}
