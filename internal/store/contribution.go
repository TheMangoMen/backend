package store

import "github.com/TheMangoMen/backend/internal/model"

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
