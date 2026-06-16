package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/ritchieridanko/apotekly/services/shared/utils/ce"
)

type Transactor struct {
	pool *pgxpool.Pool
}

func NewTransactor(p *pgxpool.Pool) *Transactor {
	return &Transactor{pool: p}
}

func (t *Transactor) WithTx(ctx context.Context, fn func(context.Context) *ce.Error) *ce.Error {
	tx := fromCtx(ctx)
	isNewTx := false

	var err error
	if tx == nil {
		tx, err = t.pool.Begin(ctx)
		if err != nil {
			return ce.NewError(ce.CodeDBTx, ce.MsgInternalServer, err)
		}

		ctx = toCtx(ctx, tx)
		isNewTx = true
	}
	if err := fn(ctx); err != nil {
		if isNewTx {
			tx.Rollback(ctx)
		}
		return err
	}
	if isNewTx {
		if err := tx.Commit(ctx); err != nil {
			return ce.NewError(ce.CodeDBTx, ce.MsgInternalServer, err)
		}
	}
	return nil
}
