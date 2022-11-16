package main

import (
	"biocadTestTask/data"
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

var checkedFiles []string
var logs []data.LogRow

func contains[T comparable](arr []T, val T) bool {
	for _, elem := range arr {
		if elem == val {
			return true
		}
	}

	return false
}

func parseTSV(filePath string) []data.LogRow {
	newLogs := []data.LogRow{}

	file, err := os.Open(filePath)
	if err != nil {
		panic(err)
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		args := strings.Split(line, "\t") // 14 things

		n := len(args)
		if n < 14 { // if not enough arguments in tsv
			for i := n; i < 14; i++ {
				if i == 0 || i == 8 {
					args = append(args, "0")
				}
				args = append(args, "") // to make sure insert won't break
			}
		}

		thisN, _ := strconv.Atoi(args[0])
		thisLevel, _ := strconv.Atoi(args[8])

		currentLog := data.LogRow{
			N:         thisN,
			Mqtt:      args[1],
			Invid:     args[2],
			Unit_guid: args[3],
			Msg_id:    args[4],
			Text:      args[5],
			Context:   args[6],
			Class:     args[7],
			Level:     thisLevel,
			Area:      args[9],
			Addr:      args[10],
			Block:     args[11],
			Typee:     args[12],
			Bit:       args[13],
		}

		logs = append(logs, currentLog)

		newLogs = append(newLogs, currentLog)
	}

	return newLogs
}

func addToDB(newLogs []data.LogRow, tableName string, db *sql.DB) {
	for _, currentLog := range newLogs {
		_, _ = db.Exec("INSERT INTO "+tableName+" (n, mqtt, invid, unit_guid, msg_id, text, context, class, level, area, addr, block, type, bit) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", currentLog.N, currentLog.Mqtt, currentLog.Invid, currentLog.Unit_guid, currentLog.Msg_id, currentLog.Text,
			currentLog.Context, currentLog.Class, currentLog.Level, currentLog.Area, currentLog.Addr, currentLog.Block, currentLog.Typee, currentLog.Bit)
	}
}

func logsToFile(newLogs []data.LogRow) {
	if _, err := os.Stat("./output"); os.IsNotExist(err) {
		// folder output doesn't exist

		err := os.Mkdir("./output", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, newLog := range newLogs {
		currentFilename := "./output/" + newLog.Unit_guid + ".doc"

		file, err := os.OpenFile(currentFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		_, err = file.Write([]byte(fmt.Sprintf("%d %s %s %s %s %s %s %s %d %s %s %s %s %s\n",
			newLog.N, newLog.Mqtt, newLog.Invid, newLog.Unit_guid, newLog.Msg_id, newLog.Text, newLog.Context,
			newLog.Class, newLog.Level, newLog.Area, newLog.Addr, newLog.Block, newLog.Typee, newLog.Bit)))
		if err != nil {
			file.Close()
			log.Fatal(err)
		}

		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}
}

func connectToDB(username, password, host, port, dbName string) *sql.DB {
	// connect to db
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbName))
	if err != nil {
		// no db connection
		log.Fatal(err)
	}

	// check if connection is successful
	err = db.Ping()
	if err != nil {
		// no db connection
		log.Fatal(err)
	}

	// See "Important settings" section. todo check
	//db.SetConnMaxLifetime(time.Minute * 3)
	//db.SetMaxOpenConns(10)
	//db.SetMaxIdleConns(10)

	return db
}

func main() {
	err := godotenv.Load()
	if err != nil {
		// no credentials
		log.Fatal(err)
	}

	// get credentials from env
	username := os.Getenv("DBUSERNAME")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")
	dbName := os.Getenv("DBNAME")
	tableName := os.Getenv("TABLENAME")
	directory := os.Getenv("DIRECTORY") // absolute path

	db := connectToDB(username, password, host, port, dbName)
	defer db.Close()

	// loop from here
	for true { // периодический осмотр директории
		files, err := os.ReadDir(directory)
		if err != nil {
			// wrong directory
			log.Fatal(err)
		}

		newLogs := []data.LogRow{}
		for _, file := range files {
			if contains(checkedFiles, file.Name()) { // if file already checked
				continue
			}

			newLogs = append(newLogs, parseTSV(directory+"\\"+file.Name())...)
		}

		addToDB(newLogs, tableName, db)

		logsToFile(newLogs)

		time.Sleep(1 * time.Minute)
	}
}
