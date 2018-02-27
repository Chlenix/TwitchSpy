package main

import (
	"net/http"
	"golang.org/x/crypto/acme/autocert"
	"github.com/gorilla/mux"
	"time"
	"crypto/tls"
	"fmt"
	"log"
	"flag"
)

var (
	flgProduction = false
)

func parseFlags() {
	flag.BoolVar(&flgProduction, "production", false, "if true, we start HTTPS server")
	flag.Parse()
}

func makeServerFromMux(router *mux.Router) *http.Server {
	// set timeouts so that a slow or malicious client doesn't
	// hold resources forever
	return &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      router,
	}
}

func setup(router *mux.Router) (*http.Server) {
	var m *autocert.Manager
	var hs *http.Server

	hs = makeServerFromMux(router)

	if flgProduction {
		// init autocert to automatically grab/cache/renew TLS certs for given hosts
		m = &autocert.Manager{
			Cache:      autocert.DirCache(CertsDir),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(HostName, "server."+HostName),
		}

		go func(){
			hs.Addr = ":https"
			hs.TLSConfig = &tls.Config{GetCertificate: m.GetCertificate}
			log.Fatal(hs.ListenAndServeTLS("", ""))
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
			hs.Addr = DevPort
			log.Fatal(hs.ListenAndServe())
		}()
	}
	return hs
}