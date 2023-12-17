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

func check_err(err error) {
	if err != nil {
		ErrorLogger.Println(err)
	}
}

func init() {
	logFile, err := os.OpenFile(os.Getenv("PATH_LOG"), os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	InfoLogger = log.New(logFile, "INFO: ", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	DebugLogger = log.New(logFile, "DEBUG: ", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	ErrorLogger = log.New(logFile, "ERROR: ", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
}

func open_connection() *sql.DB {
	DB, err := sql.Open("postgres", conURL)
	if err != nil {
		ErrorLogger.Println(err)
	}
	return DB
}

func save_to_db(data map[string]string) {
	DB := open_connection()
	insert_string := `INSERT INTO Users (ip, bank_number)
		VALUES ($1, $2);
		`
	_, errs := DB.Exec(insert_string, data["ip"], data["bank_number"])
	check_err(errs)
	InfoLogger.Println("Saving in database successfully")
}

func savingHandler(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	check_err(err)
	data := make(map[string]string)
	err = json.Unmarshal(body, &data)
	check_err(err)
	DebugLogger.Printf("request from %s, args: ip : %s, bank_number :%s", r.RemoteAddr, data["ip"], data["bank_number"])
	save_to_db(data)
}

func main() {
	DB, err := sql.Open("postgres", conURL)
	check_err(err)
	_, err = DB.Exec("CREATE TABLE IF NOT EXISTS Users(ip varchar(20), bank_number varchar(16));")
	check_err(err)
	mux := http.NewServeMux()
	httpPort := os.Getenv("PORT")
	log.Printf("the database service is running")
	mux.HandleFunc("/saving", savingHandler)
	errs := http.ListenAndServe((":" + httpPort), mux)
	if errs != nil {
		ErrorLogger.Println(err)
	}
}
