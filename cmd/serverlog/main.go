package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"

	"main.go/internal/config"
)

const (
	envLocal = "local"
	envProd  = "prod"
)

func main() {

	cfg := config.MustLoad()
	log := setupLogger(cfg.Env)
	log.Info("starting application")

	http.HandleFunc("/return", handler)
	err := listenAndServe(fmt.Sprintf("%d", cfg.Port), cfg.TlsPath)
	if err != nil {
		panic(err)
	}
}

// collects incoming request and sends it to the output
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
}

// setup level of logger info
func setupLogger(env string) *slog.Logger {
	var log *slog.Logger
	switch env {
	case envLocal:
		log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}
	return log
}

// listens to a port, choosing between a secure or unsecured connection
func listenAndServe(port, tlsPath string) error {
	var err error
	port = ":" + port
	if tlsPath == "not exist" {
		err = http.ListenAndServe(port, nil)
	} else {
		certFile := tlsPath + "fullchain.pem"
		keyFile := tlsPath + "privkey.pem"
		err = http.ListenAndServeTLS(port, certFile, keyFile, nil)
	}
	return err
}
