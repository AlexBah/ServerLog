package handler

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"time"
)

// listens to a port, choosing between a secure or unsecured connection
func ListenPort(port, tlsPath string, shutdownCh <-chan struct{}, log *slog.Logger, timeout time.Duration) {
	srv := &http.Server{Addr: port, Handler: http.HandlerFunc(handlerReturn)}
	log.Info(fmt.Sprintf("Starting listen on port %s", srv.Addr))

	if tlsPath == "not exist" {
		go func() {
			if err := srv.ListenAndServe(); err != nil {
				log.Error("Port", srv.Addr, err)
			}
		}()
	} else {
		go func() {
			certFile := tlsPath + "fullchain.pem"
			keyFile := tlsPath + "privkey.pem"
			if err := srv.ListenAndServeTLS(certFile, keyFile); err != nil {
				log.Error("Port", srv.Addr, err)
			}
		}()
	}

	stopListen(srv, shutdownCh, log, timeout)
}

// stop listen port, then come signal close application
func stopListen(srv *http.Server, shutdownCh <-chan struct{}, log *slog.Logger, timeout time.Duration) {
	go func() {
		<-shutdownCh

		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		log.Info(fmt.Sprintf("Shutting down server on port %s", srv.Addr))
		if err := srv.Shutdown(ctx); err != nil {
			log.Error("Server shutdown failed on port", srv.Addr, err)
		} else {
			log.Info(fmt.Sprintf("Server shutdown gracefully on port %s ", srv.Addr))
		}
	}()
}

// collects incoming request and sends it to the output
func handlerReturn(w http.ResponseWriter, r *http.Request) {
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
