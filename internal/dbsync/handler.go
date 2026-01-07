package dbsync

import (
	"database/sql"
	"dbreplication/internal"
)

type SyncHandler struct {
	Writer  *sql.DB   //master
	Readers []*sql.DB //slaves
}

func (h *SyncHandler) Handle(user *internal.User) bool {
	return false
}
