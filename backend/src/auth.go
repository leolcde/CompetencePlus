package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"profilsactifs/models"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const userIDKey ctxKey = "userID"

func currentUserID(req *http.Request) uint {
	return req.Context().Value(userIDKey).(uint)
}

func parseToken(req *http.Request) (uint, string, error) {
	tokenStr := strings.TrimPrefix(req.Header.Get("Authorization"), "Bearer ")
	if tokenStr == "" {
		return 0, "", fmt.Errorf("missing token")
	}
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil || !token.Valid {
		return 0, "", fmt.Errorf("invalid token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, "", fmt.Errorf("invalid claims")
	}
	return uint(claims["sub"].(float64)), fmt.Sprint(claims["role"]), nil
}

func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		userID, _, err := parseToken(req)
		if err != nil {
			http.Error(res, "unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(req.Context(), userIDKey, userID)
		next(res, req.WithContext(ctx))
	}
}

func requireRole(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		userID, r, err := parseToken(req)
		if err != nil {
			http.Error(res, "unauthorized", http.StatusUnauthorized)
			return
		}
		if r != role {
			http.Error(res, "forbidden", http.StatusForbidden)
			return
		}
		ctx := context.WithValue(req.Context(), userIDKey, userID)
		next(res, req.WithContext(ctx))
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func register(res http.ResponseWriter, req *http.Request) {
	var body struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
		Birthday string `json:"birthday"`
	}

	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(res, "invalid body", http.StatusBadRequest)
		return
	}

	if body.Name == "" || body.Email == "" || len(body.Password) < 6 {
		http.Error(res, "name, email and password (6+ chars) are required", http.StatusBadRequest)
		return
	}

	birth, err := time.Parse("2006-01-02", body.Birthday)
	if err != nil {
		http.Error(res, "invalid birthday (expected YYYY-MM-DD)", http.StatusBadRequest)
		return
	}

	years := computeAge(birth)
	if years < 16 {
		http.Error(res, "you must be at least 16 years old to register", http.StatusForbidden)
		return
	}

	status := string(models.StatusAdult)
	if years < 18 {
		status = string(models.StatusYouth)
	}

	hash, err := hashPassword(body.Password)
	if err != nil {
		http.Error(res, "server error", http.StatusInternalServerError)
		return
	}

	user := models.User{
		Name:         body.Name,
		Email:        strings.ToLower(body.Email),
		PasswordHash: hash,
		BirthDay:     birth,
		Status:       status,
		Role:         "candidate",
	}

	if err := DB.Create(&user).Error; err != nil {
		http.Error(res, "email already in use", http.StatusConflict)
		return
	}

	token, _ := generateToken(user.ID, user.Role)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(map[string]any{
		"token": token,
		"user":  map[string]any{"id": user.ID, "name": user.Name, "email": user.Email, "status": user.Status, "role": user.Role},
	})
}

func login(res http.ResponseWriter, req *http.Request) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(res, "invalid body", http.StatusBadRequest)
		return
	}

	var user models.User
	if err := DB.Where("email = ?", strings.ToLower(body.Email)).First(&user).Error; err != nil {
		http.Error(res, "invalid credentials", http.StatusUnauthorized) // user not found
		return
	}

	if !checkPassword(user.PasswordHash, body.Password) {
		http.Error(res, "invalid credentials", http.StatusUnauthorized) // wrong password
		return
	}

	token, _ := generateToken(user.ID, user.Role)

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(map[string]any{
		"token": token,
		"user":  map[string]any{"id": user.ID, "name": user.Name, "email": user.Email, "role": user.Role},
	})
}
