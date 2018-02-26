package server

import (
	"golang.org/x/crypto/acme/autocert"
	"net/http"
	"crypto/tls"
	"fmt"
	"flag"
	"time"
	"log"
	"os"
	"os/signal"
	"syscall"
	"context"
)

const (
	CertsDir = "certs"
	HostName = "yeosh.com"
	DevPort  = ":8000"
)

var (
	flgProduction = false
	testJson      = `{
"glossary": {
	"title": "example glossary",
	"GlossDiv": {
		"title": "S",
		"GlossList": {
			"GlossEntry": {
				"ID": "SGML",
				"SortAs": "SGML",
				"GlossTerm": "Standard Generalized Markup Language",
				"Acronym": "SGML",
				"Abbrev": "ISO 8879:1986",
				"GlossDef": {
					"para": "A meta-markup language, used to create markup languages such as DocBook.",
					"GlossSeeAlso": ["GML", "XML"]
				},
				"GlossSee": "markup"
			}
		}
	}
}
}`
)

func parseFlags() {
	flag.BoolVar(&flgProduction, "production", false, "if true, we start HTTPS server")
	flag.Parse()
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s\n", r.RequestURI)
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	fmt.Fprintln(w, "data:image/x-icon;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQEAYAAABPYyMiAAAABmJLR0T///////8JWPfcAAAACXBIWXMAAABIAAAASABGyWs+AAAAF0lEQVRIx2NgGAWjYBSMglEwCkbBSAcACBAAAeaR9cIAAAAASUVORK5CYII=\n")
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	for k, v := range r.Header {
		fmt.Printf("%v: %v\n", k, v)
	}

	w.Header().Add("tspy-token", "349fdk340238dkfp2191ld60")

	w.Write([]byte(testJson))
}

func makeServerFromMux(mux *http.ServeMux) *http.Server {
	// set timeouts so that a slow or malicious client doesn't
	// hold resources forever
	return &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      mux,
	}
}

func graceful(hs *http.Server, logger *log.Logger, timeout time.Duration) {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// block until signal
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger.Printf("\nShutdown with timeout: %s\n", timeout)

	if err := hs.Shutdown(ctx); err != nil {
		logger.Printf("Error: %v\n", err)
	} else {
		logger.Println("Server stopped")
	}
}

func Run() {
	parseFlags()
	var m *autocert.Manager
	var httpServ *http.Server

	log.Println("Server starting ...")

	mux := http.NewServeMux()

	// handlers
	mux.HandleFunc("/favicon.ico", faviconHandler)
	mux.HandleFunc("/", indexHandler)

	// configure and serve https
	httpServ = makeServerFromMux(mux)

	if flgProduction {
		// init autocert to automatically grab/cache/renew TLS certs for given hosts
		m = &autocert.Manager{
			Cache:      autocert.DirCache(CertsDir),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(HostName, "server."+HostName),
		}

		go func(){
			httpServ.Addr = ":https"
			httpServ.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
			log.Fatal(httpServ.ListenAndServeTLS("", ""))
		}()
	}

	if m != nil {
		// initial http-01 challenge and redirect http to httpS
		go func() {
			log.Fatal(http.ListenAndServe(":http", m.HTTPHandler(nil)))
		}()
	} else {
		go func() {
			fmt.Printf("Remember to run with -production for :443 on the server\n")
			httpServ.Addr = DevPort
			log.Fatal(httpServ.ListenAndServe())
		}()
	}

	graceful(httpServ, log.New(os.Stdout, "", 0), 5 * time.Second)
}
