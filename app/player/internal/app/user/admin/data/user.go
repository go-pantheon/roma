package data

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-pantheon/fabrica-util/errors"
	"github.com/go-pantheon/roma/app/player/internal/app/user/admin/domain"
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
	r := &userPgRepo{
		data: data,
		log:  log.NewHelper(log.With(logger, "module", "player/user/gate/data")),
	}

	return r, nil
}

func (r *userPgRepo) GetByID(ctx context.Context, user *dbv1.UserProto, mods []life.ModuleKey) (err error) {
	if user == nil {
		return errors.New("user proto is nil")
	}

	if user.Id == 0 {
		return errors.New("user id is zero")
	}

	scanargs := make([]any, 0, len(mods)+3)
	scanargs = append(scanargs, &user.Id, &user.Version, &user.ServerVersion)

	modcols, modvals, err := data.ParseQueryModuleSQLParam(mods)
	if err != nil {
		return errors.Wrapf(err, "parsing module sql param for user %d", user.Id)
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

	row := r.data.Pdb.QueryRow(ctx, sqlbuilder.String(), user.Id)

	return data.ScanUserRow(row, user, mods, scanargs, modvals)
}

func (r *userPgRepo) GetList(ctx context.Context, start, limit int64, conds map[life.ModuleKey]*dbv1.UserModuleProto, mods []life.ModuleKey) (result []*dbv1.UserProto, count int64, err error) {
	wheresql, args, err := data.BuildUserWhereSQL(conds)
	if err != nil {
		return nil, 0, errors.Wrap(err, "building where sql")
	}

	countsql := fmt.Sprintf(`SELECT COUNT(*) FROM "%s"%s`, _tableName, wheresql)
	if err := r.data.Pdb.QueryRow(ctx, countsql, args...).Scan(&count); err != nil {
		return nil, 0, errors.Wrap(err, "counting users")
	}

	if count == 0 {
		return []*dbv1.UserProto{}, 0, nil
	}

	colNames := []string{`"id"`, `"version"`, `"server_version"`}
	for _, mod := range mods {
		colNames = append(colNames, `"`+string(mod)+`"`)
	}

	querysql := fmt.Sprintf(`SELECT %s FROM "%s" %s ORDER BY id ASC LIMIT $%d OFFSET $%d`,
		strings.Join(colNames, ", "), _tableName, wheresql, len(args)+1, len(args)+2)

	args = append(args, limit, start)

	rows, err := r.data.Pdb.Query(ctx, querysql, args...)
	if err != nil {
		return nil, 0, errors.Wrap(err, "querying user list")
	}

	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, 0, errors.Wrap(err, "processing user list rows")
	}

	result = make([]*dbv1.UserProto, 0, limit)
	for rows.Next() {
		user := dbv1.UserProtoPool.Get()

		scanargs := []any{&user.Id, &user.Version, &user.ServerVersion}
		modvals := make([][]byte, 0, len(mods))

		for i := range mods {
			modvals = append(modvals, make([]byte, 0))
			scanargs = append(scanargs, &modvals[i])
		}

		if err := data.ScanUserRow(rows, user, mods, scanargs, modvals); err != nil {
			return nil, 0, errors.Wrap(err, "scanning user row")
		}

		result = append(result, user)
	}

	return result, count, nil
}
