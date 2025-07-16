package postgres

import "github.com/jackc/pgx/v5/pgconn"

func IsUniqueViolationError(err error) bool {
	pgErr, ok := err.(*pgconn.PgError)
	return ok && pgErr.Code == "23505"
}
