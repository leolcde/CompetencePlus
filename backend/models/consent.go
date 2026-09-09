package models

import "time"

type Consent struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	UserID     uint       `json:"user_id" gorm:"not null;index"`
	Version    string     `json:"version" gorm:"not null"`
	ApprouveAt time.Time  `json:"granted_at" gorm:"not null"`
	RevokedAt  *time.Time `json:"revoked_at,omitempty"`

	User User `json:"-" gorm:"constraint:OnDelete:CASCADE;"`
}
