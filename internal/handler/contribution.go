package handler

import (
	"encoding/json"
	"net/http"

	"github.com/TheMangoMen/backend/internal/model"
	"github.com/TheMangoMen/backend/internal/service"
)

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
