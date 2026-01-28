package cmd

import (
	"context"
	"dbreplication/internal"
	"dbreplication/internal/dbasync"
	"dbreplication/internal/dbsync"
	"dbreplication/internal/metrics"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func Execute() {
	logger := slog.Default()
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	mux := http.NewServeMux()

	// Prometheus registry + metrics
	reg := prometheus.NewRegistry()
	m := metrics.NewReplicationMetrics(reg)

	syncHandler := dbsync.Start(m, logger)
	asyncHandler := dbasync.Start(logger)

	// expose /metrics
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	mux.Handle("/sync-users", m.Instrument("/sync-users", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		user, err := decodeRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Optional: Set a deadline for the entire sync replication path
		reqCtx, cancel := context.WithTimeout(ctx, 100*time.Second)
		defer cancel()

		// Your handler logic
		ok, handleErr := syncHandler.Handle(reqCtx, user)
		if !ok {
			if handleErr == nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				return
			}
			http.Error(w, handleErr.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})))

	mux.Handle("/async-users", m.Instrument("/async-users", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer r.Body.Close()

		user, err := decodeRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ok, err := asyncHandler.Handle(ctx, user)
		if !ok {
			http.Error(w, "internal server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})))

	// start server
	port := 8000
	handler := recoverMiddleware(mux)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: handler,
	}

	go func() {
		fmt.Println("Starting server on port", port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("server error:", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	fmt.Println("Shutting down server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
}

func decodeRequest(r *http.Request) (*internal.User, error) {
	var user internal.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func recoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				http.Error(w, "internal server error", http.StatusInternalServerError)
				fmt.Println("panic:", rec)
			}
		}()
		next.ServeHTTP(w, r)
	})
}
