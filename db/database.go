package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

var (
	InfoLogger  *log.Logger
	DebugLogger *log.Logger
	ErrorLogger *log.Logger
	conURL      = fmt.Sprintf("host=postgres port=5432 user=%s password=%s dbname=%s sslmode=disable", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_DB"))
)

func init() {
	logFile, err := os.OpenFile(os.Getenv("PATH_LOG"), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	InfoLogger = log.New(logFile, "INFO: ", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	DebugLogger = log.New(logFile, "DEBUG: ", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	ErrorLogger = log.New(logFile, "ERROR: ", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
}

func saving(data map[string]string) {
	DB, err := sql.Open("postgres", conURL)
	if err != nil {
		ErrorLogger.Println(err)
	}
	defer DB.Close()

	err = DB.Ping()
	if err != nil {
		ErrorLogger.Println(err)
	} else {
		create_table := "CREATE TABLE IF NOT EXISTS Users(ip varchar(20), bank_number varchar(16));"
		_, err := DB.Exec(create_table)
		if err != nil {
			ErrorLogger.Println(err)
		}
		insert_string := "INSERT INTO Users (ip, bank_number) VALUES ($1, $2);"
		_, errs := DB.Exec(insert_string, data["ip"], data["args"])
		if errs != nil {
			ErrorLogger.Println(err)
		}
		DebugLogger.Println("Saving in database successfully")

	}
}

func get_data(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		ErrorLogger.Println(err)
	}
	m := make(map[string]string)
	err = json.Unmarshal(b, &m)
	if err != nil {
		ErrorLogger.Println(err)
	}
	DebugLogger.Printf("request from %s, args: %s", r.RemoteAddr, m)
	saving(m)
}

func main() {
	mux := http.NewServeMux()
	httpHost, httpPort := "0.0.0.0", os.Getenv("PORT")
	log.Printf("bank_card_validator-database starts on %s:%s", httpHost, httpPort)
	mux.HandleFunc("/saving", get_data)
	err := http.ListenAndServe((httpHost + ":" + httpPort), mux)
	if err != nil {
		ErrorLogger.Println(err)
	}
}
