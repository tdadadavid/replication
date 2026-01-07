package api

import (
	"dbreplication/internal"
	"dbreplication/internal/dbasync"
	"dbreplication/internal/dbsync"
	"encoding/json"
	"fmt"
	"net/http"
)

func Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`{"status": "ok"}`))
		w.WriteHeader(http.StatusOK)
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
		ok := dbsync.Handle(user)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	})
	mux.HandleFunc("/async-users", func(w http.ResponseWriter, r *http.Request) {
		user, err := decodeRequest(r)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ok := dbasync.Handle(user)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
	})

	port := 9000
	fmt.Println("Starting server on port", port)
	http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
}

func decodeRequest(r *http.Request) (*internal.User, error) {
	var user internal.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
