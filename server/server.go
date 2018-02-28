package main

import (
	"net/http"
	"time"
	"log"
	"os"
	"os/signal"
	"syscall"
	"context"
	"github.com/julienschmidt/httprouter"
	"TwitchSpy/server/route"
	"TwitchSpy/server/api"
	"TwitchSpy/server/mw"
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

func handleRoutes(router *httprouter.Router) {
	middleware := mw.NewMiddleware

	router.ServeFiles("/public/*filepath", http.Dir("public"))

	router.POST("/api/login", middleware(api.Login).Authorize)
	router.POST("/api/logout", api.Logout)

	router.GET("/login", route.Login)
	router.GET("/", route.Index)
}

func main() {
	parseFlags()
	var hs *http.Server

	router := httprouter.New()

	handleRoutes(router)

	hs = setup(router)

	graceful(hs, log.New(os.Stdout, "", 0), 5 * time.Second)
}
