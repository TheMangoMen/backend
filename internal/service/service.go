package service

import "github.com/TheMangoMen/backend/internal/model"

type UserService interface {
	GetUser(uID string) (model.User, error)
	CreateUser(uID string) error
}

type ContributionService interface {
	GetContribution(jID int, uID string) (model.Contribution, error)
	AddContribution(contribution model.Contribution) (err error)
}

type JobService interface {

	GetJobs(uID string) ([]model.Job, error)
	CreateWatching(uID string, jIDs []string) error
	GetJobInterviews(uID string) ([]model.Job, error)
	GetJobRankings(uID string) ([]model.Job, error)
	GetIsRankingStage() (isRankingStage bool, err error)

}

type RankingService interface {
	GetRankings(jID int) ([]model.Ranking, error)
	AddRanking(ranking model.Ranking) error
}
