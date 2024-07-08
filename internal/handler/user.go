package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/TheMangoMen/backend/internal/auth"
	"github.com/TheMangoMen/backend/internal/service"
)

func CreateUser(us service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := us.CreateUser(r.PathValue("uID")); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func GetUser(us service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		u := auth.MustFromContext(r.Context())
		user, err := us.GetUser(u.UID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "user not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&user); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
