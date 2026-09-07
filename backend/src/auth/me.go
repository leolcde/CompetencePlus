package auth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	utils "jibjob/src/utils"

	"gorm.io/gorm"
)

func meProfile(gdb *gorm.DB, res http.ResponseWriter, req *http.Request) (utils.Profile, bool) {
	claims, ok := utils.ParseProfileToken(req)
	if !ok {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode(map[string]string{"error": "unauthorized"})
		return utils.Profile{}, false
	}

	userID, err := strconv.ParseUint(claims.UserID, 10, 64)
	if err != nil {
		res.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(res).Encode(map[string]string{"error": "invalid token"})
		return utils.Profile{}, false
	}

	var profile utils.Profile
	err = gdb.First(&profile, uint(userID)).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		res.WriteHeader(http.StatusNotFound)
		json.NewEncoder(res).Encode(map[string]string{"error": "user not found"})
		return utils.Profile{}, false
	}
	if err != nil {
		res.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(res).Encode(map[string]string{"error": "internal error"})
		return utils.Profile{}, false
	}
	return profile, true
}

func Me(gdb *gorm.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		profile, ok := meProfile(gdb, res, req)
		if !ok {
			return
		}

		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(utils.MeResponse{
			ID:            int64(profile.ID),
			Email:         profile.Email,
			Role:          string(profile.Role),
			DateNaissance: profile.DateOfBirth.Format("2006-01-02"),
			Identite:      profile.Name,
			Competences:   []string(profile.Skills),
			Secteur:       profile.Sector,
			Localisation:  profile.Location,
			CreatedAt:     profile.CreatedAt.Format(time.RFC3339),
			Permissions:   utils.PermissionsFor(profile.Role),
		})
	}
}
