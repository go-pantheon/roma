package data

import (
	"context"
	"strings"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-util/data/db/postgresql/migrate"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain"
	"github.com/go-pantheon/roma/app/player/internal/app/user/gate/domain/userregister"
	"github.com/go-pantheon/roma/app/player/internal/data/pguser"
	dbv1 "github.com/go-pantheon/roma/gen/api/db/player/v1"
	"github.com/go-pantheon/roma/pkg/data/postgresdb"
	"github.com/go-pantheon/roma/pkg/universe/life"
)

const (
	_tableName = "user"
)

var _ domain.UserRepo = (*userPostgresRepo)(nil)

type userPostgresRepo struct {
	log  *log.Helper
	data *postgresdb.DB
}

func NewUserPostgresRepo(data *postgresdb.DB, logger log.Logger) (domain.UserRepo, error) {
	return newUserPostgresRepo(data, logger, userregister.AllModuleKeys())
}

func TestNewUserPostgresRepo(data *postgresdb.DB, logger log.Logger, mods []life.ModuleKey) (domain.UserRepo, error) {
	return newUserPostgresRepo(data, logger, mods)
}

func newUserPostgresRepo(data *postgresdb.DB, logger log.Logger, _ []life.ModuleKey) (domain.UserRepo, error) {
	r := &userPostgresRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "player/user/gate/data")),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := migrate.Migrate(ctx, r.data.DB, &dbv1.UserProto{}, userregister.AllModuleDBColumnsString()); err != nil {
		return nil, err
	}

	return r, nil
}

func (r *userPostgresRepo) Create(ctx context.Context, uid int64, user *dbv1.UserProto, ctime time.Time) error {
	if user == nil {
		return errors.New("user proto is nil")
	}

	modcols, modvals, modsigns, err := pguser.ParseUpsertModuleSQLParam(user, 4)
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

	if _, err = r.data.DB.Exec(ctx, sqlbuilder.String(), vals...); err != nil {
		return errors.Wrapf(err, "inserting user %d", uid)
	}

	return nil
}

func (r *userPostgresRepo) QueryByID(ctx context.Context, uid int64, user *dbv1.UserProto, mods []life.ModuleKey) error {
	if user == nil {
		return errors.New("user proto is nil")
	}

	scanargs := make([]any, 0, len(mods)+3)
	scanargs = append(scanargs, &user.Id, &user.Version, &user.ServerVersion)

	modcols, modvals, err := pguser.ParseQueryModuleSQLParam(mods)
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

	row := r.data.DB.QueryRow(ctx, sqlbuilder.String(), uid)

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

		p, err := pguser.UnmarshalUserModule(bytes, mod)
		if err != nil {
			return errors.Wrapf(err, "unmarshaling module %s", mod)
		}

		user.Modules[string(mod)] = p
	}

	return nil
}

func (r *userPostgresRepo) UpdateByID(ctx context.Context, uid int64, user *dbv1.UserProto) error {
	if user == nil {
		return errors.New("user proto is nil")
	}

	modcols, modvals, modsigns, err := pguser.ParseUpsertModuleSQLParam(user, 4)
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

	if _, err = r.data.DB.Exec(ctx, sqlbuilder.String(), vals...); err != nil {
		return errors.Wrapf(err, "updating user %d", uid)
	}

	return nil
}

func (r *userPostgresRepo) IsExist(ctx context.Context, uid int64) (bool, error) {
	sql := `SELECT EXISTS(SELECT 1 FROM "` + _tableName + `" WHERE id = $1);`

	var exists bool
	if err := r.data.DB.QueryRow(ctx, sql, uid).Scan(&exists); err != nil {
		return false, errors.Wrapf(err, "scanning user %d existence", uid)
	}

	return exists, nil
}

func (r *userPostgresRepo) IncVersion(ctx context.Context, uid int64, newVersion int64) error {
	sql := `UPDATE "user" SET version = $1 WHERE id = $2 AND version = $3;`
	_, err := r.data.DB.Exec(ctx, sql, newVersion, uid, newVersion-1)
	if err != nil {
		return errors.Wrapf(err, "incrementing user %d version", uid)
	}

	return nil
}
