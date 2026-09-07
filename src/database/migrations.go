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
		for _, stmt := range strings.Split(string(content), ";") {
			stmt = strings.TrimSpace(stmt)
			if stmt == "" {
				continue
			}
			if err := db.Exec(stmt).Error; err != nil {
				utils.Log.Errorf("Error executing migration file %s: %v", file, err)
			}
		}
	}
}