package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"time"

	"golang.org/x/crypto/acme/autocert"
	"flag"
	"fmt"
)

const (
	CertsDir = "certs"
	HostName = "yeosh.com"
	DevPort  = ":8080"
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

func main() {
	parseFlags()
	var m *autocert.Manager
	var server *http.Server

	log.Println("Server starting ...")

	mux := http.NewServeMux()

	// handlers
	mux.HandleFunc("/favicon.ico", faviconHandler)
	mux.HandleFunc("/", indexHandler)

	// configure and serve https
	server = makeServerFromMux(mux)

	if flgProduction {
		// init autocert to automatically grab/cache/renew TLS certs for given hosts
		m = &autocert.Manager{
			Cache:      autocert.DirCache(CertsDir),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(HostName, "www."+HostName),
		}

		// initial http-01 challenge and redirect http to httpS
		go http.ListenAndServe(":http", m.HTTPHandler(nil))

		server.Addr = ":https"
		server.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
		log.Fatal(server.ListenAndServeTLS("", ""))
	}

	fmt.Printf("Remember to run with -production for :443 on the server\n")
	// Local debugging
	server.Addr = DevPort
	log.Fatal(server.ListenAndServe())

}
