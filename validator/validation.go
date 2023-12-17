package main

import (
	"bytes"
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

type RequestBodyWeb struct {
	Err string `json:"err"`
}

type RequestBodyDB struct {
	Ip          string `json:"ip"`
	Bank_number string `json:"bank_number"`
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

func check(err error) {
	if err != nil {
		ErrorLogger.Println(err)
	}
}

func send_resp(res bool, err string) []byte {
	var url string
	switch {
	case res:
		url = "http://web:3336/success"
	default:
		url = "http://web:3336/failed"
	}
	data := RequestBodyWeb{Err: err}
	body, er := json.Marshal(data)
	check(er)
	resp, errs := http.Post(url, "application/json", bytes.NewBuffer(body))
	check(errs)
	defer resp.Body.Close()
	body, ers := io.ReadAll(resp.Body)
	check(ers)
	return body
}

func handler(w http.ResponseWriter, r *http.Request) {
	var resp []byte
	param := r.URL.Query().Get("bank_number")
	DebugLogger.Printf("%s request from %s, args: %s", r.Method, r.RemoteAddr, param)
	if len(param) < 16 {
		resp = send_resp(false, "The number of characters in the bank card must be 16")
		InfoLogger.Printf("%s request, result: %t\n", r.Method, false)
	} else {
		res, err := lynh.AlgLynh(param)
		InfoLogger.Printf("%s request, result: %t\n", r.Method, res)
		if res {
			resp = send_resp(res, err)
			url := "http://db:3335/saving"
			data := RequestBodyDB{Ip: r.RemoteAddr,
				Bank_number: param}
			reqBody, err := json.Marshal(data)
			check(err)

			_, err = http.Post(url, "application/json", bytes.NewBuffer(reqBody))
			check(err)
		} else {
			resp = send_resp(res, err)
		}
		InfoLogger.Printf("%s request, result: %t\n", r.Method, res)
	}
	w.Write(resp)
}

func main() {
	mux := http.NewServeMux()
	httpPort := os.Getenv("PORT")
	log.Printf("the validation service is running")
	mux.HandleFunc("/validation", handler)
	err := http.ListenAndServe((":" + httpPort), mux)
	check(err)
}
