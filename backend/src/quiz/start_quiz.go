package quiz

import (
	"encoding/json"
	"fmt"
	"net/http"

	utils "profilsactifs/src/utils"

	"gorm.io/gorm"
)

// Start remet à zéro les réponses d'un profil.
// POST /quiz/start  { "profile_id": 1 }
func Start(gdb *gorm.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var body struct {
			ProfileID uint `json:"profile_id"`
		}
		json.NewDecoder(req.Body).Decode(&body)

		gdb.Where("profile_id = ?", body.ProfileID).Delete(&utils.QuestionsAnswer{})

		fmt.Fprintf(res, "Profile %d: quiz reset", body.ProfileID)
	}
}
