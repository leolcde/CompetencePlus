package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	utils "jibjob/src/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Role = utils.Role

const (
	Candidate = utils.Candidate
	Recruiter = utils.Recruiter
	Admin     = utils.Admin
)

type RegisterRequest struct {
	Name        string   `json:"name"`
	Mail        string   `json:"mail"`
	Password    string   `json:"password"`
	Role        Role     `json:"role"`
	DateOfBirth string   `json:"date_of_birth"` // formatté comme ça : 2006-01-02
	Skills      []string `json:"skills"`
	Sector      string   `json:"sector"`
	Location    string   `json:"location"`
}

func checkRole(r Role) bool {
	switch r {
	case "", Candidate, Recruiter, Admin:
		return true
	}
	return false
}

func RegisterErrorHandling(res http.ResponseWriter, req *http.Request) (RegisterRequest, time.Time, bool) {
	var r RegisterRequest
	if err := json.NewDecoder(req.Body).Decode(&r); err != nil {
		http.Error(res, `{"error":"invalid request"}`, http.StatusBadRequest)
		return r, time.Time{}, false
	}
	// erreur valeur manquante
	if r.Name == "" || r.Mail == "" || r.Password == "" {
		http.Error(res, `{"error":"missing value"}`, http.StatusUnprocessableEntity)
		return r, time.Time{}, false
	}
	// l'anniv c'est pas correct
	out, err := time.Parse("2006-01-02", r.DateOfBirth)
	if err != nil {
		http.Error(res, `{"error":"invalid birthday"}`, http.StatusUnprocessableEntity)
		return r, time.Time{}, false
	}
	// évite les erreur bcrypt 500, > a 72octets
	if len(r.Password) > 72 {
		http.Error(res, `{"error":"password too long"}`, http.StatusUnprocessableEntity)
		return r, time.Time{}, false
	}
	// format les emails
	if !strings.Contains(r.Mail, "@") {
		http.Error(res, `{"error":"invalid mail"}`, http.StatusUnprocessableEntity)
		return r, time.Time{}, false
	}
	// check si le role est bien, candidate ou recruite ou admin
	if !checkRole(r.Role) {
		http.Error(res, `{"error":"invalid role"}`, http.StatusUnprocessableEntity)
		return r, time.Time{}, false
	}
	return r, out, true
}

func Register(gdb *gorm.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		r, dob, ok := RegisterErrorHandling(res, req)
		if !ok {
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(res, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}

		profil := utils.Profile{
			Name:         r.Name,
			Email:        r.Mail,
			PasswordHash: string(hash),
			DateOfBirth:  dob,
			Skills:       r.Skills,
			Sector:       r.Sector,
			Location:     r.Location,
			Role:         r.Role,
		}

		if err := gdb.Create(&profil).Error; err != nil {
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				http.Error(res, `{"error":"mail already used"}`, http.StatusUnprocessableEntity)
				return
			}
			http.Error(res, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}

		res.Header().Set("Content-Type", "application/json")
		res.WriteHeader(http.StatusCreated)
		json.NewEncoder(res).Encode(map[string]any{
			"id":   profil.ID,
			"mail": profil.Email,
			"role": profil.Role,
		})
	}
}
