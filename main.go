package main

import (
	"bufio"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

var checkedFiles []string
var logs []logRow
var db *sql.DB

type logRow struct {
	n         int
	mqtt      string
	invid     string
	unit_guid string // globally unique identifier
	msg_id    string
	text      string
	context   string
	class     string
	level     int
	area      string
	addr      string
	block     string
	typee     string // can't use word 'type' in Go
	bit       string
}

func contains[T comparable](arr []T, val T) bool {
	for _, elem := range arr {
		if elem == val {
			return true
		}
	}

	return false
}

func parseTSV(filePath string) []logRow {
	newLogs := []logRow{}

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

		currentLog := logRow{
			n:         thisN,
			mqtt:      args[1],
			invid:     args[2],
			unit_guid: args[3],
			msg_id:    args[4],
			text:      args[5],
			context:   args[6],
			class:     args[7],
			level:     thisLevel,
			area:      args[9],
			addr:      args[10],
			block:     args[11],
			typee:     args[12],
			bit:       args[13],
		}

		logs = append(logs, currentLog)

		newLogs = append(newLogs, currentLog)
	}

	return newLogs
}

func addToDB(newLogs []logRow, tableName string) {
	for _, currentLog := range newLogs {
		_, _ = db.Exec("INSERT INTO "+tableName+" (n, mqtt, invid, unit_guid, msg_id, text, context, class, level, area, addr, block, type, bit) "+
			"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", currentLog.n, currentLog.mqtt, currentLog.invid, currentLog.unit_guid, currentLog.msg_id, currentLog.text,
			currentLog.context, currentLog.class, currentLog.level, currentLog.area, currentLog.addr, currentLog.block, currentLog.typee, currentLog.bit)
	}
}

func logsToFile(newLogs []logRow) {
	if _, err := os.Stat("./output"); os.IsNotExist(err) {
		// folder output doesn't exist

		err := os.Mkdir("./output", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}

	for _, newLog := range newLogs {
		currentFilename := "./output/" + newLog.unit_guid + ".doc"

		file, err := os.OpenFile(currentFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		_, err = file.Write([]byte(fmt.Sprintf("%d %s %s %s %s %s %s %s %d %s %s %s %s %s\n",
			newLog.n, newLog.mqtt, newLog.invid, newLog.unit_guid, newLog.msg_id, newLog.text, newLog.context,
			newLog.class, newLog.level, newLog.area, newLog.addr, newLog.block, newLog.typee, newLog.bit)))
		if err != nil {
			file.Close()
			log.Fatal(err)
		}

		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
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
	tableName := os.Getenv("TABLENAME")
	directory := os.Getenv("DIRECTORY") // absolute path

	// connect to db
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbName))
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

	// See "Important settings" section. todo check
	//db.SetConnMaxLifetime(time.Minute * 3)
	//db.SetMaxOpenConns(10)
	//db.SetMaxIdleConns(10)

	files, err := os.ReadDir(directory)
	if err != nil {
		// wrong directory
		log.Fatal(err)
	}

	newLogs := []logRow{}
	for _, file := range files {
		if contains(checkedFiles, file.Name()) { // if file already checked
			continue
		}

		newLogs = append(newLogs, parseTSV(directory+"\\"+file.Name())...)
	}

	addToDB(newLogs, tableName)

	logsToFile(newLogs)
}
