package main

import (
	"net/http"
	"time"
	"log"
	"os"
	"github.com/julienschmidt/httprouter"
	"TwitchSpy/server/route"
	"TwitchSpy/server/api"
	"TwitchSpy/server/mw"
	"github.com/kelseyhightower/envconfig"
	"TwitchSpy/config"
)

const (
	CertsDir       = "certs"
	StaticFilesDir = "public"
	ConfigPrefix   = "server"
)

var (
	hs *http.Server
	conf *config.ServerConfig
	flgProduction = false
)

func handleRoutes(router *httprouter.Router) {
	middleware := mw.NewMiddleware

	router.ServeFiles("/public/*filepath", http.Dir(StaticFilesDir))

	router.POST("/api/login", middleware(api.Login, conf).Login)
	router.POST("/api/logout", api.Logout)

	router.GET("/login", route.Login)
	router.GET("/", middleware(route.Index, conf).Auth)
}

func main() {
	parseFlags()
	conf = &config.ServerConfig{}

	if err := envconfig.Process(ConfigPrefix, conf); err != nil {
		panic(err)
	}

	router := httprouter.New()

	handleRoutes(router)

	hs = setup(router)

	graceful(hs, log.New(os.Stdout, "", 0), 5*time.Second)
}
