package config

import (
	"log"

	"course-management-system/models"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

// ConnectDatabase initializes the SQLite database and performs migrations
func ConnectDatabase(dbPath string) *gorm.DB {
	if dbPath == "" {
		dbPath = "course_management.db"
	}

	database, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Enable foreign key constraints in SQLite
	database.Exec("PRAGMA foreign_keys = ON;")

	log.Println("Running AutoMigration for Course, Student, and Join Tables...")
	err = database.AutoMigrate(&models.Course{}, &models.Student{})
	if err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	DB = database
	log.Println("Database connection and migration completed successfully.")
	return DB
}
