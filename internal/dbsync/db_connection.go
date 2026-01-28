package dbsync

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

func connectFollower1() *pgxpool.Pool {
	ctx := context.Background()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", "postgres", "postgres", "localhost", 5433, "replication")
	db := mustPool(ctx, dsn)
	err := db.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	return db
}

func connectFollower2() *pgxpool.Pool {
	ctx := context.Background()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", "postgres", "postgres", "localhost", 5434, "replication")
	db := mustPool(ctx, dsn)
	err := db.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	return db
}

func connectFollower3() *pgxpool.Pool {
	ctx := context.Background()
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", "postgres", "postgres", "localhost", 5435, "replication")

	db := mustPool(ctx, dsn)
	err := db.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	return db
}
