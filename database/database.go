package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var Database *gorm.DB

func Connect() error {
	var err error
	host := os.Getenv("DB_HOST")
	username := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	databaseName := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Africa/Lagos", host, username, password, databaseName, port)
	Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	//create database if does not exist
	if err != nil {
		// Try to connect without specifying a database name
		dsn = fmt.Sprintf("host=%s user=%s password=%s port=%s sslmode=disable TimeZone=Africa/Lagos", host, username, password, port)
		Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err != nil {
			return err
		}

		// Create the database
		err = Database.Exec(fmt.Sprintf("CREATE DATABASE %s;", databaseName)).Error

		if err != nil {
			return err
		}

		// Close the connection and reconnect with the new database
		sqlDB, err := Database.DB()

		if err != nil {
			return err
		}

		sqlDB.Close()

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Africa/Lagos", host, username, password, databaseName, port)
		Database, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

		if err != nil {
			return err
		}

	}

	return nil
}
