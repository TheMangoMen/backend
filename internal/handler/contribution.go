package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/TheMangoMen/backend/internal/model"
	"github.com/TheMangoMen/backend/internal/service"
)

func GetContribution(cs service.ContributionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jIDStr := r.URL.Query().Get("jID")
		uID := r.URL.Query().Get("uID")
		if jIDStr == "" || uID == "" {
			http.Error(w, "missing job id or user id", http.StatusBadRequest)
			return
		}
		jID, err := strconv.Atoi(jIDStr)
		if err != nil {
			http.Error(w, "invalid job id", http.StatusBadRequest)
		}
		contribution, err := cs.GetContribution(jID, uID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(&contribution); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

func AddContribution(cs service.ContributionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := model.Contribution{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := cs.AddContribution(body); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
