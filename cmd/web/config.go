package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var connectDBAttempts = 3

func openDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connectToDB() *gorm.DB {
	host := os.Getenv("POSTGRES_HOST")
	user := os.Getenv("POSTGRES_USER")
	password := os.Getenv("POSTGRES_PASSWORD")
	dbname := os.Getenv("POSTGRES_DB")
	port := os.Getenv("POSTGRES_PORT")
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s",
		host, user, password, dbname, port)
	app.InfoLog.Println(dsn)

	i := connectDBAttempts
	for {
		conn, err := openDB(dsn)
		if err != nil {
			log.Println("Postgsres is not ready...")
			i--
		} else {
			return conn
		}
		if i == 0 {
			log.Println(err)
			return nil

		}
		time.Sleep(2 * time.Second)
	}
}
