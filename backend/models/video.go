package models

import "time"

type Video struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	User      User      `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
	URL       string    `json:"url" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
}
