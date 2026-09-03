package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	utils "jibjob/src/utils"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func jwtSecret() []byte {
	s := os.Getenv("JWT_SECRET")
	if s == "" {
		panic("JWT_SECRET manquant")
	}
	return []byte(s)
}

func LoginErrorHandling(gdb *gorm.DB, res http.ResponseWriter, req *http.Request) (utils.Profile, bool) {
	var body utils.LoginRequest
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(res, `{"error":"invalid request"}`, http.StatusBadRequest)
		return utils.Profile{}, false
	}
	if body.Mail == "" || body.Password == "" {
		http.Error(res, `{"error":"missing value"}`, http.StatusUnprocessableEntity)
		return utils.Profile{}, false
	}

	var profil utils.Profile
	if err := gdb.Where("email = ?", body.Mail).First(&profil).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(res, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
			return utils.Profile{}, false
		}
		http.Error(res, `{"error":"internal error"}`, http.StatusInternalServerError)
		return utils.Profile{}, false
	}

	if bcrypt.CompareHashAndPassword([]byte(profil.PasswordHash), []byte(body.Password)) != nil {
		http.Error(res, `{"error":"invalid credentials"}`, http.StatusUnauthorized)
		return utils.Profile{}, false
	}

	return profil, true
}

func Login(gdb *gorm.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		profil, ok := LoginErrorHandling(gdb, res, req)
		if !ok {
			return
		}

		claims := utils.Claims{
			UserID: fmt.Sprint(profil.ID),
			Role:   string(profil.Role),
			RegisteredClaims: jwt.RegisteredClaims{
				IssuedAt:  jwt.NewNumericDate(time.Now()),
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			},
		}
		token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(jwtSecret())
		if err != nil {
			http.Error(res, `{"error":"internal error"}`, http.StatusInternalServerError)
			return
		}

		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(map[string]any{
			"token": token,
			"id":    profil.ID,
			"name":  profil.Name,
			"mail":  profil.Email,
			"role":  profil.Role,
		})
	}
}
