package models

import (
	"time"

	"github.com/lib/pq"
)

type Role string
type Status string

const (
	Candidate Role = "candidate"
	Recruiter Role = "recruiter"
	Admin     Role = "admin"
)

const (
	StatusYouth Status = "youth"
	StatusAdult Status = "adult"
)

type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Name         string    `json:"name" gorm:"not null"`
	Email        string    `json:"email" gorm:"uniqueIndex;not null"`
	PasswordHash string    `json:"-" gorm:"not null"`
	BirthDay     time.Time `json:"birthday" gorm:"not null"`
	Status       string    `json:"status" gorm:"not null;default:'adult'"`

	Skills pq.StringArray `json:"skills" gorm:"type:text[]"`
	Sector string         `json:"sector"`
	City   string         `json:"city"`
	Role   string         `json:"role" gorm:"not null;default:'candidate'"`

	CreatedAt time.Time `json:"created_at"`
}
