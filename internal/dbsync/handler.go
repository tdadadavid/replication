package dbsync

import (
	"context"
	"dbreplication/internal"
	"errors"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type SyncHandler struct {
	Leader    *pgx.Conn   //master
	Followers []*pgx.Conn //slaves
	log       *slog.Logger
}

// For synchronous replication we are going to first write the master db
// then write to the followers db, we are going to ensure all writes are complete before sending a response
// otel spans are registered to watch performance of requests and db
func (h *SyncHandler) Handle(ctx context.Context, user *internal.User) (bool, error) {

	// 1. Write to master db
	args := pgx.NamedArgs{
		"email":   user.Email,
		"balance": user.Balance,
		"age":     user.Age,
	}
	query := "INSERT INTO users (email, balance, age) VALUES (@email, @balance, @age)"
	_, err := h.Leader.Exec(ctx, query, args)
	if err != nil {
		h.log.Info("failed to insert into leader", "error", err.Error())
		return false, errors.New("failed to insert information into leader")
	}

	// 2. Write to each follower db
	for idx, follower := range h.Followers {
		_, err := follower.Exec(ctx, query, args)
		if err != nil {
			h.log.Info("failed to replicate data. stopping all operations", "follower-id", idx, "host", follower.Config().Host)
			return false, fmt.Errorf("failed to replicate data to follower-%d, conn-info=%s. stopping all operations", idx, follower.Config().Host)
		}
	}

	return true, nil
}
