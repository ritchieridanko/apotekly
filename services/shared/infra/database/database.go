package database

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type Database struct {
	pool *pgxpool.Pool
}

func NewDatabase(p *pgxpool.Pool) *Database {
	return &Database{pool: p}
}

func (d *Database) Execute(ctx context.Context, query string, args ...any) error {
	res, err := d.executor(ctx).Exec(ctx, query, args...)
	if err != nil {
		return err
	}
	if ra := res.RowsAffected(); ra == 0 {
		return ce.ErrDBAffectNoRows
	}
	return nil
}

func (d *Database) Query(ctx context.Context, query string, args ...any) pgx.Row {
	return d.executor(ctx).QueryRow(ctx, query, args...)
}

func (d *Database) QueryAll(ctx context.Context, query string, args ...any) (pgx.Rows, error) {
	return d.executor(ctx).Query(ctx, query, args...)
}

func (d *Database) WithinTx(ctx context.Context) bool {
	return fromCtx(ctx) != nil
}
