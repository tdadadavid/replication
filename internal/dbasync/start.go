package dbasync

import (
	"github.com/jackc/pgx/v5"
)

func Start() *AsyncHandler {
	leader := connectLeader()
	follower1 := connectFollower1()
	follower2 := connectFollower2()
	follower3 := connectFollower3()

	return &AsyncHandler{
		Leader:    leader,
		Followers: []*pgx.Conn{follower1, follower2, follower3},
	}
}
