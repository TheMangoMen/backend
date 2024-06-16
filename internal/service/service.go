package service

import "github.com/TheMangoMen/backend/internal/model"

type UserService interface {
	GetUser(uID string) (model.User, error)
	CreateUser(uID string) error
}

type ContributionService interface {
	GetContribution(uID string, jID string) (model.Contribution, error)
	CreateContribution(uID string, jID string, oa bool, interviewStage int, offerCall bool) error
}
