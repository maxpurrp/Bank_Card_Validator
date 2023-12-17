package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"text/template"
)

var (
	InfoLogger  *log.Logger
	DebugLogger *log.Logger
	ErrorLogger *log.Logger
)

func check(err error) {
	ErrorLogger.Println(err)
}

func init() {
	logFile, err := os.OpenFile("../logs.log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	InfoLogger = log.New(logFile, "INFO: ", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	DebugLogger = log.New(logFile, "DEBUG: ", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
	ErrorLogger = log.New(logFile, "ERROR: ", log.LstdFlags|log.Lshortfile|log.Lmicroseconds)
}

type Data struct {
	Error string
}

func web(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./web/templates/index.html")
	check(err)
	t.Execute(w, nil)
}

func successHalndler(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("./web/templates/success.html")
	check(err)
	t.Execute(w, nil)
}

func failHalndler(w http.ResponseWriter, r *http.Request) {

	tmpl := template.Must(template.ParseFiles("./web/templates/failed.html"))
	body, err := io.ReadAll(r.Body)
	check(err)
	var data map[string]string
	err = json.Unmarshal(body, &data)
	check(err)
	res := Data{Error: data["err"]}
	err = tmpl.Execute(w, res)
	check(err)
}

func main() {
	mux := http.NewServeMux()
	httpPort := os.Getenv("PORT")
	log.Printf("the web service is running")
	mux.HandleFunc("/", web)
	mux.HandleFunc("/success", successHalndler)
	mux.HandleFunc("/failed", failHalndler)
	err := http.ListenAndServe((":" + httpPort), mux)
	check(err)
}
