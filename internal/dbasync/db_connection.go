package dbasync

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

func mustPool(ctx context.Context, dsn string) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		panic(err)
	}
	return pool
}

func connectLeader() *pgxpool.Pool {
	ctx := context.Background()
	// postgres://<user>:<password>@<localhost>:<port>/<table_name>
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", "postgres", "postgres", "localhost", 5432, "replication")
	db := mustPool(ctx, dsn)
	err := db.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	return db
}
