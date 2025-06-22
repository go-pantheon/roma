//go:build integration

package data_test

import (
	"context"
	"os"
	"testing"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-util/data/db/mongo"
	"github.com/go-pantheon/roma/app/room/internal/app/room/admin/data"
	"github.com/go-pantheon/roma/app/room/internal/app/room/admin/domain"
	pdata "github.com/go-pantheon/roma/app/room/internal/data"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/room/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAdminRoomTestRepo(t *testing.T) (domain.RoomRepo, func()) {
	t.Helper()

	dbName := "room_admin_test"
	dbsn := "127.0.0.1:27017"

	mdb, mdbCleanup, err := mongo.New(context.Background(), dbsn, dbName)
	require.NoError(t, err)

	err = mdb.Drop(context.Background())
	require.NoError(t, err)

	// Seed data using the DB-level proto
	roomToSeed := &dbv1.RoomProto{
		Id:       3001,
		Sid:      1,
		RoomType: 1,
		Version:  1,
	}
	_, err = mdb.Collection("room").InsertOne(context.Background(), roomToSeed)
	require.NoError(t, err)

	d := &pdata.Data{Mdb: mdb}
	logger := log.NewStdLogger(os.Stdout)

	repo, err := data.NewRoomMongoRepo(d, logger)
	require.NoError(t, err)

	return repo, mdbCleanup
}

func TestAdminRoomMongoRepo(t *testing.T) {
	repo, cleanup := setupAdminRoomTestRepo(t)
	defer cleanup()

	ctx := context.Background()
	roomID := int64(3001)

	t.Run("GetByID", func(t *testing.T) {
		// Test get existing room
		room, err := repo.GetByID(ctx, roomID)
		require.NoError(t, err)
		require.NotNil(t, room)
		assert.Equal(t, roomID, room.Id)
		assert.Equal(t, uint64(1), room.Sid)
		assert.Equal(t, uint64(1), room.RoomType)
		assert.Equal(t, int64(1), room.Version)

		// Test get non-existent room
		_, err = repo.GetByID(ctx, roomID+1)
		require.Error(t, err)

		// Test get with zero ID
		_, err = repo.GetByID(ctx, 0)
		require.Error(t, err)
	})
}
