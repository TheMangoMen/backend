package service

import "github.com/TheMangoMen/backend/internal/model"

type UserService interface {
	GetUser(uID string) (model.User, error)
	CreateUser(uID string) error
}

type ContributionService interface {
	GetContribution(jID int, uID string) (model.Contribution, model.ContributionTags, error)
	AddContribution(contribution model.Contribution, tags model.ContributionTags) (err error)
}

type JobService interface {
	CreateWatching(uID string, jIDs []int) error
	DeleteWatching(uID string, jID int) error
	GetJobInterviews(uID string) ([]model.Job, error)
	GetJobRankings(uID string) ([]model.Job, error)
	GetIsRankingStage() (isRankingStage bool, err error)
}

type RankingService interface {
	GetRankings(jID int) ([]model.Ranking, error)
	AddRanking(ranking model.Ranking) error
}
