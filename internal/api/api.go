package api

import (
	"dbreplication/internal"
	"dbreplication/internal/dbasync"
	"dbreplication/internal/dbsync"
	"encoding/json"
	"net/http"
)

func Start() {
	mux := http.NewServeMux()

	mux.HandleFunc("/heathlz", func(w http.ResponseWriter, r *http.Request) {})

	mux.HandleFunc("/sync-users", func(w http.ResponseWriter, r *http.Request) {
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

	http.ListenAndServe(":8000", mux)
}

func decodeRequest(r *http.Request) (*internal.User, error) {
	var user internal.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
