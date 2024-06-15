package handlers

import (
	"net/http"

	"github.com/TheMangoMen/backend/internal/service"
)

func CreateUser(us service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := us.CreateUser(r.PathValue("uID")); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		}
		w.WriteHeader(http.StatusOK)
	}
}
