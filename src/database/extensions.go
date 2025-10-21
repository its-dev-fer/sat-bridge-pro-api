package database

import (
	"gorm.io/gorm"
	"os"
	"path/filepath"
	"app/src/utils"
)

func DefineExtensions(db *gorm.DB) {
	files, err := filepath.Glob("src/database/extensions/*.sql")
	if err != nil {
		utils.Log.Errorf("Error finding extension files: %v", err)
		return
	}
	for _, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			utils.Log.Errorf("Error reading extension file %s: %v", file, err)
			return
		}
		// Execute the SQL commands in the extension file
		if err := db.Exec(string(content)).Error; err != nil {
			utils.Log.Errorf("Error executing extension file %s: %v", file, err)
			return
		}
	}
}