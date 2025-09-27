package migration

import (
	"context"
	_ "embed"
	"github.com/jackc/pgx/v5"
)

//go:embed schema.sql
var DDL string

func Run(ctx context.Context, conn *pgx.Conn) error {
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := conn.Exec(ctx, DDL); err != nil {
		return err
	}

	return tx.Commit(ctx)
}
