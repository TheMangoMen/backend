package store

import (
	"database/sql"

	"github.com/TheMangoMen/backend/internal/model"
)

func (s *Store) GetContribution(jID int, uID string) (model.Contribution, model.ContributionTags, error) {
	var contribution model.Contribution
	var contributionTags model.ContributionTags
	query1 := `
SELECT * FROM Contributions NATURAL JOIN Tags
WHERE jid = $1 AND uid = $2;
`
	err1 := s.db.Get(&contribution, query1, jID, uID)
	if err1 != nil {
		if err1 == sql.ErrNoRows {
			// we are guaranteed to look for a valid jID
			return model.Contribution{
					UID: uID,
					JID: jID,
				},
				model.ContributionTags{}, nil
		}
		return model.Contribution{}, model.ContributionTags{}, err1
	}

	query2 := `
SELECT * FROM Tags
WHERE jid = $1 AND uid = $2;
`
	err2 := s.db.Get(&contributionTags, query2, jID, uID)
	if err2 != nil {
		if err2 == sql.ErrNoRows {
			// we are guaranteed to look for a valid jID
			return contribution,
				model.ContributionTags{
					UID: uID,
					JID: jID,
				}, nil
		}
		return contribution, model.ContributionTags{}, err2
	}

	return contribution, contributionTags, nil
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
