package profils

import (
	"encoding/json"
	"net/http"

	utils "jibjob/src/utils"

	"gorm.io/gorm"
)

const (
	defaultJob  = "Profil ProfilsActifs"
	defaultCity = "France"
)

type feedProfile struct {
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

// GET /profils -> [ { id, name, job, city, skills, role, ... }, ... ]
func List(gdb *gorm.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		res.Header().Set("Content-Type", "application/json")

		var profiles []utils.Profile
		if err := gdb.Order("created_at DESC, id DESC").Find(&profiles).Error; err != nil {
			res.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(res).Encode(map[string]string{"error": "internal error"})
			return
		}

		var videos []utils.Video
		gdb.Order("created_at DESC, id DESC").Find(&videos)
		videoByProfile := make(map[uint]string, len(videos))
		for _, v := range videos {
			if _, seen := videoByProfile[v.ProfileID]; !seen {
				videoByProfile[v.ProfileID] = v.URL
			}
		}

		out := make([]feedProfile, 0, len(profiles))
		for _, p := range profiles {
			skills := []string(p.Skills)
			if skills == nil {
				skills = []string{}
			}
			videoURL, hasVideo := videoByProfile[p.ID]
			role := string(p.Role)
			if role == "" {
				role = string(utils.Candidate)
			}
			out = append(out, feedProfile{
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

		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(out)
	}
}
