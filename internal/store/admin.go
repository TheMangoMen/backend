package store

import "github.com/TheMangoMen/backend/internal/model"

func (s *Store) UpdateStage(isRankingStage bool) error {
	query := `
UPDATE Stage
SET isRankingStage = $1;
`
	_, err := s.db.Exec(query, isRankingStage)
	return err
}

func (s *Store) UpdateYear(year int) error {
	query := `
UPDATE Year
SET year = $1;
`
	_, err := s.db.Exec(query, year)
	return err
}

func (s *Store) UpdateSeason(season string) error {
	query := `
UPDATE Season
SET season = $1;
`
	_, err := s.db.Exec(query, season)
	return err
}

func (s *Store) UpdateCycle(cycle int) error {
	query := `
UPDATE Cycle
SET cycle = $1;
`
	_, err := s.db.Exec(query, cycle)
	return err
}

// TODO: add pagination
func (s *Store) GetContributionLogs() ([]model.ContributionLog, error) {
	var contributionLogs []model.ContributionLog
	query := `
SELECT *
FROM ContributionsLogs;
`
	err := s.db.Select(&contributionLogs, query)
	if err != nil {
		return []model.ContributionLog{}, err
	}
	return contributionLogs, nil
}
