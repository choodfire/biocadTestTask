package main

import (
	"biocadTestTask/data"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"strconv"
)

type Env struct {
	db *sql.DB
}

func (e *Env) getData(c *gin.Context) {
	tableName := os.Getenv("TABLENAME")

	unit_guid := c.Param("unit_guid")

	// SELECT * FROM maintable WHERE unit_guid = 'cold7_Defrost_status';
	res, err := e.db.Query("SELECT * FROM `" + tableName + "` WHERE unit_guid = '" + unit_guid + "'")
	if err != nil {
		log.Fatal(err)
	}

	results := []data.LogRow{}
	for res.Next() {
		var currentLog data.LogRow

		err := res.Scan(&currentLog.N, &currentLog.Mqtt, &currentLog.Invid, &currentLog.Unit_guid, &currentLog.Msg_id, &currentLog.Text, &currentLog.Context,
			&currentLog.Class, &currentLog.Level, &currentLog.Area, &currentLog.Addr, &currentLog.Block, &currentLog.Typee, &currentLog.Bit)

		if err != nil {
			log.Fatal(err)
		}

		results = append(results, currentLog)
	}

	limit := 10
	pageStr := c.Param("page")
	page, _ := strconv.Atoi(pageStr)

	lowerBound := limit*page - (limit - 1)
	upperBound := limit*page + 1

	if upperBound > len(results) {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Page not found."})
	}
	c.IndentedJSON(http.StatusOK, results[lowerBound:upperBound])
}

func connectToDB(username, password, host, port, dbName string) *sql.DB {
	// connect to db
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbName))
	if err != nil {
		// shut down no db connection
		log.Fatal(err)
	}

	// check if connection is successful
	err = db.Ping()
	if err != nil {
		// shut down no db connection
		log.Fatal(err)
	}

	return db
}

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
	dbName := os.Getenv("DBNAME")

	db := connectToDB(username, password, host, port, dbName)
	defer db.Close()

	router := gin.Default()
	env := &Env{db: db}

	router.GET("/data/:unit_guid/:page", env.getData)
	router.Run("localhost:8080")
}
