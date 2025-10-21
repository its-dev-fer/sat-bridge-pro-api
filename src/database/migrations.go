package database

import (
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"app/src/utils"
)

func RunMigrations(db *gorm.DB) {
	files, err := filepath.Glob("src/database/migrations/*.sql")
	if err != nil {
		utils.Log.Errorf("Error finding migration files: %v", err)
		return
	}
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			utils.Log.Errorf("Error reading migration file %s: %v", file, err)
			return
		}
		// Execute the SQL commands in the migration file
		if err := db.Exec(string(content)).Error; err != nil {
			utils.Log.Errorf("Error executing migration file %s: %v", file, err)
			return
		}
	}
}