package db

import (
	"gpt-load/internal/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// V1_3_0_AddTokensColumns adds tokens-related columns to request_logs table
func V1_3_0_AddTokensColumns(db *gorm.DB) error {
	// Check if migration is needed
	if db.Migrator().HasColumn(&models.RequestLog{}, "prompt_tokens") {
		logrus.Info("Tokens columns already exist, skipping v1.3.0...")
		return nil
	}

	logrus.Info("Running migration v1.3.0: Adding tokens columns to request_logs table...")

	if err := db.AutoMigrate(&models.RequestLog{}); err != nil {
		return err
	}

	logrus.Info("Migration v1.3.0 completed successfully")
	return nil
}
