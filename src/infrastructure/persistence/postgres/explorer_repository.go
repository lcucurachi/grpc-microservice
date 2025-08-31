package postgres

import (
	"context"
	"fmt"

	"github.com/lokker96/grpc_project/domain/entity"
	domainError "github.com/lokker96/grpc_project/domain/error"
	"github.com/lokker96/grpc_project/domain/repository"

	"errors"

	"gorm.io/gorm"
)

// The explorer repository implements the method we can use to access data from the DB.
type explorerRepository struct {
	db *gorm.DB
}

func NewExplorerRepository(db *gorm.DB) repository.ExplorerRepository {
	return &explorerRepository{
		db: db,
	}
}

// Helper function to setup some dummy data
func (r *explorerRepository) CreateUser(ctx context.Context, user *entity.User) error {
	// using gorm transactions to make it easier to rollback if there are any issues,
	// this is typically used in more complex repository methods to preserve data integrity
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(&user).Error; err != nil {
			return fmt.Errorf("error on creating user in db: %w", err)
		}

		return nil
	})
}

func (r *explorerRepository) CreateDecision(ctx context.Context, decision *entity.Decision) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.WithContext(ctx).Create(&decision).Error; err != nil {
			return fmt.Errorf("error on creating decision in db: %w", err)
		}

		return nil
	})
}

func (r *explorerRepository) GetDecisionsForRecipientId(ctx context.Context, userID int, liked *bool) ([]entity.Decision, error) {
	var result []entity.Decision

	queryBuilder := r.db.WithContext(ctx).Model(&entity.Decision{})

	queryBuilder = queryBuilder.Where("recipient_id = ?", uint(userID))

	if liked != nil {
		queryBuilder = queryBuilder.Where("liked = ?", *liked)
	}

	err := queryBuilder.Find(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainError.NewDecisionNotFoundErr()
		} else {
			return nil, fmt.Errorf("error searching for decisions liked for user id: %w", err)
		}
	}

	return result, nil
}

func (r *explorerRepository) GetDecisionsForUserId(ctx context.Context, userID int, liked *bool) ([]entity.Decision, error) {
	var result []entity.Decision

	queryBuilder := r.db.WithContext(ctx).Model(&entity.Decision{})

	queryBuilder = queryBuilder.Where("author_id = ?", uint(userID))

	if liked != nil {
		queryBuilder = queryBuilder.Where("liked = ?", *liked)
	}

	err := queryBuilder.Find(&result).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domainError.NewDecisionNotFoundErr()
		} else {
			return nil, fmt.Errorf("error searching for decisions for user id: %w", err)
		}
	}

	return result, nil
}

func (r *explorerRepository) GetLikesCountByProfileId(ctx context.Context, profileID int) int64 {
	var count int64

	queryBuilder := r.db.WithContext(ctx).Model(&entity.Decision{})

	queryBuilder = queryBuilder.Where("recipient_id = ?", uint(profileID))
	queryBuilder = queryBuilder.Where("liked = ?", true)

	queryBuilder.Count(&count)

	return count
}

func (r *explorerRepository) UpdateDecision(ctx context.Context, userID int, recipientUserId int, liked bool) error {
	queryBuilder := r.db.WithContext(ctx).Model(&entity.Decision{})

	queryBuilder = queryBuilder.Where("author_id = ?", uint(userID))
	queryBuilder = queryBuilder.Where("recipient_id = ?", uint(recipientUserId))

	result := queryBuilder.Update("liked", liked)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return domainError.NewDecisionNotFoundErr()
		} else {
			return fmt.Errorf("error updating decision: %w", result.Error)
		}
	}

	return nil
}

func (r *explorerRepository) FindMutualLike(ctx context.Context, userID int, recipientUserID int) bool {
	var actorLikesCount int64
	var recipientLikesCount int64

	queryBuilder := r.db.WithContext(ctx).Model(&entity.Decision{})

	queryBuilder = queryBuilder.Where("author_id = ?", userID)
	queryBuilder = queryBuilder.Where("recipient_id = ?", recipientUserID)
	queryBuilder = queryBuilder.Where("liked = ?", true)

	queryBuilder.Count(&actorLikesCount)

	queryBuilder = r.db.WithContext(ctx).Model(&entity.Decision{})

	queryBuilder = queryBuilder.Where("author_id = ?", recipientUserID)
	queryBuilder = queryBuilder.Where("recipient_id = ?", userID)
	queryBuilder = queryBuilder.Where("liked = ?", true)

	queryBuilder.Count(&recipientLikesCount)

	// We are assuming that we will only store 1 decision between actors and recipients
	return actorLikesCount == 1 && recipientLikesCount == 1
}
