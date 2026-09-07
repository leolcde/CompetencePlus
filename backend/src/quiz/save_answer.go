package quiz

import (
	"encoding/json"
	"fmt"
	"net/http"

	utils "jibjob/src/utils"

	"gorm.io/gorm"
)

// SaveAnswer enregistre la réponse d'un profil à une question.
// POST /quiz/answer  { "profile_id": 1, "question_id": 2, "options": ["Oui"] }
func SaveAnswer(gdb *gorm.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var body struct {
			ProfileID  uint     `json:"profile_id"`
			QuestionID uint     `json:"question_id"`
			Options    []string `json:"options"`
		}
		json.NewDecoder(req.Body).Decode(&body)

		reponse := utils.QuestionsAnswer{
			ProfileID:  body.ProfileID,
			QuestionID: body.QuestionID,
			Options:    body.Options,
		}

		if err := gdb.Create(&reponse).Error; err != nil {
			fmt.Println("Erreur d'insertion :", err)
			http.Error(res, "answer not save", http.StatusBadRequest)
			return
		}

		fmt.Fprintf(res, "answer save for question: %d", body.QuestionID)
	}
}
