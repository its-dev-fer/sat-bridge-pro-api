package database

import (
	"app/src/utils"
	"os"
	"path/filepath"
	"strings"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	files, err := filepath.Glob("src/database/migrations/*.sql")
	if err != nil {
		utils.Log.Errorf("Error finding migration files: %v", err)
		return
	}
	for _, file := range files {
		if strings.HasSuffix(file, ".down.sql") {
			continue
		}
		content, err := os.ReadFile(file)
		if err != nil {
			utils.Log.Errorf("Error reading migration file %s: %v", file, err)
			return
		}
		if err := db.Exec(string(content)).Error; err != nil {
			utils.Log.Errorf("Error executing migration file %s: %v", file, err)
			return
		}
	}
}