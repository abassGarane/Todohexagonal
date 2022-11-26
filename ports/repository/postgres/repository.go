package postgres

import (
	"context"
	"database/sql"
)

type postgresRepository struct {
	ctx context.Context
	db  *sql.DB
}
