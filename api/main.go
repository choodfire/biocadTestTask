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
)

var db *sql.DB

func getData(c *gin.Context) {
	tableName := os.Getenv("TABLENAME")

	unit_guid := c.Param("unit_guid")

	// SELECT * FROM maintable WHERE unit_guid = 'cold7_Defrost_status';
	res, err := db.Query("SELECT * FROM " + tableName + " WHERE unit_guid = '" + unit_guid + "'")
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

	fmt.Println(results)

	c.IndentedJSON(http.StatusOK, results)
}

func connectToDB(username, password, host, port, dbName string) {
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

	connectToDB(username, password, host, port, dbName)

	router := gin.Default()
	router.GET("/data/:unit_guid", getData)
	router.Run("localhost:8080")
}
