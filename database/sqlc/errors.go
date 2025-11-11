package db

import (
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
)

const (
	ForeignKeyViolation string = "23503"
)

func ErrorCode(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}
