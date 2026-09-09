package main

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"profilsactifs/models"
)

const consentVersion = "1.0"

func deleteUserVideos(userID uint) {
	var videos []models.Video
	DB.Where("user_id = ?", userID).Find(&videos)

	DB.Where("user_id = ?", userID).Delete(&models.Video{})

	for _, v := range videos {
		if strings.HasPrefix(v.URL, "/uploads/") {
			os.Remove(filepath.Join("uploads", filepath.Base(v.URL)))
		}
	}
}

func getConsent(res http.ResponseWriter, req *http.Request) {
	userID := currentUserID(req)

	var c models.Consent
	err := DB.Where("user_id = ? AND version = ? AND revoked_at IS NULL",
		userID, consentVersion).First(&c).Error

	res.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(res).Encode(map[string]any{"active": false, "version": consentVersion})
		return
	}
	json.NewEncoder(res).Encode(map[string]any{"active": true, "consent": c})
}

func approuveConsent(res http.ResponseWriter, req *http.Request) {
	userID := currentUserID(req)

	c := models.Consent{
		UserID:     userID,
		Version:    consentVersion,
		ApprouveAt: time.Now(),
	}
	if err := DB.Create(&c).Error; err != nil {
		http.Error(res, "db error", http.StatusInternalServerError)
		return
	}

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(http.StatusCreated)
	json.NewEncoder(res).Encode(c)
}

func revokeConsent(res http.ResponseWriter, req *http.Request) {
	userID := currentUserID(req)

	now := time.Now()
	result := DB.Model(&models.Consent{}).
		Where("user_id = ? AND revoked_at IS NULL", userID).
		Update("revoked_at", now)

	if result.RowsAffected == 0 {
		http.Error(res, "no active consent", http.StatusBadRequest)
		return
	}

	deleteUserVideos(userID)

	res.WriteHeader(http.StatusNoContent)
}
