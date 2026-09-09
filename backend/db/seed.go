package db

import (
	"log"

	"profilsactifs/models"

	"github.com/lib/pq"
	"gorm.io/gorm"
)

func SeedQuestions(g *gorm.DB) error {
	var count int64
	if err := g.Model(&models.Question{}).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return nil
	}

	yesOrNo := pq.StringArray{"Yes", "No"}
	contents := []string{
		"I communicate easily with my colleagues.",
		"I am comfortable speaking in public.",
		"I can express myself clearly in writing.",
		"I ask questions when I do not understand.",
		"I express myself with confidence.",
		"I stay calm when facing disagreement.",
		"I respect agreed schedules.",
		"I organize my working time efficiently.",
		"I meet the deadlines I am given.",
		"I prioritize my tasks efficiently.",
		"I handle several tasks at the same time.",
		"I am rigorous in my work.",
		"I organize myself without constant supervision.",
		"I adapt quickly to a new environment.",
		"I accept constructive criticism.",
		"I stay calm in stressful situations.",
		"I work well under pressure.",
		"I adapt to organizational changes.",
		"I handle the unexpected without panicking.",
		"I am ready to step out of my comfort zone.",
	}

	questions := make([]models.Question, len(contents))
	for i, c := range contents {
		questions[i] = models.Question{Content: c, Options: yesOrNo, Weight: 1}
	}

	if err := g.Create(&questions).Error; err != nil {
		return err
	}
	log.Printf("seed: %d questions added", len(questions))
	return nil
}
