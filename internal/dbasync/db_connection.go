package dbasync

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func connectLeader() *pgx.Conn {
	// postgres://<user>:<password>@<localhost>:<port>/<table_name>
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", "postgres", "rootpassword", "localhost", 4500, "replication")
	db, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		panic(err)
	}
	err = db.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	return db
}

func connectFollower1() *pgx.Conn {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", "postgres", "rootpassword", "localhost", 4500, "replication")

	db, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		panic(err)
	}
	err = db.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	return db
}

func connectFollower2() *pgx.Conn {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", "postgres", "rootpassword", "localhost", 4500, "replication")
	db, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		panic(err)
	}
	err = db.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	return db
}

func connectFollower3() *pgx.Conn {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", "postgres", "rootpassword", "localhost", 4500, "replication")

	db, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		panic(err)
	}
	err = db.Ping(context.Background())
	if err != nil {
		panic(err)
	}

	return db
}
