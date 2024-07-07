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
		uID := auth.MustFromContext(r.Context())
		jID, err := strconv.Atoi(r.PathValue("jID"))
		if err != nil {
			http.Error(w, "invalid job id", http.StatusBadRequest)
		}
		contribution, contributionTags, err := cs.GetContribution(jID, uID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(&contribution); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if err := json.NewEncoder(w).Encode(&contributionTags); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}

type AddContributionBody struct {
	JID            int    `json:"jid" db:"jid"`
	OA             bool   `json:"oa" db:"oa"`
	InterviewStage int    `json:"interviewcount" db:"interviewstage"`
	OfferCall      bool   `json:"offercall" db:"offercall"`
	OADifficulty   string `json:"oadifficulty" db:"oa1"`
	OALength       string `json:"oalength" db:"oa2"`
	InterviewVibe  string `json:"interviewvibe" db:"int1"`
	InterviewTech  string `json:"interviewtechnical" db:"int2"`
	OfferComp      int    `json:"compensation" db:"offer1"`
}

func AddContribution(cs service.ContributionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uID := auth.MustFromContext(r.Context())
		body := AddContributionBody{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		contribution := model.Contribution{
			UID:            uID,
			JID:            body.JID,
			OA:             body.OA,
			InterviewStage: body.InterviewStage,
			OfferCall:      body.OfferCall,
		}

		contributionTags := model.ContributionTags{
			UID:           uID,
			JID:           body.JID,
			OADifficulty:  body.OADifficulty,
			OALength:      body.OALength,
			InterviewVibe: body.InterviewVibe,
			InterviewTech: body.InterviewTech,
			OfferComp:     body.OfferComp,
		}
		//mapping to create a contributionTags
		if err := cs.AddContribution(contribution, contributionTags); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}
}
