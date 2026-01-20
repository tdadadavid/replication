package dbasync

import (
	"dbreplication/internal"

	"github.com/jackc/pgx/v5"
)

type AsyncHandler struct {
	Leader    *pgx.Conn   //master
	Followers []*pgx.Conn //slaves
}

func (h *AsyncHandler) Handle(users *internal.User) bool {
	return false
}
