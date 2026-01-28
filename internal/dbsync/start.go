package dbsync

import (
	"dbreplication/internal/metrics"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
)

func Start(m *metrics.ReplicationMetrics, logger *slog.Logger) *SyncHandler {
	leader := connectLeader()
	follower1 := connectFollower1()
	follower2 := connectFollower2()
	follower3 := connectFollower3()

	return &SyncHandler{
		Leader:    leader,
		Followers: []*pgxpool.Pool{follower1, follower2, follower3},
		log:       logger,
		Metrics:   m,
	}
}
