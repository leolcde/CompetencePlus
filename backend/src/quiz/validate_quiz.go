package quiz

import (
	"encoding/json"
	"net/http"

	utils "profilsactifs/src/utils"

	"gorm.io/gorm"
)

func Validate(gdb *gorm.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var body struct {
			ProfileID uint `json:"profile_id"`
		}
		json.NewDecoder(req.Body).Decode(&body)

		var answers []utils.QuestionsAnswer
		gdb.Where("profile_id = ?", body.ProfileID).Find(&answers)

		score := 0
		for _, ans := range answers {
			if len(ans.Options) > 0 && ans.Options[0] == "Oui" {
				var q utils.Question
				gdb.First(&q, ans.QuestionID)
				score += q.Weight
			}
		}

		var maxScore int
		gdb.Model(&utils.Question{}).Select("COALESCE(SUM(weight), 0)").Scan(&maxScore)
		badge := maxScore > 0 && score*100 >= maxScore*60

		result := utils.CertificationResult{
			ProfileID:   body.ProfileID,
			TotalScore:  score,
			BadgeEarned: badge,
		}
		gdb.Create(&result)

		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(result)
	}
}
