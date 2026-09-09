package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"

	"profilsactifs/models"

	"gorm.io/gorm/clause"
)

func listQuestions(res http.ResponseWriter, req *http.Request) {
	var questions []models.Question
	DB.Find(&questions)
	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(questions)
}

func startQuestionnaire(res http.ResponseWriter, req *http.Request) {
	userID := currentUserID(req)

	DB.Where("user_id = ?", userID).Delete(&models.Answer{})

	fmt.Fprintf(res, "Questionnaire reset for user %d", userID)
}

func submitAnswer(res http.ResponseWriter, req *http.Request) {
	userID := currentUserID(req)

	var body struct {
		QuestionID uint   `json:"question_id"`
		Choice     string `json:"choice"`
	}
	if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
		http.Error(res, "invalid body", http.StatusBadRequest)
		return
	}

	var q models.Question
	if err := DB.First(&q, body.QuestionID).Error; err != nil {
		http.Error(res, "unknown question", http.StatusBadRequest)
		return
	}

	if !slices.Contains(q.Options, body.Choice) {
		http.Error(res, "invalid choice", http.StatusBadRequest)
		return
	}

	answer := models.Answer{
		UserID:     userID,
		QuestionID: body.QuestionID,
		Choice:     body.Choice,
	}

	if err := DB.Create(&answer).Error; err != nil {
		http.Error(res, "could not save answer (already answered?)", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(res, "Answer saved for question %d", body.QuestionID)
}

func validateQuestionnaire(res http.ResponseWriter, req *http.Request) {
	userID := currentUserID(req)

	var answers []models.Answer
	DB.Where("user_id = ?", userID).Find(&answers)

	score := 0
	for _, answer := range answers {
		if answer.Choice == "Yes" {
			var q models.Question
			DB.First(&q, answer.QuestionID)
			score += q.Weight
		}
	}

	result := models.BadgeResult{
		UserID: userID,
		Score:  score,
		Badge:  score > 50,
	}
	DB.Clauses(clause.OnConflict{UpdateAll: true}).Create(&result)

	res.Header().Set("Content-Type", "application/json")
	json.NewEncoder(res).Encode(result)
}
