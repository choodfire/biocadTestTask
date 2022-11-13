package main

import (
	"github.com/joho/godotenv"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	username := os.Getenv("DBUSERNAME")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	dbname := os.Getenv("DBNAME")
	directory := os.Getenv("DIRECTORY") // absolute path

}
