package data

import (
	"context"
	"strconv"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-util/errors"
	basicobj "github.com/go-pantheon/roma/app/player/internal/app/basic/gate/domain/object"
	statusobj "github.com/go-pantheon/roma/app/player/internal/app/status/gate/domain/object"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/userregister"
	"github.com/go-pantheon/roma/app/player/internal/data"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

const (
	_tableName = "user"
	_jsonbType = "JSONB"
)

// specialColumnTypeMap defines specific data types for certain user modules.
// Modules not listed here will default to BYTEA (for raw bytes).
var specialColumnTypeMap = map[life.ModuleKey]string{
	basicobj.ModuleKey:  _jsonbType,
	statusobj.ModuleKey: _jsonbType,
}

var _ domain.UserRepo = (*userPostgresRepo)(nil)

type userPostgresRepo struct {
	log  *log.Helper
	data *data.Data
}

func NewUserPostgresRepo(data *data.Data, logger log.Logger) (domain.UserRepo, error) {
	return newUserPostgresRepo(data, logger, userregister.AllModuleKeys())
}

func TestNewUserPostgresRepo(data *data.Data, logger log.Logger, mods []life.ModuleKey) (domain.UserRepo, error) {
	return newUserPostgresRepo(data, logger, mods)
}

func newUserPostgresRepo(data *data.Data, logger log.Logger, mods []life.ModuleKey) (domain.UserRepo, error) {
	r := &userPostgresRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "player/user/gate/data")),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := r.initDB(ctx, mods); err != nil {
		return nil, err
	}

	return r, nil
}

type Column struct {
	Name    string
	Type    string
	NotNull bool
}

func (r *userPostgresRepo) GetTableColumns(ctx context.Context, tableName string) ([]Column, error) {
	query := `
		SELECT column_name, data_type, is_nullable
		FROM information_schema.columns
		WHERE table_schema = 'public' AND table_name = $1
		ORDER BY ordinal_position;
	`

	rows, err := r.data.Pdb.Query(ctx, query, tableName)
	if err != nil {
		return nil, errors.Wrapf(err, "querying columns for table %s", tableName)
	}
	defer rows.Close()

	var columns []Column
	for rows.Next() {
		var col Column
		var isNullable string
		if err := rows.Scan(&col.Name, &col.Type, &isNullable); err != nil {
			return nil, errors.Wrap(err, "scanning column info")
		}
		col.NotNull = isNullable == "NO"
		columns = append(columns, col)
	}

	if err := rows.Err(); err != nil {
		return nil, errors.Wrap(err, "processing rows for table columns")
	}

	return columns, nil
}

func (r *userPostgresRepo) initDB(ctx context.Context, mods []life.ModuleKey) error {
	if err := r.createTable(ctx); err != nil {
		return err
	}

	if err := r.updateColumns(ctx, mods); err != nil {
		return err
	}
	return nil
}

func (r *userPostgresRepo) createTable(ctx context.Context) error {
	const createTableSQL = `
	CREATE TABLE IF NOT EXISTS "user" (
		"id" BIGINT PRIMARY KEY,
		"version" BIGINT NOT NULL DEFAULT 0,
		"server_version" BIGINT NOT NULL DEFAULT 0
	);
	`
	_, err := r.data.Pdb.Exec(ctx, createTableSQL)
	if err != nil {
		return errors.Wrap(err, "failed to create user table")
	}

	r.log.Infof("[PostgreSQL] created user table")

	return nil
}

func (r *userPostgresRepo) updateColumns(ctx context.Context, mods []life.ModuleKey) error {
	const queryCols = `
	SELECT column_name FROM information_schema.columns
	WHERE table_schema = 'public' AND table_name = $1;
	`
	rows, err := r.data.Pdb.Query(ctx, queryCols, _tableName)
	if err != nil {
		return errors.Wrap(err, "querying user table columns")
	}
	defer rows.Close()

	existingCols := make(map[string]struct{}, len(mods))

	for rows.Next() {
		var col string
		if err := rows.Scan(&col); err != nil {
			return errors.Wrap(err, "scanning column name")
		}
		existingCols[col] = struct{}{}
	}

	for _, mod := range mods {
		col := string(mod)
		if _, exists := existingCols[col]; exists {
			continue
		}

		colType, ok := specialColumnTypeMap[mod]
		if !ok {
			colType = "BYTEA"
		}

		alterSQL := `ALTER TABLE "user" ADD COLUMN IF NOT EXISTS "` + col + `" ` + colType + `;`
		if _, err := r.data.Pdb.Exec(ctx, alterSQL); err != nil {
			return errors.Wrapf(err, "adding column %s with type %s", col, colType)
		}

		r.log.Infof("[PostgreSQL] added column '%s' with type %s in table 'user'", col, colType)
	}

	return nil
}

func (r *userPostgresRepo) Create(ctx context.Context, uid int64, user *dbv1.UserProto, ctime time.Time) error {
	if user == nil {
		return errors.New("user proto is nil")
	}

	modcols, modvalues, modsigns, err := r.parseModuleSQLParam(user, 4)
	if err != nil {
		return errors.Wrapf(err, "parsing module sql param")
	}

	values := []any{uid, user.Version, user.ServerVersion}
	values = append(values, modvalues...)

	queryBuilder := strings.Builder{}
	queryBuilder.WriteString("INSERT INTO ")
	queryBuilder.WriteString(`"` + _tableName + `"`)
	queryBuilder.WriteString(" (")

	queryBuilder.WriteString("id, version, server_version")

	if len(modcols) > 0 {
		queryBuilder.WriteString(", ")
		queryBuilder.WriteString(strings.Join(modcols, ", "))
	}

	queryBuilder.WriteString(") VALUES ($1, $2, $3")

	if len(modsigns) > 0 {
		queryBuilder.WriteString(", ")
		queryBuilder.WriteString(strings.Join(modsigns, ", "))
	}

	queryBuilder.WriteString(")")

	if _, err = r.data.Pdb.Exec(ctx, queryBuilder.String(), values...); err != nil {
		return errors.Wrapf(err, "inserting user %d", uid)
	}

	return nil
}

func (r *userPostgresRepo) QueryByID(ctx context.Context, uid int64, user *dbv1.UserProto, mods []life.ModuleKey) error {
	if user == nil {
		return errors.New("user proto is nil")
	}

	modcols := make([]string, 0, len(mods))
	modValues := make([][]byte, len(mods))
	scanArgs := make([]any, len(mods)+3)
	scanArgs[0] = &user.Id
	scanArgs[1] = &user.Version
	scanArgs[2] = &user.ServerVersion

	for i, mod := range mods {
		modcols = append(modcols, `"`+string(mod)+`"`)
		scanArgs[i+3] = &modValues[i]
	}

	querySQLBuilder := strings.Builder{}
	querySQLBuilder.WriteString("SELECT id, version, server_version")

	if len(modcols) > 0 {
		querySQLBuilder.WriteString(", ")
		querySQLBuilder.WriteString(strings.Join(modcols, ", "))
	}

	querySQLBuilder.WriteString(` FROM "`)
	querySQLBuilder.WriteString(_tableName)
	querySQLBuilder.WriteString(`" WHERE id = $1`)

	row := r.data.Pdb.QueryRow(ctx, querySQLBuilder.String(), uid)

	if err := row.Scan(scanArgs...); err != nil {
		return errors.Wrapf(err, "scanning user %d", uid)
	}

	if user.Modules == nil {
		user.Modules = make(map[string]*dbv1.UserModuleProto)
	}

	for i, mod := range mods {
		rawBytes := modValues[i]
		if len(rawBytes) == 0 {
			continue
		}

		p := &dbv1.UserModuleProto{}

		if specialColumnTypeMap[mod] == _jsonbType {
			if err := protojson.Unmarshal(rawBytes, p); err != nil {
				return errors.Wrapf(err, "unmarshaling module %s", mod)
			}
		} else {
			if err := proto.Unmarshal(rawBytes, p); err != nil {
				return errors.Wrapf(err, "unmarshaling module %s", mod)
			}
		}
		user.Modules[string(mod)] = p
	}

	return nil
}

func (r *userPostgresRepo) UpdateByID(ctx context.Context, uid int64, user *dbv1.UserProto) error {
	if user == nil {
		return errors.New("user proto is nil")
	}

	modcols, modvalues, modsigns, err := r.parseModuleSQLParam(user, 4)
	if err != nil {
		return errors.Wrapf(err, "parsing module sql param")
	}

	querySQLBuilder := strings.Builder{}
	querySQLBuilder.WriteString("UPDATE ")
	querySQLBuilder.WriteString(`"` + _tableName + `"`)
	querySQLBuilder.WriteString(" SET version = $1 ")

	for i := 0; i < len(modcols); i++ {
		querySQLBuilder.WriteString(", ")
		querySQLBuilder.WriteString(modcols[i])
		querySQLBuilder.WriteString(" = ")
		querySQLBuilder.WriteString(modsigns[i])
	}

	querySQLBuilder.WriteString(" WHERE id = $2 AND version = $3")

	values := []any{user.Version, uid, user.Version - 1}
	values = append(values, modvalues...)

	if _, err = r.data.Pdb.Exec(ctx, querySQLBuilder.String(), values...); err != nil {
		return errors.Wrapf(err, "updating user %d", uid)
	}

	return nil
}

func (r *userPostgresRepo) IsExist(ctx context.Context, uid int64) (bool, error) {
	querySQL := `SELECT EXISTS(SELECT 1 FROM "` + _tableName + `" WHERE id = $1);`

	var exists bool
	if err := r.data.Pdb.QueryRow(ctx, querySQL, uid).Scan(&exists); err != nil {
		return false, errors.Wrapf(err, "scanning user %d existence", uid)
	}

	return exists, nil
}

func (r *userPostgresRepo) IncVersion(ctx context.Context, uid int64, newVersion int64) error {
	querySQL := `UPDATE "user" SET version = $1 WHERE id = $2 AND version = $3;`
	_, err := r.data.Pdb.Exec(ctx, querySQL, newVersion, uid, newVersion-1)
	if err != nil {
		return errors.Wrapf(err, "incrementing user %d version", uid)
	}

	return nil
}

func (r *userPostgresRepo) parseModuleSQLParam(user *dbv1.UserProto, i int) (cols []string, values []any, signs []string, err error) {
	values = make([]any, 0, len(user.Modules))

	cols = make([]string, 0, len(user.Modules))
	signs = make([]string, 0, len(user.Modules))

	for k, p := range user.Modules {
		cols = append(cols, `"`+string(k)+`"`)
		signs = append(signs, "$"+strconv.Itoa(i))
		i++

		v, err := r.marshalModule(life.ModuleKey(k), p)
		if err != nil {
			return nil, nil, nil, errors.Wrapf(err, "marshaling module %s", k)
		}

		values = append(values, v)
	}

	return cols, values, signs, nil
}

func (r *userPostgresRepo) marshalModule(key life.ModuleKey, p proto.Message) (ret []byte, err error) {
	if specialColumnTypeMap[life.ModuleKey(key)] == _jsonbType {
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
