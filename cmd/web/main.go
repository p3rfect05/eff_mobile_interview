package main

import (
	"log"
	"net/http"
	"os"

	"gorm.io/gorm"
)

type AppConfig struct {
	InfoLog  *log.Logger
	ErrorLog *log.Logger
	DB       *gorm.DB
}

var app AppConfig

func main() {

	infoLog := log.New(os.Stdout, "[INFO]\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stdout, "[ERROR]\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog
	app.InfoLog = infoLog

	conn := connectToDB()
	if conn == nil {
		app.ErrorLog.Fatalln("Could not connect to database")
	}

	app.DB = conn
	app.InfoLog.Println("Connected to database")
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router(),
	}

	InitConn(app.DB) // чтобы впоследствии в файле models.go обращаться через db.
	_, err := MigrateModels()
	if err != nil {
		app.ErrorLog.Fatalln("Failed to migrate tables...")
	}
	app.InfoLog.Println("Starting the server on port 80")
	srv.ListenAndServe()
}
