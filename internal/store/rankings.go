package store

import "github.com/TheMangoMen/backend/internal/model"

func (s *Store) GetRankings(jID int) ([]model.Ranking, error) {
	var rankings []model.Ranking
	query := "SELECT * FROM Rankings WHERE jid = $1;"
	err := s.db.Select(&rankings, query, jID)
	if err != nil {
		return nil, err
	}
	return rankings, nil
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
