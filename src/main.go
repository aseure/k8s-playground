package main

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ctx := getContext()

	mux := http.NewServeMux()
	mux.Handle("/live", http.HandlerFunc(liveHandler))
	mux.Handle("/ready", http.HandlerFunc(readyHandler))
	mux.Handle("/", http.HandlerFunc(rootHandler))

	server := &http.Server{Addr: ":8080", Handler: mux}

	go func() {
		_ = server.ListenAndServe()
	}()

	<-ctx.Done()
	_ = server.Shutdown(context.Background())
}

func liveHandler(w http.ResponseWriter, r *http.Request) {
	log(r)
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "live")
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	log(r)
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "ready")
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log(r)
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, "OK")
}

func log(r *http.Request) {
	fmt.Printf("[%s] %s - %s\n", time.Now().Format(time.RFC3339), r.Method, r.URL.Path)
}

func getContext() context.Context {
	signals := make(chan os.Signal, 2)

	signal.Notify(signals, syscall.SIGTERM, syscall.SIGKILL)
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		s := <-signals
		switch s {
		case syscall.SIGTERM:
			fmt.Println("Stopping application: received SIGTERM")
		case syscall.SIGKILL:
			fmt.Println("Stopping application: received SIGKILL")
		}
		cancel()
	}()

	return ctx
}
