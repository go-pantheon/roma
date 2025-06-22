//go:build integration

package data_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-util/data/db/mongo"
	"github.com/go-pantheon/roma/app/room/internal/app/room/gate/data"
	"github.com/go-pantheon/roma/app/room/internal/app/room/gate/domain"
	pdata "github.com/go-pantheon/roma/app/room/internal/data"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/room/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRoomTestRepo(t *testing.T) (domain.RoomRepo, func()) {
	t.Helper()

	dbName := "room_gate_test"
	dbsn := "127.0.0.1:27017"

	mdb, mdbCleanup, err := mongo.New(context.Background(), dbsn, dbName)
	require.NoError(t, err)

	err = mdb.Drop(context.Background())
	require.NoError(t, err)

	d := &pdata.Data{Mdb: mdb}
	logger := log.NewStdLogger(os.Stdout)

	repo, err := data.NewRoomMongoRepo(d, logger)
	require.NoError(t, err)

	return repo, mdbCleanup
}

func TestRoomMongoRepo(t *testing.T) {
	repo, cleanup := setupRoomTestRepo(t)
	defer cleanup()

	ctx := context.Background()
	roomID := int64(2001)
	ctime := time.Now()

	defaultRoom := &dbv1.RoomProto{
		Id:       roomID,
		Sid:      1,
		RoomType: 1,
		Members: []*dbv1.RoomMemberProto{
			{Id: 1001, JoinedAt: ctime.Unix()},
		},
		Version: 1,
	}

	t.Run("Create", func(t *testing.T) {
		err := repo.Create(ctx, defaultRoom, ctime)
		require.NoError(t, err)

		// Test creating a room that already exists (_id is unique)
		err = repo.Create(ctx, defaultRoom, ctime)
		require.Error(t, err)

		// Test creating with nil room proto
		err = repo.Create(ctx, nil, ctime)
		require.Error(t, err)

		// Test creating with zero ID
		err = repo.Create(ctx, &dbv1.RoomProto{Id: 0}, ctime)
		require.Error(t, err)
	})

	t.Run("Exist", func(t *testing.T) {
		exists, err := repo.Exist(ctx, roomID)
		require.NoError(t, err)
		assert.True(t, exists)

		exists, err = repo.Exist(ctx, roomID+1)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("QueryByID", func(t *testing.T) {
		queriedRoom := &dbv1.RoomProto{}
		err := repo.QueryByID(ctx, roomID, queriedRoom)
		require.NoError(t, err)
		assert.Equal(t, defaultRoom.Id, queriedRoom.Id)
		assert.Equal(t, defaultRoom.Sid, queriedRoom.Sid)
		assert.Equal(t, defaultRoom.RoomType, queriedRoom.RoomType)
		assert.Equal(t, defaultRoom.Version, queriedRoom.Version)
		require.Len(t, queriedRoom.Members, 1)
		assert.Equal(t, defaultRoom.Members[0].Id, queriedRoom.Members[0].Id)

		err = repo.QueryByID(ctx, roomID+1, &dbv1.RoomProto{})
		require.Error(t, err)
	})

	t.Run("IncVersion", func(t *testing.T) {
		newVersion := defaultRoom.Version + 1
		err := repo.IncVersion(ctx, roomID, newVersion)
		require.NoError(t, err)

		queriedRoom := &dbv1.RoomProto{}
		err = repo.QueryByID(ctx, roomID, queriedRoom)
		require.NoError(t, err)
		assert.Equal(t, newVersion, queriedRoom.Version)

		// Test increment with wrong version
		err = repo.IncVersion(ctx, roomID, newVersion)
		require.Error(t, err)
	})

	t.Run("UpdateByID", func(t *testing.T) {
		currentRoom := &dbv1.RoomProto{}
		err := repo.QueryByID(ctx, roomID, currentRoom)
		require.NoError(t, err)

		updatedRoom := &dbv1.RoomProto{
			Version: currentRoom.Version + 1,
			Members: []*dbv1.RoomMemberProto{
				{Id: 1001, JoinedAt: ctime.Unix()},
				{Id: 1002, JoinedAt: ctime.Unix() + 1},
			},
			RoomType: 2,
		}

		err = repo.UpdateByID(ctx, roomID, updatedRoom)
		require.NoError(t, err)

		queriedAfterUpdate := &dbv1.RoomProto{}
		err = repo.QueryByID(ctx, roomID, queriedAfterUpdate)
		require.NoError(t, err)
		assert.Equal(t, updatedRoom.Version, queriedAfterUpdate.Version)
		assert.Equal(t, updatedRoom.RoomType, queriedAfterUpdate.RoomType)
		require.Len(t, queriedAfterUpdate.Members, 2)
		assert.Equal(t, int64(1002), queriedAfterUpdate.Members[1].Id)

		err = repo.UpdateByID(ctx, roomID, updatedRoom) // Stale version
		require.Error(t, err)

		err = repo.UpdateByID(ctx, roomID+1, updatedRoom)
		require.Error(t, err)
	})
}
