package database

import (
	"log"
	"os"
	"to-do-backend/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Dbinstance struct {
	Db *gorm.DB
}

var Database Dbinstance

func ConnectDb() {
	db, err := gorm.Open(sqlite.Open("api.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to DB! \n", err.Error())
		os.Exit(2)
	}
	log.Print("Connected to the database successfully")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running Migrations")

	// Add Migrations
	db.AutoMigrate(&models.User{}, &models.Task{})
	Database = Dbinstance{Db: db}
}
