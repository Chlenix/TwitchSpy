package main

import (
	"net/http"
	"golang.org/x/crypto/acme/autocert"
	"time"
	"crypto/tls"
	"fmt"
	"log"
	"flag"
	"github.com/julienschmidt/httprouter"
	"os"
	"os/signal"
	"syscall"
	"context"
)

func parseFlags() {
	flag.BoolVar(&flgProduction, "production", false, "if true, we start HTTPS server")
	flag.Parse()
}

func makeServerFromMux(router *httprouter.Router) *http.Server {
	// set timeouts so that a slow or malicious client doesn't
	// hold resources forever
	return &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		Handler:      router,
	}
}

func setup(router *httprouter.Router) (*http.Server) {
	var m *autocert.Manager
	var hs *http.Server

	hs = makeServerFromMux(router)

	if flgProduction {
		// init autocert to automatically grab/cache/renew TLS certs for given hosts
		m = &autocert.Manager{
			Cache:      autocert.DirCache(CertsDir),
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(conf.HostName, "www."+conf.HostName),
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
			hs.Addr = conf.Port
			log.Fatal(hs.ListenAndServe())
		}()
	}
	return hs
}

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