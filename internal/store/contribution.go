package store

import (
	"database/sql"

	"github.com/TheMangoMen/backend/internal/model"
)

func (s *Store) GetContribution(jID int, uID string) (model.Contribution, error) {
	var contribution model.Contribution
	query := `
SELECT *
FROM Contributions
WHERE jid = $1 AND uid = $2;
`
	err := s.db.Get(&contribution, query, jID, uID)
	if err != nil {
		if err == sql.ErrNoRows {
			// we are guaranteed to look for a valid jID
			return model.Contribution{
				UID: uID,
				JID: jID,
			}, nil
		}
		return model.Contribution{}, err
	}
	return contribution, nil
}
func (s *Store) AddContribution(contribution model.Contribution) (err error) {
	query := `
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
	_, err = s.db.NamedExec(query, contribution)
	return err
}
