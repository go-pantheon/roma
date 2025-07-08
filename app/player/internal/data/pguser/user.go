package pguser

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/go-pantheon/fabrica-kit/xerrors"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/userregister"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/data/xpg"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"github.com/jackc/pgx/v5"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

func BuildUserWhereSQL(conds map[life.ModuleKey]*dbv1.UserModuleProto) (sql string, args []any, err error) {
	i := 1
	clauses := make([]string, 0, len(conds))
	args = make([]any, 0, len(conds))

	for mod, cond := range conds {
		jsonb, err := MarshalUserModule(mod, cond)
		if err != nil {
			return "", nil, errors.Wrap(err, "marshaling cond")
		}

		if string(jsonb) == "{}" {
			continue
		}

		clauses = append(clauses, fmt.Sprintf(`"%s" @> $%d::jsonb`, mod, i))
		args = append(args, jsonb)
		i++
	}

	if len(clauses) > 0 {
		sql = " WHERE " + strings.Join(clauses, " AND ")
	}

	return sql, args, nil
}

func ScanUserRow(row pgx.Row, user *dbv1.UserProto, mods []life.ModuleKey, scanargs []any, modvals [][]byte) error {
	if err := row.Scan(scanargs...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return xerrors.ErrDBRecordNotFound
		}

		return errors.Wrapf(err, "scanning user %d", user.Id)
	}

	if user.Modules == nil {
		user.Modules = make(map[string]*dbv1.UserModuleProto)
	}

	for i, mod := range mods {
		bytes := modvals[i]
		if len(bytes) == 0 {
			continue
		}

		p, err := UnmarshalUserModule(bytes, mod)
		if err != nil {
			return errors.Wrapf(err, "unmarshaling module %s", mod)
		}

		user.Modules[string(mod)] = p
	}

	return nil
}

func ParseUpsertModuleSQLParam(user *dbv1.UserProto, firstModIndex int) (cols []string, values []any, signs []string, err error) {
	values = make([]any, 0, len(user.Modules))

	cols = make([]string, 0, len(user.Modules))
	signs = make([]string, 0, len(user.Modules))

	for k, p := range user.Modules {
		cols = append(cols, `"`+k+`"`)
		signs = append(signs, "$"+strconv.Itoa(firstModIndex))
		firstModIndex++

		v, err := MarshalUserModule(life.ModuleKey(k), p)
		if err != nil {
			return nil, nil, nil, errors.Wrapf(err, "marshaling module %s", k)
		}

		values = append(values, v)
	}

	return cols, values, signs, nil
}

func ParseQueryModuleSQLParam(mods []life.ModuleKey) (cols []string, values [][]byte, err error) {
	cols = make([]string, 0, len(mods))
	values = make([][]byte, len(mods))

	for i, mod := range mods {
		cols = append(cols, `"`+string(mod)+`"`)
		values[i] = make([]byte, 0)
	}

	return cols, values, nil
}

// UnmarshalUserModule unmarshals a user module from raw bytes.
// UserModuleProto should be returned to the pool after use by calling dbv1.UserModuleProtoPool.Put.
func UnmarshalUserModule(rawBytes []byte, modKey life.ModuleKey) (*dbv1.UserModuleProto, error) {
	p := dbv1.UserModuleProtoPool.Get()

	if userregister.GetPGColumnType(modKey) == xpg.JSONType {
		if err := protojson.Unmarshal(rawBytes, p); err != nil {
			return nil, errors.Wrapf(err, "unmarshaling jsonb module %s", modKey)
		}
	} else {
		if err := proto.Unmarshal(rawBytes, p); err != nil {
			return nil, errors.Wrapf(err, "unmarshaling bytes module %s", modKey)
		}
	}

	return p, nil
}

func MarshalUserModule(key life.ModuleKey, p proto.Message) (ret []byte, err error) {
	if userregister.GetPGColumnType(key) == xpg.JSONType {
		jsonb, err := protojson.Marshal(p)
		if err != nil {
			return nil, errors.Wrapf(err, "marshaling jsonb module %s", key)
		}

		return jsonb, nil
	}

	bytes, err := proto.Marshal(p)
	if err != nil {
		return nil, errors.Wrapf(err, "marshaling bytes module %s", key)
	}

	return bytes, nil
}
