package dbsync

import (
	"github.com/jackc/pgx/v5"
)

func Start() *SyncHandler {
	leader := connectLeader()
	follower1 := connectFollower1()
	follower2 := connectFollower2()
	follower3 := connectFollower3()

	return &SyncHandler{
		Leader:    leader,
		Followers: []*pgx.Conn{follower1, follower2, follower3},
	}
}
