package repository

import (
	"context"

	"github.com/lokker96/grpc_project/domain/entity"
)

type ExplorerRepository interface {
	CreateUser(ctx context.Context, user *entity.User) error
	CreateDecision(ctx context.Context, decision *entity.Decision) error
	GetDecisionsForRecipientId(ctx context.Context, userID int, liked *bool) ([]entity.Decision, error)
	GetDecisionsForUserId(ctx context.Context, userID int, liked *bool) ([]entity.Decision, error)
	GetLikesCountByProfileId(ctx context.Context, profileID int) int64
	UpdateDecision(ctx context.Context, userID int, recipientUserId int, liked bool) error
	FindMutualLike(ctx context.Context, userID int, recipientUserID int) bool
}
