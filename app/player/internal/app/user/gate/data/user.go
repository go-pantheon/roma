package data

import (
	"context"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/userregister"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/userregister/modulecolumn"
	"github.com/go-pantheon/roma/app/player/internal/data"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/universe/life"
)

const (
	_tableName = "user"
)

var _ domain.UserRepo = (*userPgRepo)(nil)

type userPgRepo struct {
	log  *log.Helper
	data *data.Data
}

func NewUserPgRepo(data *data.Data, logger log.Logger) (domain.UserRepo, error) {
	return newUserPgRepo(data, logger, userregister.AllModuleKeys())
}

func TestNewUserPgRepo(data *data.Data, logger log.Logger, mods []life.ModuleKey) (domain.UserRepo, error) {
	return newUserPgRepo(data, logger, mods)
}

func newUserPgRepo(data *data.Data, logger log.Logger, mods []life.ModuleKey) (domain.UserRepo, error) {
	r := &userPgRepo{
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

func (r *userPgRepo) GetTableColumns(ctx context.Context, tableName string) ([]Column, error) {
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

func (r *userPgRepo) initDB(ctx context.Context, mods []life.ModuleKey) error {
	if err := r.createTable(ctx); err != nil {
		return err
	}

	if err := r.updateColumns(ctx, mods); err != nil {
		return err
	}
	return nil
}

func (r *userPgRepo) createTable(ctx context.Context) error {
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

func (r *userPgRepo) updateColumns(ctx context.Context, mods []life.ModuleKey) error {
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

		colType, ok := modulecolumn.ColumnTypeMap[mod]
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

func (r *userPgRepo) Create(ctx context.Context, uid int64, user *dbv1.UserProto, ctime time.Time) error {
	if user == nil {
		return errors.New("user proto is nil")
	}

	modcols, modvals, modsigns, err := data.ParseUpsertModuleSQLParam(user, 4)
	if err != nil {
		return errors.Wrapf(err, "parsing module sql param")
	}

	vals := []any{uid, user.Version, user.ServerVersion}
	vals = append(vals, modvals...)

	sqlbuilder := strings.Builder{}
	sqlbuilder.WriteString("INSERT INTO ")
	sqlbuilder.WriteString(`"` + _tableName + `"`)
	sqlbuilder.WriteString(" (")

	sqlbuilder.WriteString("id, version, server_version")

	if len(modcols) > 0 {
		sqlbuilder.WriteString(", ")
		sqlbuilder.WriteString(strings.Join(modcols, ", "))
	}

	sqlbuilder.WriteString(") VALUES ($1, $2, $3")

	if len(modsigns) > 0 {
		sqlbuilder.WriteString(", ")
		sqlbuilder.WriteString(strings.Join(modsigns, ", "))
	}

	sqlbuilder.WriteString(")")

	if _, err = r.data.Pdb.Exec(ctx, sqlbuilder.String(), vals...); err != nil {
		return errors.Wrapf(err, "inserting user %d", uid)
	}

	return nil
}

func (r *userPgRepo) QueryByID(ctx context.Context, uid int64, user *dbv1.UserProto, mods []life.ModuleKey) error {
	if user == nil {
		return errors.New("user proto is nil")
	}

	scanargs := make([]any, 0, len(mods)+3)
	scanargs = append(scanargs, &user.Id, &user.Version, &user.ServerVersion)

	modcols, modvals, err := data.ParseQueryModuleSQLParam(mods)
	if err != nil {
		return errors.Wrapf(err, "parsing module sql param")
	}

	for i := range modvals {
		scanargs = append(scanargs, &modvals[i])
	}

	sqlbuilder := strings.Builder{}
	sqlbuilder.WriteString("SELECT id, version, server_version")

	if len(modcols) > 0 {
		sqlbuilder.WriteString(", ")
		sqlbuilder.WriteString(strings.Join(modcols, ", "))
	}

	sqlbuilder.WriteString(` FROM "`)
	sqlbuilder.WriteString(_tableName)
	sqlbuilder.WriteString(`" WHERE id = $1`)

	row := r.data.Pdb.QueryRow(ctx, sqlbuilder.String(), uid)

	if err := row.Scan(scanargs...); err != nil {
		return errors.Wrapf(err, "scanning user %d", uid)
	}

	if user.Modules == nil {
		user.Modules = make(map[string]*dbv1.UserModuleProto)
	}

	for i, mod := range mods {
		bytes := modvals[i]
		if len(bytes) == 0 {
			continue
		}

		p, err := data.UnmarshalUserModule(bytes, mod)
		if err != nil {
			return errors.Wrapf(err, "unmarshaling module %s", mod)
		}

		user.Modules[string(mod)] = p
	}

	return nil
}

func (r *userPgRepo) UpdateByID(ctx context.Context, uid int64, user *dbv1.UserProto) error {
	if user == nil {
		return errors.New("user proto is nil")
	}

	modcols, modvals, modsigns, err := data.ParseUpsertModuleSQLParam(user, 4)
	if err != nil {
		return errors.Wrapf(err, "parsing module sql param")
	}

	sqlbuilder := strings.Builder{}
	sqlbuilder.WriteString("UPDATE ")
	sqlbuilder.WriteString(`"` + _tableName + `"`)
	sqlbuilder.WriteString(" SET version = $1 ")

	for i := range modcols {
		sqlbuilder.WriteString(", ")
		sqlbuilder.WriteString(modcols[i])
		sqlbuilder.WriteString(" = ")
		sqlbuilder.WriteString(modsigns[i])
	}

	sqlbuilder.WriteString(" WHERE id = $2 AND version = $3")

	vals := []any{user.Version, uid, user.Version - 1}
	vals = append(vals, modvals...)

	if _, err = r.data.Pdb.Exec(ctx, sqlbuilder.String(), vals...); err != nil {
		return errors.Wrapf(err, "updating user %d", uid)
	}

	return nil
}

func (r *userPgRepo) IsExist(ctx context.Context, uid int64) (bool, error) {
	sql := `SELECT EXISTS(SELECT 1 FROM "` + _tableName + `" WHERE id = $1);`

	var exists bool
	if err := r.data.Pdb.QueryRow(ctx, sql, uid).Scan(&exists); err != nil {
		return false, errors.Wrapf(err, "scanning user %d existence", uid)
	}

	return exists, nil
}

func (r *userPgRepo) IncVersion(ctx context.Context, uid int64, newVersion int64) error {
	sql := `UPDATE "user" SET version = $1 WHERE id = $2 AND version = $3;`
	_, err := r.data.Pdb.Exec(ctx, sql, newVersion, uid, newVersion-1)
	if err != nil {
		return errors.Wrapf(err, "incrementing user %d version", uid)
	}

	return nil
}
