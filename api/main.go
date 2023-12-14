package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	InfoLogger  *log.Logger
	DebugLogger *log.Logger
	ErrorLogger *log.Logger
)

type Validation struct {
	Valid  bool   `json:"valid"`
	Errors string `json:"errors"`
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

func saving_info(data map[string]string) {
	reqURL := "http://host.docker.internal:3335/saving"
	body, err := json.Marshal(data)
	if err != nil {
		ErrorLogger.Println(err)
	}
	client := http.Client{}
	_, err = client.Post(reqURL, "application/json", bytes.NewBuffer(body))
	if err != nil {
		ErrorLogger.Println(err)
	}

}

func web(w http.ResponseWriter, r *http.Request) {
	web_req := "http://host.docker.internal:3336"
	client := http.Client{}
	resp, err := client.Get(web_req)
	if err != nil {
		ErrorLogger.Println(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	w.Write([]byte(body))
}

func get_info(w http.ResponseWriter, r *http.Request) {
	card_number := r.URL.Query().Get("bank_number")
	data := make(map[string]string)
	data["number"] = card_number
	if data["number"] == "" {
		DebugLogger.Printf("%s request from %s, result :%t", r.Method, r.RemoteAddr, false)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("The bank card number is expected"))
	} else {
		// SENDING REQUEST
		requestURL := "http://host.docker.internal:3334/validation"
		body, err := json.Marshal(data)
		if err != nil {
			ErrorLogger.Println(err)
		}
		client := http.Client{}
		req, err := client.Post(requestURL, "application/json", bytes.NewBuffer(body))

		if err != nil {
			ErrorLogger.Println(err)
			return
		}

		// READ RESPONSE
		b, err := io.ReadAll(req.Body)
		if err != nil {
			ErrorLogger.Println(err)
		}
		req.Body.Close()
		m := map[string]interface{}{}
		// UMARSHALL DATA
		err = json.Unmarshal(b, &m)
		if err != nil {
			ErrorLogger.Println(err)
		}
		// PARSE FROM MAP
		valid := m["valid"]
		errsa := m["errors"]
		res, ok := valid.(bool)
		if !ok {
			ErrorLogger.Println(ok)
		}
		ers, ok := errsa.(string)
		if !ok {
			ErrorLogger.Println(ok)
		}
		// SAVING IN DB
		resp := Validation{Valid: res, Errors: ers}
		jsondata, err := json.Marshal(resp)
		if err != nil {
			ErrorLogger.Println(err)
		}
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsondata)
		InfoLogger.Printf("%s request, result: %t\n", r.Method, res)
		data := make(map[string]string)
		data["ip"] = r.RemoteAddr
		data["args"] = card_number
		if res {
			saving_info(data)
		}
	}
}

func main() {
	mux := http.NewServeMux()
	httpHost, httpPort := "0.0.0.0", os.Getenv("PORT")
	log.Printf("bank_card_validator-api starts on %s:%s", httpHost, httpPort)
	mux.HandleFunc("/", web)
	mux.HandleFunc("/check_number", get_info)
	err := http.ListenAndServe((httpHost + ":" + httpPort), mux)
	if err != nil {
		ErrorLogger.Println(err)
	}
}
