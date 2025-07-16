package postgres

import (
	"context"
	"github.com/jackc/pgx/v5"
	"github.com/mkolibaba/gophermart/postgres/gen"
)

type DBX struct {
	*postgres.Queries
	conn *pgx.Conn
}

func NewDBX(conn *pgx.Conn) *DBX {
	return &DBX{
		Queries: postgres.New(conn),
		conn:    conn,
	}
}

func (q *DBX) DoInTx(ctx context.Context, fn func(qtx postgres.Querier) error) error {
	tx, err := q.conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	qtx := q.Queries.WithTx(tx)

	if err = fn(qtx); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
