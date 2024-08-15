package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"zappy.sh/models"
)

var (
	DBConn *gorm.DB
)

func ConnectDB() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=require", os.Getenv("PGHOST"), os.Getenv("PGUSER"), os.Getenv("PGPASSWORD"), os.Getenv("PGDATABASE"), os.Getenv("PGPORT"))
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to db: ", err)
		os.Exit(2)
	}

	log.Println("Connected to database")

	if err := db.AutoMigrate(&models.Alias{}, &models.Request{}); err != nil {
		log.Fatal("Failed to migrate db: ", err)
	}

	DBConn = db
}
