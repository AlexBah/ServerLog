package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var mu sync.Mutex
var count int

func main() {
	http.HandleFunc("/return", handler)
	http.HandleFunc("/count", counter)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	requestString := fmt.Sprintf("%s %s %s\n", r.Method, r.URL, r.Proto)
	for k, v := range r.Header {
		requestString += fmt.Sprintf("%q : %q\n", k, v)
	}
	requestString += fmt.Sprintf("\nHost = %q\nRemoteAddr = %q\n\n", r.Host, r.RemoteAddr)
	if err := r.ParseForm(); err != nil {
		log.Print(err)
	}
	for k, v := range r.Form {
		requestString += fmt.Sprintf("Form[%q] = %q\n", k, v)
	}
	fmt.Fprintf(w, "%s", requestString)
	writeLogFile(requestString)
	mu.Lock()
	count++
	mu.Unlock()
}

func counter(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	fmt.Fprintf(w, "Count = %d\n", count)
	mu.Unlock()
}

func writeLogFile(requestString string) {
	mu.Lock()
	fileName := time.Now().GoString() + ".log"
	mu.Unlock()
	os.Chdir("Log")
	logFile, err := os.Create(fileName)
	fmt.Fprintf(logFile, "%s", requestString)
	logFile.Close()
	os.Chdir("..")
	if err != nil {
		log.Print(err)
	}

}
