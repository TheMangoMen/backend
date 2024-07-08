package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/TheMangoMen/backend/internal/auth"
	"github.com/TheMangoMen/backend/internal/model"
	"github.com/TheMangoMen/backend/internal/service"
)

func GetContribution(cs service.ContributionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := auth.MustFromContext(r.Context())
		jID, err := strconv.Atoi(r.PathValue("jID"))
		if err != nil {
			http.Error(w, "invalid job id", http.StatusBadRequest)
		}
		contribution, err := cs.GetContribution(jID, user.UID)
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

type AddContributionBody struct {
	JID            int  `json:"jid" db:"jid"`
	OA             bool `json:"oa" db:"oa"`
	InterviewStage int  `json:"interviewstage" db:"interviewstage"`
	OfferCall      bool `json:"offercall" db:"offercall"`
}

func AddContribution(cs service.ContributionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := auth.MustFromContext(r.Context())
		body := AddContributionBody{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		contribution := model.Contribution{
			UID:            user.UID,
			JID:            body.JID,
			OA:             body.OA,
			InterviewStage: body.InterviewStage,
			OfferCall:      body.OfferCall,
		}
		if err := cs.AddContribution(contribution); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
