package handler

import (
	"encoding/json"
	"net/http"
	"github.com/TheMangoMen/backend/internal/service"
	"github.com/TheMangoMen/backend/internal/model"
)


type GetJobsBody struct {
	Jobs	[]model.Job		`json:"jobs"`
}

func GetJobs(js service.JobService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jobs, err := js.GetJobs(r.PathValue("uID"))
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