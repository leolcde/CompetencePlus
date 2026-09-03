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

type Video struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	ProfileID uint      `json:"-"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
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

type MeResponse struct {
	ID                  int64       `json:"id"`
	Email               string      `json:"email"`
	Role                string      `json:"role"`
	DateNaissance       string      `json:"date_naissance"`
	Identite            string      `json:"identite,omitempty"`
	Competences         []string    `json:"competences,omitempty"`
	Secteur             string      `json:"secteur,omitempty"`
	Localisation        string      `json:"localisation,omitempty"`
	StatutCertification string      `json:"statut_certification,omitempty"`
	CreatedAt           string      `json:"created_at"`
	Permissions         Permissions `json:"permissions"`
}

// PERMISSIONS

type Permissions struct {
	CanPublishVideo bool `json:"can_publish_video"`
	CanContact      bool `json:"can_contact"`
	CanLike         bool `json:"can_like"`
}

func PermissionsFor(role Role) Permissions {
	switch role {
	case Recruiter:
		return Permissions{CanContact: true, CanLike: true}
	case Candidate:
		return Permissions{CanPublishVideo: true}
	default:
		return Permissions{}
	}
}
