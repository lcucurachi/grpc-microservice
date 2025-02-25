package repository

import (
	"github.com/lokker96/grpc_project/domain/entity"
)

type ExplorerRepository interface {
	CreateUser(user *entity.User) error
	CreateDecision(decision *entity.Decision) error
	GetDecisionsForRecipientId(userID int, liked *bool) ([]entity.Decision, error)
	GetDecisionsForUserId(userID int, liked *bool) ([]entity.Decision, error)
	GetLikesCountByProfileId(profileID int) int64
	UpdateDecision(userID int, recipientUserId int, liked bool) error
	FindMutualLike(userID int, recipientUserID int) bool
}
