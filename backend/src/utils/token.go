package utils

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func ParseProfileToken(req *http.Request) (*Claims, bool) {
	header := req.Header.Get("Authorization")
	raw := strings.TrimSpace(strings.TrimPrefix(header, "Bearer "))
	if raw == "" || raw == header {
		return nil, false
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, false
	}

	var claims Claims
	token, err := jwt.ParseWithClaims(raw, &claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(secret), nil
	})
	if err != nil || !token.Valid {
		return nil, false
	}
	return &claims, true
}

func RequireRole(role Role, next http.HandlerFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		claims, ok := ParseProfileToken(req)
		if !ok {
			res.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(res).Encode(map[string]string{"error": "unauthorized"})
			return
		}
		if Role(claims.Role) != role {
			res.WriteHeader(http.StatusForbidden)
			json.NewEncoder(res).Encode(map[string]string{"error": "forbidden"})
			return
		}
		next(res, req)
	}
}
