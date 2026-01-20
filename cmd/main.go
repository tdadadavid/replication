package cmd

import (
	"context"
	"dbreplication/internal"
	"dbreplication/internal/dbasync"
	"dbreplication/internal/dbsync"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
)

func Execute() {
	ctx := context.Background()
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, os.Kill)

	// register otel spans and other neccessary things for observability
	syncHandler := dbsync.Start()
	asyncHandler := dbasync.Start()

	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	mux.HandleFunc("/sync-users", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		user, err := decodeRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ok, err := syncHandler.Handle(ctx, user)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/async-users", func(w http.ResponseWriter, r *http.Request) {
		user, err := decodeRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ok := asyncHandler.Handle(user)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	})

	port := 9000
	fmt.Println("Starting server on port", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
	<-ctx.Done()
	fmt.Println("Closing server on port", port)
	cancel()
}

func decodeRequest(r *http.Request) (*internal.User, error) {
	var user internal.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
