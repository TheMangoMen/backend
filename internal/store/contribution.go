package store

import (
	"database/sql"

	"github.com/TheMangoMen/backend/internal/model"
)

func (s *Store) GetContribution(jID int, uID string) (model.ContributionCombined, error) {
	var contributionCombined model.ContributionCombined
	query1 := `
SELECT * FROM Contributions NATURAL JOIN Tags
WHERE jid = $1 AND uid = $2;
`
	err1 := s.db.Get(&contributionCombined, query1, jID, uID)
	if err1 != nil {
		if err1 == sql.ErrNoRows {
			// we are guaranteed to look for a valid jID
			return model.ContributionCombined{
				UID: uID,
				JID: jID,
			}, nil
		}
		return model.ContributionCombined{}, err1
	}

	return contributionCombined, nil
}
func (s *Store) AddContribution(contribution model.Contribution, tags model.ContributionTags) (err error) {
	query1 := `
INSERT INTO
    Contributions (uid, jid, oa, interviewstage, offercall)
VALUES
    (:uid, :jid, :oa, :interviewstage, :offercall) ON CONFLICT (uid, jid) DO
UPDATE
SET
    oa = EXCLUDED.oa,
    interviewstage = EXCLUDED.interviewstage,
    offercall = EXCLUDED.offercall;
`
	_, err = s.db.NamedExec(query1, contribution)
	if err != nil {
		return err
	}

	if !contribution.OA && contribution.InterviewStage == 0 && !contribution.OfferCall {
		return err
	}

	query2 := `
INSERT INTO
    Tags (uid, jid, oadifficulty, oalength, interviewvibe, interviewtechnical, offercomp)
VALUES
    (:uid, :jid, :oadifficulty, :oalength, :interviewvibe, :interviewtechnical, :offercomp) ON CONFLICT (uid, jid) DO
UPDATE
SET
    oadifficulty = EXCLUDED.oadifficulty,
    oalength = EXCLUDED.oalength,
    interviewvibe = EXCLUDED.interviewvibe,
    interviewtechnical = EXCLUDED.interviewtechnical,
    offercomp = EXCLUDED.offercomp;
`
	_, err = s.db.NamedExec(query2, tags)
	return err
}
