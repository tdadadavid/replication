package dbasync

import (
	"log/slog"
)

func Start(logger *slog.Logger) *AsyncHandler {
	leader := connectLeader()

	return &AsyncHandler{
		Leader: leader,
		log:    logger,
	}
}
