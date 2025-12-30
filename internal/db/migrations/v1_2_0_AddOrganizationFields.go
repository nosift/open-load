package db

import (
	"gpt-load/internal/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// V1_2_0_AddOrganizationFields adds organization-related fields to api_keys table
func V1_2_0_AddOrganizationFields(db *gorm.DB) error {
	// Check if migration is needed
	var needMigrateCount int64

	// Check if the column already exists
	if db.Migrator().HasColumn(&models.APIKey{}, "is_organization_key") {
		logrus.Info("Organization fields already exist, skipping v1.2.0...")
		return nil
	}

	logrus.Info("Running migration v1.2.0: Adding organization fields to api_keys table...")

	// Add the new columns
	type APIKey struct {
		IsOrganizationKey bool   `gorm:"not null;default:false"`
		OrganizationID    string `gorm:"type:varchar(255);default:''"`
		OrganizationName  string `gorm:"type:varchar(255);default:''"`
	}

	if err := db.AutoMigrate(&models.APIKey{}); err != nil {
		return err
	}

	// Count keys that need to be checked
	db.Model(&models.APIKey{}).Count(&needMigrateCount)

	if needMigrateCount > 0 {
		logrus.Infof("Migration v1.2.0: Added organization fields. %d existing keys will be checked on next validation.", needMigrateCount)
	}

	logrus.Info("Migration v1.2.0 completed successfully")
	return nil
}
