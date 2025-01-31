package dao

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectDatabase() {

	username := os.Getenv("db_user")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("db_name")
	dbHost := os.Getenv("db_host")

	dbUri := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=5432 sslmode=disable", username, password, dbName, dbHost)
	fmt.Println(dbUri)
	conn, err := gorm.Open(postgres.Open(dbUri), &gorm.Config{})
	if err != nil {
		fmt.Print(err)
	}

	db = conn
}

func GetDB() *gorm.DB {
	return db
}

func SetDb(testdb *gorm.DB) {
	db = testdb
}
