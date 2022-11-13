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

func parseTSV(filePath string) {
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

		// insert row to db
		addToDB(currentLog)
	}
}

func addToDB(currentLog logRow) {
	_, _ = db.Exec("INSERT INTO data"+
		" (n, mqtt, invid, unit_guid, msg_id, text, context, class, level, area, addr, block, type, bit) "+
		"VALUES (?, ?, ?)", currentLog.n, currentLog.mqtt, currentLog.invid, currentLog.unit_guid, currentLog.msg_id, currentLog.text,
		currentLog.context, currentLog.class, currentLog.level, currentLog.area, currentLog.addr, currentLog.block, currentLog.typee, currentLog.bit)
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
	dbname := os.Getenv("DBNAME")
	directory := os.Getenv("DIRECTORY") // absolute path

	// connect to db
	db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, dbname))
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

	for _, file := range files {
		if contains(checkedFiles, file.Name()) { // if file already checked
			continue
		}
		parseTSV(directory + "\\" + file.Name())
	}
}
