package store

import (
	"database/sql"

	"github.com/TheMangoMen/backend/internal/model"
)

// func (s *Store) GetRanking(jID int) (model.Ranking, error) {
// 	var ranking model.Ranking
// 	query := "SELECT * FROM Rankings WHERE jid = $1;"
// 	err := s.db.Select(&rankings, query, jID)
// 	if err != nil {
// 		return nil, err
// 	}
// 	return ranking, nil
// }

func (s *Store) GetRanking(jID int, uID string) (model.Ranking, error) {
	var ranking model.Ranking
	query := `
SELECT * FROM Rankings
WHERE jid = $1 AND uid = $2;
`
	err := s.db.Get(&ranking, query, jID, uID)
	if err != nil {
		if err == sql.ErrNoRows {
			// we are guaranteed to look for a valid jID
			return model.Ranking{
				UID: uID,
				JID: jID,
			}, nil
		}
		return model.Ranking{}, err
	}

	return ranking, nil
}

func (s *Store) AddRanking(ranking model.Ranking) (err error) {
	query := `
INSERT INTO
    Rankings (uid, jid, userranking, employerranking)
VALUES
(:uid, :jid, :userranking, :employerranking) ON CONFLICT (uid, jid) DO
UPDATE
SET
    userranking = EXCLUDED.userranking,
    employerranking = EXCLUDED.employerranking;
`
	_, err = s.db.NamedExec(query, ranking)
	return err
}
