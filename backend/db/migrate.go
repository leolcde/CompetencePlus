package db

import (
	"profilsactifs/models"

	"gorm.io/gorm"
)

func AutoMigrate(g *gorm.DB) error {
	return g.AutoMigrate(
		&models.User{},
		&models.Video{},
		&models.Question{},
		&models.Answer{},
		&models.BadgeResult{},
		&models.Consent{},
	)
}
