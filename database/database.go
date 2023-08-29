package database

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"zappy.sh/models"
)

var (
	DBConn *gorm.DB
)

func ConnectDB() {
	var err error
	dsn := fmt.Sprintf("%s&parseTime=True", os.Getenv("DSN"))
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		TranslateError:                           true,
	})

	if err != nil {
		log.Fatal("Failed to connect to db: ", err)
		os.Exit(2)
	}

	log.Println("Connected to database")

	if err := db.AutoMigrate(&models.Request{}, &models.Alias{}); err != nil {
		log.Fatal("Failed to migrate db: ", err)
	}

	DBConn = db
}
