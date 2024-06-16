package store

import "github.com/TheMangoMen/backend/internal/model"

func (s *Store) GetRankings(jID int) (ranking []model.Ranking, err error) {
	var rankings []model.Ranking
	query := "SELECT * FROM Rankings WHERE JID = $1;"
	err = s.db.Select(&rankings, query, jID)
	if err != nil {
		return nil, err
	}
	return rankings, nil
}

func (s *Store) AddRanking(ranking model.Ranking) (err error) {
	query := `
INSERT INTO
    Rankings (UID, JID, UserRanking, EmployerRanking)
VALUES
(:UID, :JID, :UserRanking, :EmployerRanking) ON CONFLICT (UID, JID) DO
UPDATE
SET
    UserRanking = EXCLUDED.UserRanking,
    EmployerRanking = EXCLUDED.EmployerRanking;
`
	_, err = s.db.NamedExec(query, ranking)
	return err
}
