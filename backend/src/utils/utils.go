package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/lib/pq"
)

// REGISTER

type Role string

const (
	Candidate Role = "candidate"
	Recruiter Role = "recruiter"
	Admin     Role = "admin"
)

type Profile struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name"`
	Email        string         `json:"email" gorm:"uniqueIndex"`
	PasswordHash string         `json:"-"`
	DateOfBirth  time.Time      `json:"date_of_birth"`
	Skills       pq.StringArray `json:"skills" gorm:"type:text[]"`
	Sector       string         `json:"sector"`
	Location     string         `json:"location"`
	Role         Role           `json:"role" gorm:"type:user_role;default:candidate"`
	CreatedAt    time.Time      `json:"created_at"`
}

// LOGIN

type LoginRequest struct {
	Mail     string `json:"mail"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Role  string `json:"role"`
}

type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}
