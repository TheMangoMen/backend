package handler

import (
	"encoding/json"
	"net/http"

	"github.com/TheMangoMen/backend/internal/auth"
	"github.com/TheMangoMen/backend/internal/model"
	"github.com/TheMangoMen/backend/internal/service"
)

func GetWatchedStatusCounts(as service.AnalyticsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uID := auth.MustFromContext(r.Context())
		watchedJobsStatusCounts, err := as.GetWatchedJobsStatusCounts(uID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		watchedCompaniesStatusCounts, err := as.GetWatchedCompaniesStatusCounts(uID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		res := struct {
			Jobs      []model.StatusCount `json:"jobs"`
			Companies []model.StatusCount `json:"companies"`
		}{
			Jobs:      watchedJobsStatusCounts,
			Companies: watchedCompaniesStatusCounts,
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(&res); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
