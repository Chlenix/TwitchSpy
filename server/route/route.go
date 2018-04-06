package route

import (
	"net/http"
	"fmt"
	"os"
	"path/filepath"
	"html/template"
	"github.com/julienschmidt/httprouter"
)

const (
	HtmlExtension = ".gohtml"
)

func renderTemplate(w http.ResponseWriter, tmpl string, data interface{}) {
	cwd, _ := os.Getwd()
	t, err := template.ParseFiles(filepath.Join(cwd, "./view/" + tmpl + HtmlExtension))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = t.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Printf("%s\n", r.RequestURI)

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	renderTemplate(w, "login", nil)
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func NotFound(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}