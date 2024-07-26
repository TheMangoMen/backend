package service

import "github.com/TheMangoMen/backend/internal/model"

type UserService interface {
	GetUser(uID string) (model.User, error)
	GetIsAdmin(uID string) (bool, error)
	CreateUser(uID string) error
}

type ContributionService interface {
	GetContribution(jID int, uID string) (model.ContributionCombined, error)
	AddContribution(contribution model.Contribution, tags model.ContributionTags) (err error)
}

type JobService interface {
	CreateWatching(uID string, jID int) error
	DeleteWatching(uID string, jID int) error
	GetJobInterviews(uID string) ([]model.Job, error)
	GetJobRankings(uID string) ([]model.Job, error)
	GetIsRankingStage() (isRankingStage bool, err error)
}

type RankingService interface {
	GetRanking(jID int, uID string) (model.Ranking, error)
	AddRanking(ranking model.Ranking) error
}

type AnalyticsService interface {
	GetWatchedJobsStatusCounts(uID string) ([]model.StatusCount, error)
	GetWatchedCompaniesStatusCounts(uID string) ([]model.StatusCount, error)
}

type AdminService interface {
	UpdateStage(isRankingStage bool) (bool, error)
	UpdateYear(year int) (int, error)
	UpdateSeason(season string) (string, error)
	UpdateCycle(cycle int) (int, error)
	GetStage() (bool, error)
	GetYear() (int, error)
	GetSeason() (string, error)
	GetCycle() (int, error)
	GetContributionLogs() ([]model.ContributionLog, error)
}
