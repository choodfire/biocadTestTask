package main

import (
	"database/sql"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		// shut down cause no info
		log.Fatal(err)
	}

	// get credentials from env
	username := os.Getenv("DBUSERNAME")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	dbname := os.Getenv("DBNAME")
	directory := os.Getenv("DIRECTORY") // absolute path

	// connect to db
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname))
	if err != nil {
		// shut down no db connection
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		// shut down no db connection
		log.Fatal(err)
	}

	// See "Important settings" section. todo check
	//db.SetConnMaxLifetime(time.Minute * 3)
	//db.SetMaxOpenConns(10)
	//db.SetMaxIdleConns(10)

}
