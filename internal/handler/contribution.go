package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/TheMangoMen/backend/internal/service"
)

type CreateContributionBody struct {
	UID string
}

func CreateContribution(us service.UserService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := CreateContributionBody{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := us.CreateUser(body.UID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

type UpdateContributionBody struct {
	UID            string `json:"uid"`
	JID            string `json:"jid"`
	OA             bool   `json:"oa"`
	InterviewStage int    `json:"interview_stage"`
	OfferCall      bool   `json:"offer_call"`
}

func UpdateContribution(cs service.ContributionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body := UpdateContributionBody{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := cs.CreateContribution(body.UID, body.JID, body.OA, body.InterviewStage, body.OfferCall); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
