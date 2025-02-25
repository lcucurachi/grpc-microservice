package entity

import (
	"time"
)

type Decision struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	AuthorID    uint      // Author who made the decision
	RecipientID uint      // Profile that was presented to the author
	Liked       bool      // True if liked, false if not. This ideally would be an enum with types PASS and LIKE
	Author      User      // gorm uses the author_id to fill this structure with the relational data
	Recipient   User      // gorm uses the profile_id to fill this structure with the relational data
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`
}

func (Decision) TableName() string {
	return "decisions"
}
