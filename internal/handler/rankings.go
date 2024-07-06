package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/TheMangoMen/backend/internal/auth"
	"github.com/TheMangoMen/backend/internal/model"
	"github.com/TheMangoMen/backend/internal/service"
)

func GetRankings(rs service.RankingService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		jID, err := strconv.Atoi(r.PathValue("jID"))
		if err != nil {
			http.Error(w, "invalid job id", http.StatusBadRequest)
		}
		ranking, err := rs.GetRankings(jID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				http.Error(w, "404 job not found", http.StatusNotFound)
			} else {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&ranking); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

type AddRankingBody struct {
	JID             int    `json:"jid"`
	UserRanking     int    `json:"userranking"`
	EmployerRanking string `json:"employerranking"`
}

func AddRanking(rs service.RankingService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uID := auth.MustFromContext(r.Context())
		body := AddRankingBody{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		ranking := model.Ranking{
			UID:             uID,
			JID:             body.JID,
			UserRanking:     body.UserRanking,
			EmployerRanking: body.EmployerRanking,
		}
		if err := rs.AddRanking(ranking); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
