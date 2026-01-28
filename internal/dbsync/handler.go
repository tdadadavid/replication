package dbsync

import (
	"context"
	"dbreplication/internal"
	"dbreplication/internal/metrics"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"log/slog"
	"time"

	"github.com/jackc/pgx/v5"
	"golang.org/x/sync/errgroup"
)

type SyncHandler struct {
	Leader    *pgxpool.Pool
	Followers []*pgxpool.Pool
	log       *slog.Logger

	Metrics *metrics.ReplicationMetrics

	// Controls how many follower writes can run concurrently.
	// If you have 10 followers, you may not want 10 concurrent DB writes per request.
	MaxFollowerConcurrency int

	// Timeout per follower write (teaches tail latency + failure behavior)
	FollowerWriteTimeout time.Duration
}

func (h *SyncHandler) Handle(ctx context.Context, user *internal.User) (bool, error) {
	const op = "insert_user"

	args := pgx.NamedArgs{
		"email":   user.Email,
		"balance": user.Balance,
		"age":     user.Age,
	}
	query := "INSERT INTO users (email, balance, age) VALUES (@email, @balance, @age)"

	// master write.
	masterT := metrics.StartTimer()
	_, err := h.Leader.Exec(ctx, query, args)
	h.observeMaster(op, masterT.Seconds(), err)

	if err != nil {
		h.log.Info("failed to insert into leader", "error", err.Error())
		return false, errors.New("failed to insert information into leader")
	}
	h.log.Info("inserted into master successfully")

	// Replication [Follower Write]
	overheadStart := time.Now()

	// Use defaults if not set
	maxConc := h.MaxFollowerConcurrency
	if maxConc <= 0 {
		maxConc = 4 // reasonable default
	}
	timeout := h.FollowerWriteTimeout
	if timeout <= 0 {
		timeout = 2 * time.Second
	}

	// If any follower fails, errgroup cancels the derived context
	g, gctx := errgroup.WithContext(ctx)

	// semaphore to bound concurrency
	sem := make(chan struct{}, maxConc)

	for idx := range h.Followers {
		idx := idx // capture loop variable
		follower := h.Followers[idx]
		followerLabel := fmt.Sprintf("follower-%d", idx)

		g.Go(func() error {
			// acquire slot
			sem <- struct{}{}
			defer func() { <-sem }()

			// Each follower gets its own timeout.
			fCtx, cancel := context.WithTimeout(gctx, timeout)
			defer cancel()

			followerT := metrics.StartTimer()
			_, ferr := follower.Exec(fCtx, query, args)
			h.observeFollower(op, followerLabel, followerT.Seconds(), ferr)

			if ferr != nil {
				h.log.Info(
					"failed to replicate to follower",
					"follower-id", idx,
					"port", follower.Config().ConnConfig.Port,
					"error", ferr.Error(),
				)
				return fmt.Errorf("follower-%d failed: %w", idx, ferr)
			}

			h.log.Info(
				"replicated successfully",
				"follower-id", idx,
				"port", follower.Config().ConnConfig.Port,
			)
			return nil
		})
	}

	// Wait for all follower writes. If one fails, returns early and cancels others.
	if err := g.Wait(); err != nil {
		// MASTER succeeded, FOLLOWER failed => partial commit (consistency debt)
		h.observePartialCommit(op)

		// In strict sync replication, the request must fail.
		return false, fmt.Errorf("sync replication failed: %w", err)
	}

	h.observeOverhead(op, time.Since(overheadStart).Seconds())

	return true, nil
}

func (h *SyncHandler) observeMaster(op string, seconds float64, err error) {
	if h.Metrics == nil {
		return
	}
	h.Metrics.MasterWriteDuration.WithLabelValues(op).Observe(seconds)
}

func (h *SyncHandler) observeFollower(op, follower string, seconds float64, err error) {
	if h.Metrics == nil {
		return
	}

	h.Metrics.FollowerWriteDuration.WithLabelValues(op, follower).Observe(seconds)

	result := "ok"
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			result = "timeout"
			h.Metrics.FollowerTimeoutsTotal.WithLabelValues(op, follower).Inc()
		} else {
			result = "error"
		}
	}

	h.Metrics.FollowerWriteTotal.WithLabelValues(op, follower, result).Inc()
}

func (h *SyncHandler) observeOverhead(op string, seconds float64) {
	if h.Metrics == nil {
		return
	}
	h.Metrics.SyncOverheadDuration.WithLabelValues(op).Observe(seconds)
}

func (h *SyncHandler) observePartialCommit(op string) {
	if h.Metrics == nil {
		return
	}
	h.Metrics.PartialCommitTotal.WithLabelValues(op).Inc()
}
