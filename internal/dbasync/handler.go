package dbasync

import (
	"database/sql"
	"dbreplication/internal"
)

type AsyncHandler struct {
	Writer  *sql.DB   //master
	Readers []*sql.DB //slaves
}

func (h *AsyncHandler) Handle(users *internal.User) bool {
	return false
}
