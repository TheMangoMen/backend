package handler

import (
	"encoding/json"
	"net/http"

	"github.com/TheMangoMen/backend/internal/auth"
	"github.com/TheMangoMen/backend/internal/model"
	"github.com/TheMangoMen/backend/internal/service"
	"github.com/jackc/pgx/v5/pgconn"
)

type GetJobsBody struct {
	Jobs []model.Job `json:"jobs"`
}

func GetJobs(js service.JobService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := auth.FromContext(r.Context())
		uID := user.UID
		if !ok {
			// TODO: this mechanism needs to be better
			uID = ""
		}

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
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&jobs); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type UpdateWatchingBody struct {
	JIDs   []int `json:"jids"`
	Delete bool  `json:"delete"`
}

func UpdateWatching(js service.JobService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := UpdateWatchingBody{}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, ok := auth.FromContext(r.Context())
		if !ok {
			http.Error(w, "Error deleting watching", http.StatusBadRequest)
			return
		}

		updateFunc := js.CreateWatching
		if body.Delete {
			updateFunc = js.DeleteWatching
		}

		const foreignKeyViolationErrorCode = "23503"
		for _, jID := range body.JIDs {
			if err := updateFunc(user.UID, jID); err != nil {
				if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code != foreignKeyViolationErrorCode {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}
