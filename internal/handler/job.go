package handler

import (
	"encoding/json"
	"net/http"

	"github.com/TheMangoMen/backend/internal/auth"
	"github.com/TheMangoMen/backend/internal/model"
	"github.com/TheMangoMen/backend/internal/service"
)

type GetJobsBody struct {
	Jobs []model.Job `json:"jobs"`
}

func GetJobs(js service.JobService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uID, ok := auth.FromContext(r.Context())
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

type CreateWatchingBody struct {
	UID  string `json:"uid"`
	JIDs []int  `json:"jids"`
}

type DeleteWatchingBody struct {
	JID    int  `json:"jid"`
	Delete bool `json:"delete"`
}

func UpdateWatching(js service.JobService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := DeleteWatchingBody{}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		uID, ok := auth.FromContext(r.Context())
		if !ok {
			http.Error(w, "Error deleting watching", http.StatusBadRequest)
			return
		}

		if body.Delete {
			if err := js.DeleteWatching(uID, body.JID); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		} else {
			jids := []int{body.JID}
			if err := js.CreateWatching(uID, jids); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusOK)
		}
	}
}

func CreateWatching(js service.JobService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := CreateWatchingBody{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := js.CreateWatching(body.UID, body.JIDs); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
