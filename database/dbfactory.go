package database

import (
	"fmt"
	"github.com/dilly3/doctors-portal/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"os"
)

func SetupDB() *gorm.DB {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	password := os.Getenv("DB_PASSWORD")
	dbDatabase := os.Getenv("DB")
	root := os.Getenv("DB_ROOTS")
	port := os.Getenv("DB_PORT")
	dsn := fmt.Sprintf("host=localhost user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai", root, password, dbDatabase, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect database: ", err)

	}

	//This creates our table for this model
	err = db.AutoMigrate(&models.Doctor{}, &models.Patient{}, &models.Appointment{})
	if err != nil {
		log.Fatalf("failed to migrate: %v", err)
	}
	return db
}
