package profils

import (
	"encoding/json"
	"errors"
	"net/http"

	utils "jibjob/src/utils"

	"gorm.io/gorm"
)

type detailProfile struct {
	ID          uint     `json:"id"`
	Name        string   `json:"name"`
	Job         string   `json:"job"`
	City        string   `json:"city"`
	Skills      []string `json:"skills"`
	Role        string   `json:"role"`
	IsCertified bool     `json:"isCertified"`
	Score       int      `json:"score"`
	VideoUrl    string   `json:"videoUrl"`
	HasConsent  bool     `json:"hasConsent"`
}

// Detail renvoie un seul profil.
// GET /profils/{id} -> { id, name, job, city, skills, role, ... }
func Detail(gdb *gorm.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		id := req.PathValue("id")

		var p utils.Profile
		if err := gdb.First(&p, "id = ?", id).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				res.WriteHeader(http.StatusNotFound)
				json.NewEncoder(res).Encode(map[string]string{"error": "not found"})
				return
			}
			res.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(res).Encode(map[string]string{"error": "internal error"})
			return
		}

		var video utils.Video
		videoURL := ""
		hasVideo := false
		if err := gdb.Where("profile_id = ?", p.ID).Order("created_at DESC, id DESC").First(&video).Error; err == nil {
			videoURL = video.URL
			hasVideo = true
		}

		skills := []string(p.Skills)
		if skills == nil {
			skills = []string{}
		}
		role := string(p.Role)
		if role == "" {
			role = string(utils.Candidate)
		}

		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(detailProfile{
			ID:          p.ID,
			Name:        utils.Or(p.Name, "Profil sans nom"),
			Job:         utils.Or(p.Sector, defaultJob),
			City:        utils.Or(p.Location, defaultCity),
			Skills:      skills,
			Role:        role,
			IsCertified: false,
			Score:       0,
			VideoUrl:    videoURL,
			HasConsent:  hasVideo,
		})
	}
}
