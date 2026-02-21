package database

import (
	"log"
	"os"
	"todoApp/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Db *gorm.DB

func InitDB() {
	//Check i fthe database url variable exists (render)
	connStr := os.Getenv("DATABASE_URL")

    //If it's empty,it is running on the local host
    if connStr == "" {
        connStr = "user=postgres host=localhost password=password dbname=todoApp port=5432 sslmode=disable"
        log.Println("Connecting to Local Database...")
    } else {
        log.Println("Connecting to Render Database...")
    }

	var err error
	Db, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		panic (err)
	}
	//migrate the tables
	err = Db.AutoMigrate(&models.User{}, &models.Todo{})
	if err != nil {
		log.Fatal("error migrating table", err)
	}
}