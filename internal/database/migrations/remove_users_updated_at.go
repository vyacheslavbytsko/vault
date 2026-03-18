package migrations

import (
	"vault/internal/database/models"

	"gorm.io/gorm"
)

func RemoveUsersUpdatedAt(db *gorm.DB) error {
	if !db.Migrator().HasColumn(&models.User{}, "updated_at") {
		return nil
	}

	return db.Migrator().DropColumn(&models.User{}, "updated_at")
}
