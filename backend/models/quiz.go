package models

import (
	"time"

	"github.com/lib/pq"
)

type Question struct {
	ID      uint           `json:"id" gorm:"primaryKey"`
	Content string         `json:"content" gorm:"size:255;not null"`
	Options pq.StringArray `json:"options" gorm:"type:text[];not null;default:'{}'"`
	Weight  int            `json:"weight" gorm:"not null;default:1"`
}

type Answer struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_user_question"`
	QuestionID uint      `json:"question_id" gorm:"not null;uniqueIndex:idx_user_question"`
	Choice     string    `json:"choice" gorm:"size:120;not null"`
	CreatedAt  time.Time `json:"created_at"`

	User     User     `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
	Question Question `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
}

type BadgeResult struct {
	UserID    uint      `json:"user_id" gorm:"primaryKey"`
	Score     int       `json:"score" gorm:"not null;default:0"`
	Badge     bool      `json:"badge" gorm:"not null;default:false"`
	UpdatedAt time.Time `json:"updated_at"`

	User User `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
}
