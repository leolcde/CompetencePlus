package quiz

import (
	"encoding/json"
	"net/http"

	utils "profilsactifs/src/utils"

	"gorm.io/gorm"
)

// List renvoie toutes les questions du quiz.
// GET /quiz
func List(gdb *gorm.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var questions []utils.Question
		gdb.Find(&questions)

		res.Header().Set("Content-Type", "application/json")
		json.NewEncoder(res).Encode(questions)
	}
}
