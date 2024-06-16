package handler

import (
	"encoding/json"
	"net/http"

	"github.com/TheMangoMen/backend/internal/model"
	"github.com/TheMangoMen/backend/internal/service"
)

type GetJobsBody struct {
	Jobs []model.Job `json:"jobs"`
}

func GetJobs(js service.JobService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uID := r.URL.Query().Get("uID")
		isRankingStage, err := js.GetIsRankingStage()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var jobs interface{}
		if !isRankingStage {
			jobs, err = js.GetJobInterviews(uID)
		} else {
			jobs, err = js.GetJobRankings(uID)
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(&jobs); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
