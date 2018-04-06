package api

import (
	"net/http"
	"fmt"
	"github.com/julienschmidt/httprouter"
)

const (
	DefaultContentType = "application/json; charset=utf-8"
)

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s\n", r.RequestURI)

	w.Header().Set("Content-Type", DefaultContentType)
	w.WriteHeader(http.StatusOK)
}

func Logout(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// TODO: clear sessions
	http.Redirect(w, r, "/", http.StatusFound)
}
