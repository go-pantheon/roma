//go:build integration

package data_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-util/data/db/mongo"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain"
	"github.com/go-pantheon/roma/app/player/internal/data"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestRepo connects to a test MongoDB instance, creates a repo, and provides a cleanup function.
// It ensures the database is clean for each test run.
func setupTestRepo(t *testing.T) (domain.UserRepo, func()) {
	t.Helper()

	// NOTE: This requires a running MongoDB instance at "127.0.0.1:27017".
	dbName := "player_test"
	dbsn := "127.0.0.1:27017"

	// Connect to mongo
	mdb, mdbCleanup, err := mongo.New(context.Background(), dbsn, dbName)
	require.NoError(t, err)

	// Drop the database to ensure a clean state
	err = mdb.Drop(context.Background())
	require.NoError(t, err)

	d := &data.Data{Mdb: mdb}
	logger := log.NewStdLogger(os.Stdout)

	repo, err := NewUserMongoRepo(d, logger)
	require.NoError(t, err)

	return repo, mdbCleanup
}

func TestUserMongoRepo(t *testing.T) {
	repo, cleanup := setupTestRepo(t)
	defer cleanup()

	ctx := context.Background()
	uid := int64(1001)
	ctime := time.Now()

	defaultUser := &dbv1.UserProto{
		Version:       1,
		ServerVersion: "1.0.0",
		Modules: map[string]*dbv1.UserModuleProto{
			"basic": {
				Module: &dbv1.UserModuleProto_Basic{
					Basic: &dbv1.UserBasicProto{Name: "test-user"},
				},
			},
			"system": {
				Module: &dbv1.UserModuleProto_System{
					System: &dbv1.UserSystemProto{CurrentGenId: 10},
				},
			},
		},
	}

	t.Run("Create", func(t *testing.T) {
		// Test successful creation
		err := repo.Create(ctx, uid, defaultUser, ctime)
		require.NoError(t, err)

		// Test creating a user that already exists
		err = repo.Create(ctx, uid, defaultUser, ctime)
		require.Error(t, err) // Should fail due to unique index

		// Test creating with nil user proto
		err = repo.Create(ctx, uid+1, nil, ctime)
		require.Error(t, err)
	})

	t.Run("IsExist", func(t *testing.T) {
		// Test existing user
		exists, err := repo.IsExist(ctx, uid)
		require.NoError(t, err)
		assert.True(t, exists)

		// Test non-existent user
		exists, err = repo.IsExist(ctx, uid+1)
		require.NoError(t, err)
		assert.False(t, exists)
	})

	t.Run("QueryByID", func(t *testing.T) {
		// Test full query
		queriedUser := &dbv1.UserProto{}
		err := repo.QueryByID(ctx, uid, queriedUser, nil)
		require.NoError(t, err)
		assert.Equal(t, uid, queriedUser.Id)
		assert.Equal(t, defaultUser.Version, queriedUser.Version)
		assert.Equal(t, defaultUser.ServerVersion, queriedUser.ServerVersion)
		assert.Equal(t, defaultUser.Modules["basic"].GetBasic().Name, queriedUser.Modules["basic"].GetBasic().Name)
		assert.Equal(t, defaultUser.Modules["system"].GetSystem().CurrentGenId, queriedUser.Modules["system"].GetSystem().CurrentGenId)

		// Test partial query with projection
		queriedUserPartial := &dbv1.UserProto{}
		mods := []life.ModuleKey{"basic"}
		err = repo.QueryByID(ctx, uid, queriedUserPartial, mods)
		require.NoError(t, err)
		assert.Equal(t, uid, queriedUserPartial.Id)
		assert.Contains(t, queriedUserPartial.Modules, "basic")
		assert.NotContains(t, queriedUserPartial.Modules, "system")

		// Test querying a non-existent user
		err = repo.QueryByID(ctx, uid+1, &dbv1.UserProto{}, nil)
		require.Error(t, err)
	})

	t.Run("IncVersion", func(t *testing.T) {
		// Test successful increment
		newVersion := defaultUser.Version + 1
		err := repo.IncVersion(ctx, uid, newVersion)
		require.NoError(t, err)

		queriedUser := &dbv1.UserProto{}
		err = repo.QueryByID(ctx, uid, queriedUser, nil)
		require.NoError(t, err)
		assert.Equal(t, newVersion, queriedUser.Version)

		// Test increment with wrong version
		err = repo.IncVersion(ctx, uid, newVersion)
		require.Error(t, err)

		// Test increment on non-existent user
		err = repo.IncVersion(ctx, uid+1, 2)
		require.Error(t, err)
	})

	t.Run("UpdateByID", func(t *testing.T) {
		// First, get the current user state
		currentUser := &dbv1.UserProto{}
		err := repo.QueryByID(ctx, uid, currentUser, nil)
		require.NoError(t, err)

		// Prepare an update
		updatedUser := &dbv1.UserProto{
			Version: currentUser.Version + 1, // Correct next version
			Modules: map[string]*dbv1.UserModuleProto{
				"basic": { // Update existing module
					Module: &dbv1.UserModuleProto_Basic{
						Basic: &dbv1.UserBasicProto{Name: "updated-name"},
					},
				},
				"dev": { // Add a new module
					Module: &dbv1.UserModuleProto_Dev{
						Dev: &dbv1.UserDevProto{TimeOffset: 3600},
					},
				},
			},
		}

		// Test successful update
		err = repo.UpdateByID(ctx, uid, updatedUser)
		require.NoError(t, err)

		// Verify the update
		queriedAfterUpdate := &dbv1.UserProto{}
		err = repo.QueryByID(ctx, uid, queriedAfterUpdate, nil)
		require.NoError(t, err)
		assert.Equal(t, updatedUser.Version, queriedAfterUpdate.Version)
		assert.Equal(t, updatedUser.Modules["basic"].GetBasic().Name, queriedAfterUpdate.Modules["basic"].GetBasic().Name)
		assert.Equal(t, updatedUser.Modules["dev"].GetDev().TimeOffset, queriedAfterUpdate.Modules["dev"].GetDev().TimeOffset)
		// The system module should still exist as we do partial updates with $set
		assert.Contains(t, queriedAfterUpdate.Modules, "system")

		// Test update with wrong version
		err = repo.UpdateByID(ctx, uid, updatedUser) // Version is now stale
		require.Error(t, err)

		// Test update for non-existent user
		err = repo.UpdateByID(ctx, uid+1, updatedUser)
		require.Error(t, err)

		// Test update with nil user
		err = repo.UpdateByID(ctx, uid, nil)
		require.Error(t, err)
	})
}
