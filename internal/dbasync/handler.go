package dbasync

import (
	"context"
	"dbreplication/internal"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"

	"github.com/jackc/pgx/v5"
)

type AsyncHandler struct {
	Leader *pgxpool.Pool //master
	log    *slog.Logger
}

func (h *AsyncHandler) Handle(ctx context.Context, user *internal.User) (bool, error) {

	args := pgx.NamedArgs{
		"email":   user.Email,
		"balance": user.Balance,
		"age":     user.Age,
	}
	query := "INSERT INTO users (email, balance, age) VALUES (@email, @balance, @age)"

	// master write.
	_, err := h.Leader.Exec(ctx, query, args)

	if err != nil {
		h.log.Info("failed to insert into leader", "error", err.Error())
		return false, errors.New("failed to insert information into leader")
	}
	h.log.Info("inserted into master successfully")
	return true, nil
}
