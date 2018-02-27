package main

import (
	"net/http"
	"time"
	"log"
	"os"
	"os/signal"
	"syscall"
	"context"
	"github.com/gorilla/mux"
	"TwitchSpy/server/route"
)

const (
	CertsDir = "certs"
	HostName = "yeosh.com"
	DevPort  = ":8000"
)

func graceful(hs *http.Server, logger *log.Logger, timeout time.Duration) {
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	// block until signal
	<-sigs

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger.Printf("\nShutdown with timeout: %s\n", timeout)

	if err := hs.Shutdown(ctx); err != nil {
		logger.Printf("Error: %v\n", err)
	} else {
		logger.Println("Server stopped")
	}
}

func assign(router *mux.Router) {
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(http.Dir("public"))))

	router.HandleFunc("/login", route.Login).Methods(http.MethodGet, http.MethodPost)
	router.HandleFunc("/logout", route.Logout).Methods(http.MethodPost)

	router.HandleFunc("/favicon.ico", route.Favicon).Methods(http.MethodGet)
	router.HandleFunc("/", route.Index).Methods(http.MethodGet, http.MethodPost)
}

func main() {
	parseFlags()
	var hs *http.Server

	router := mux.NewRouter()

	assign(router)

	hs = setup(router)

	graceful(hs, log.New(os.Stdout, "", 0), 5 * time.Second)
}
