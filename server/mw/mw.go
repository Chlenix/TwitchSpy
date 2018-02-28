package mw

import (
	"net/http"
	"log"
	"github.com/julienschmidt/httprouter"
)

type Middleware struct {
	next http.HandlerFunc
	message string
}

func NewMiddleware(next http.HandlerFunc) *Middleware {
	return &Middleware{next: next, message: "Hello from mw!"}
}

func (mw *Middleware) Authorize(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// We can modify the request here; for simplicity, we will just log a message
	log.Printf("msg: %s, Method: %s, URI: %s\n", mw.message, r.Method, r.RequestURI)
	mw.next.ServeHTTP(w, r)
	// We can modify the response here
}