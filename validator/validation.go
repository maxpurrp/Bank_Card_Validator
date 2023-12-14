package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	lynh "validator/validation"
)

type Validation struct {
	Valid  bool   `json:"valid"`
	Errors string `json:"errors"`
}

var (
	InfoLogger  *log.Logger
	DebugLogger *log.Logger
	ErrorLogger *log.Logger
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

func handler(w http.ResponseWriter, r *http.Request) {
	b, err := io.ReadAll(r.Body)
	if err != nil {
		ErrorLogger.Println(err)
	}
	m := make(map[string]string)
	err = json.Unmarshal(b, &m)
	if err != nil {
		ErrorLogger.Println(err)
	}
	DebugLogger.Printf("%s request from %s, args: %s", r.Method, r.RemoteAddr, m)
	if len(m["number"]) < 16 {
		resp := Validation{Valid: false, Errors: "Bad card number"}
		jsondata, err := json.Marshal(resp)
		if err != nil {
			ErrorLogger.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsondata)
		InfoLogger.Printf("%s request, result: %t\n", r.Method, false)
	} else {
		res, er := lynh.AlgLynh(m["number"])
		resp := Validation{Valid: res, Errors: er}
		jsondata, err := json.Marshal(resp)
		if err != nil {
			ErrorLogger.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsondata)
		InfoLogger.Printf("%s request, result: %t\n", r.Method, res)
	}
}

func main() {
	mux := http.NewServeMux()
	httpHost, httpPort := "0.0.0.0", os.Getenv("PORT")
	log.Printf("bank_card_validator-validator starts on %s:%s", httpHost, httpPort)
	mux.HandleFunc("/validation", handler)
	err := http.ListenAndServe((httpHost + ":" + httpPort), mux)
	if err != nil {
		ErrorLogger.Println(err)
	}

}
