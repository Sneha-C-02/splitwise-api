package config

import (
	"log"
	"splitwise-api/models"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	database, err := gorm.Open(sqlite.Open("splitwise.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!")
	}

	DB = database

	// Auto-migrate all models
	database.AutoMigrate(
		&models.User{},
		&models.Group{},
		&models.GroupMember{},
		&models.Expense{},
		&models.ExpenseSplit{},
	)

	log.Println("Database connected & migrated successfully 🚀")
}
