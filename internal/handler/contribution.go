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
			return
		}
		contributionCombined, err := cs.GetContribution(jID, user.UID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&contributionCombined); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}

type AddContributionBody struct {
	JID            int     `json:"jid" db:"jid"`
	OA             bool    `json:"oa" db:"oa"`
	Interview      bool    `json:"interview"`
	InterviewStage int     `json:"interviewcount" db:"interviewstage"`
	OfferCall      bool    `json:"offercall" db:"offercall"`
	OADifficulty   string  `json:"oadifficulty" db:"oadifficulty"`
	OALength       string  `json:"oalength" db:"oalength"`
	InterviewVibe  string  `json:"interviewvibe" db:"interviewvibe"`
	InterviewTech  string  `json:"interviewtechnical" db:"interviewtechnical"`
	OfferComp      float32 `json:"compensation" db:"offercomp"`
}

func AddContribution(cs service.ContributionService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := auth.MustFromContext(r.Context())
		body := AddContributionBody{}
		interviewCnt := 0
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if body.Interview {
			interviewCnt = body.InterviewStage
		}
		contribution := model.Contribution{
			UID:            user.UID,
			JID:            body.JID,
			OA:             body.OA,
			InterviewStage: interviewCnt,
			OfferCall:      body.OfferCall,
		}

		contributionTags := model.ContributionTags{
			UID:           user.UID,
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
