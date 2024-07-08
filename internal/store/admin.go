package store

import "github.com/TheMangoMen/backend/internal/model"

func (s *Store) UpdateStage(isRankingStage bool) (bool, error) {
	var updatedRankingStage bool
	query := `
UPDATE Stage
SET isRankingStage = $1
RETURNING *;
`
	err := s.db.QueryRow(query, isRankingStage).Scan(&updatedRankingStage)
	if err != nil {
		return false, err
	}
	return updatedRankingStage, nil
}

func (s *Store) UpdateYear(year int) (int, error) {
	var updatedYear int
	query := `
UPDATE Year
SET year = $1
RETURNING *;
`
	err := s.db.QueryRow(query, year).Scan(&updatedYear)
	if err != nil {
		return 0, err
	}
	return updatedYear, err
}

func (s *Store) UpdateSeason(season string) (string, error) {
	var updatedYear string
	query := `
UPDATE Season
SET season = $1
RETURNING *;
`
	err := s.db.QueryRow(query, season).Scan(&updatedYear)
	if err != nil {
		return "", err
	}
	return updatedYear, nil
}

func (s *Store) UpdateCycle(cycle int) (int, error) {
	var updatedCycle int
	query := `
UPDATE Cycle
SET cycle = $1
RETURNING *;
`
	err := s.db.QueryRow(query, cycle).Scan(&updatedCycle)
	if err != nil {
		return 0, err
	}
	return updatedCycle, nil
}

func (s *Store) GetStage() (bool, error) {
	var isRankingStage bool
	query := `
SELECT *
FROM Stage;
`
	err := s.db.Get(&isRankingStage, query)
	if err != nil {
		return false, err
	}
	return isRankingStage, nil
}

func (s *Store) GetYear() (int, error) {
	var year int
	query := `
SELECT *
FROM Year;
`
	err := s.db.Get(&year, query)
	if err != nil {
		return 0, err
	}
	return year, nil

}

func (s *Store) GetSeason() (string, error) {
	var season string
	query := `
SELECT *
FROM Season;
`
	err := s.db.Get(&season, query)
	if err != nil {
		return "", err
	}
	return season, nil

}

func (s *Store) GetCycle() (int, error) {
	var cycle int
	query := `
SELECT *
FROM Cycle;
`
	err := s.db.Get(&cycle, query)
	if err != nil {
		return 0, err
	}
	return cycle, nil

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
