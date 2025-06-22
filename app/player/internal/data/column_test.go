package data

import (
	"reflect"
	"testing"

	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/userregister/modulecolumn"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

// mockPgxRow is a mock for the pgx.Row interface.
type mockPgxRow struct {
	values []any
	err    error
}

// Scan mocks the Scan method of pgx.Row.
func (m *mockPgxRow) Scan(dest ...any) error {
	if m.err != nil {
		return m.err
	}

	if len(dest) != len(m.values) {
		return errors.Errorf("expected %d destination arguments, got %d", len(m.values), len(dest))
	}

	for i, v := range m.values {
		if err := convertAssign(dest[i], v); err != nil {
			return err
		}
	}

	return nil
}

func convertAssign(dest, src any) error {
	destVal := reflect.ValueOf(dest)
	if destVal.Kind() != reflect.Ptr {
		return errors.New("destination not a pointer")
	}

	srcVal := reflect.ValueOf(src)

	// Handle special case for []byte
	if destVal.Elem().Kind() == reflect.Slice && destVal.Elem().Type().Elem().Kind() == reflect.Uint8 {
		if s, ok := src.([]byte); ok {
			destVal.Elem().Set(reflect.ValueOf(s))
			return nil
		}
	}

	if destVal.Elem().Type() == srcVal.Type() {
		destVal.Elem().Set(srcVal)
		return nil
	}

	return errors.Errorf("cannot assign %T to %T", src, dest)
}

//nolint:paralleltest
func TestBuildUserWhereSQL(t *testing.T) {
	modulecolumn.ColumnTypeMap = map[life.ModuleKey]string{
		"module1": "JSONB",
		"module2": "BYTES",
		"module3": "JSONB",
	}

	//nolint:paralleltest
	t.Run("non-empty conditions", func(t *testing.T) {
		conds := map[life.ModuleKey]*dbv1.UserModuleProto{
			"module1": {Module: &dbv1.UserModuleProto_Basic{Basic: &dbv1.UserBasicProto{}}},
			"module3": {Module: &dbv1.UserModuleProto_Basic{Basic: &dbv1.UserBasicProto{}}},
		}

		sql, args, err := BuildUserWhereSQL(conds)
		require.NoError(t, err)
		assert.Len(t, args, 2)

		expectedSQL1 := ` WHERE "module1" @> $1::jsonb AND "module3" @> $2::jsonb`
		expectedSQL2 := ` WHERE "module3" @> $1::jsonb AND "module1" @> $2::jsonb`
		assert.Contains(t, []string{expectedSQL1, expectedSQL2}, sql)
	})

	//nolint:paralleltest
	t.Run("empty conditions", func(t *testing.T) {
		conds := map[life.ModuleKey]*dbv1.UserModuleProto{}
		sql, args, err := BuildUserWhereSQL(conds)
		require.NoError(t, err)
		assert.Empty(t, sql)
		assert.Empty(t, args)
	})

	//nolint:paralleltest
	t.Run("condition with empty proto", func(t *testing.T) {
		conds := map[life.ModuleKey]*dbv1.UserModuleProto{
			"module1": {}, // This will be marshaled to "{}" and skipped.
			"module3": {Module: &dbv1.UserModuleProto_Basic{Basic: &dbv1.UserBasicProto{}}},
		}

		sql, args, err := BuildUserWhereSQL(conds)
		require.NoError(t, err)
		assert.Equal(t, ` WHERE "module3" @> $1::jsonb`, sql)
		assert.Len(t, args, 1)
	})
}

//nolint:paralleltest
func TestParseUpsertModuleSQLParam(t *testing.T) {
	modulecolumn.ColumnTypeMap = map[life.ModuleKey]string{
		"module1": "JSONB",
		"module2": "BYTES",
	}

	user := &dbv1.UserProto{
		Modules: map[string]*dbv1.UserModuleProto{
			"module1": {Module: &dbv1.UserModuleProto_Basic{Basic: &dbv1.UserBasicProto{}}},
			"module2": {Module: &dbv1.UserModuleProto_Basic{Basic: &dbv1.UserBasicProto{}}},
		},
	}

	cols, values, signs, err := ParseUpsertModuleSQLParam(user, 1)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{`"module1"`, `"module2"`}, cols)
	assert.Len(t, values, 2)
	assert.ElementsMatch(t, []string{"$1", "$2"}, signs)
}

//nolint:paralleltest
func TestParseQueryModuleSQLParam(t *testing.T) {
	mods := []life.ModuleKey{"module1", "module2"}
	cols, values, err := ParseQueryModuleSQLParam(mods)
	require.NoError(t, err)

	assert.Equal(t, []string{`"module1"`, `"module2"`}, cols)
	assert.Len(t, values, 2)
}

//nolint:paralleltest
func TestMarshalUnmarshalUserModule(t *testing.T) {
	modulecolumn.ColumnTypeMap = map[life.ModuleKey]string{
		"jsonb_module": "JSONB",
		"bytes_module": "BYTES",
	}

	//nolint:paralleltest
	t.Run("jsonb module", func(t *testing.T) {
		p := dbv1.UserModuleProtoPool.Get()
		p.Module = &dbv1.UserModuleProto_Basic{Basic: &dbv1.UserBasicProto{}}
		defer dbv1.UserModuleProtoPool.Put(p)

		rawBytes, err := MarshalUserModule("jsonb_module", p)
		require.NoError(t, err)
		assert.NotEmpty(t, rawBytes)

		unmarshaled, err := UnmarshalUserModule(rawBytes, "jsonb_module")
		require.NoError(t, err)
		defer dbv1.UserModuleProtoPool.Put(unmarshaled)

		assert.True(t, proto.Equal(p, unmarshaled))
	})

	//nolint:paralleltest
	t.Run("bytes module", func(t *testing.T) {
		p := dbv1.UserModuleProtoPool.Get()
		p.Module = &dbv1.UserModuleProto_Basic{Basic: &dbv1.UserBasicProto{}}
		defer dbv1.UserModuleProtoPool.Put(p)

		rawBytes, err := MarshalUserModule("bytes_module", p)
		require.NoError(t, err)
		assert.NotEmpty(t, rawBytes)

		unmarshaled, err := UnmarshalUserModule(rawBytes, "bytes_module")
		require.NoError(t, err)
		defer dbv1.UserModuleProtoPool.Put(unmarshaled)

		assert.True(t, proto.Equal(p, unmarshaled))
	})
}

//nolint:paralleltest
func TestScanUserRow(t *testing.T) {
	user := &dbv1.UserProto{Id: 1}
	mods := []life.ModuleKey{"module1"}
	modvals := make([][]byte, 1)
	scanargs := make([]any, 1)
	scanargs[0] = &modvals[0]

	//nolint:paralleltest
	t.Run("successful scan", func(t *testing.T) {
		modulecolumn.ColumnTypeMap = map[life.ModuleKey]string{"module1": "JSONB"}
		p := dbv1.UserModuleProtoPool.Get()
		p.Module = &dbv1.UserModuleProto_Basic{Basic: &dbv1.UserBasicProto{}}
		defer dbv1.UserModuleProtoPool.Put(p)

		marshaled, err := MarshalUserModule("module1", p)
		require.NoError(t, err)

		row := &mockPgxRow{
			values: []any{marshaled},
		}

		err = ScanUserRow(row, user, mods, scanargs, modvals)
		require.NoError(t, err)

		require.NotNil(t, user.Modules)
		require.Contains(t, user.Modules, "module1")
		assert.True(t, proto.Equal(p, user.Modules["module1"]))
	})

	//nolint:paralleltest
	t.Run("scan error", func(t *testing.T) {
		row := &mockPgxRow{err: errors.New("scan failed")}

		err := ScanUserRow(row, user, mods, scanargs, modvals)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "scanning user 1")
	})

	//nolint:paralleltest
	t.Run("unmarshal error", func(t *testing.T) {
		modulecolumn.ColumnTypeMap = map[life.ModuleKey]string{"module1": "JSONB"}
		row := &mockPgxRow{
			values: []any{[]byte("invalid json")},
		}

		err := ScanUserRow(row, user, mods, scanargs, modvals)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "unmarshaling module module1")
	})
}
