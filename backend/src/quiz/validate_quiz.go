package quiz

import (
	"encoding/json"
	"net/http"

	utils "profilsactifs/src/utils"

	"gorm.io/gorm"
)

// Validate calcule le score du profil et enregistre le résultat de certification.
// POST /quiz/valider  { "profile_id": 1 }
func Validate(gdb *gorm.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var body struct {
			ProfileID uint `json:"profile_id"`
		}
		json.NewDecoder(req.Body).Decode(&body)

		var reponses []utils.QuestionsAnswer
		gdb.Where("profile_id = ?", body.ProfileID).Find(&reponses)

		score := 0
		for _, rep := range reponses {
			if len(rep.Options) > 0 && rep.Options[0] == "Oui" {
				var q utils.Question
				gdb.First(&q, rep.QuestionID)
				score += q.Weight
			}
		}
		badge := score > 50

		resultat := utils.CertificationResult{
			ProfileID:   body.ProfileID,
			TotalScore:  score,
			BadgeEarned: badge,
		}
		gdb.Create(&resultat)

		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(resultat)
	}
}
