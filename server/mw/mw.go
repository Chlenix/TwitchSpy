package mw

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"github.com/gorilla/sessions"
	_ "github.com/gorilla/securecookie"
	"TwitchSpy/config"
)

//var secureStore = sessions.NewCookieStore(securecookie.GenerateRandomKey(64))
var secureStore = sessions.NewCookieStore([]byte("secret"))

const (
	// Cookies
	SessionToken  = "session-id"
	MaxSessionAge = 7776000

	// Headers
	AccessToken = "access-token"

	TempKey = "3929efdkvk23-ekvdsas02i9"
)

type Middleware struct {
	next   http.HandlerFunc
	config *config.ServerConfig
}

func NewMiddleware(next http.HandlerFunc, config *config.ServerConfig) *Middleware {
	return &Middleware{next: next, config: config}
}

func (mw *Middleware) AuthAPI(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	session, _ := secureStore.Get(r, SessionToken)

	if session.IsNew {
		http.Error(w, "No Session", http.StatusBadRequest)
		return
	}

	accessToken := session.Values[AccessToken].(string)
	if headerToken := r.Header.Get(AccessToken); headerToken != accessToken {
		http.Error(w, "Bad Access Token", http.StatusBadRequest)
		return
	}

	mw.next.ServeHTTP(w, r)
}

func (mw *Middleware) Auth(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	session, err := secureStore.Get(r, SessionToken)

	if err != nil {
		http.Error(w, err.Error(), http.StatusForbidden)
		return
	}

	if session.IsNew {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// TODO: check vs real database
	if session.Values["Password"] != TempKey {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	mw.next.ServeHTTP(w, r)
}

func (mw *Middleware) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// We can modify the request here; for simplicity, we will just log a message
	session, _ := secureStore.Get(r, SessionToken)

	if session.IsNew {
		session.Options = &sessions.Options{
			Path:     "/",
			Secure:   mw.config.Secure,
			HttpOnly: true,
			MaxAge:   MaxSessionAge,
			Domain:   "",
		}
	}

	r.ParseForm()

	session.Values["Password"] = r.FormValue("pswd")
	session.Values["Token"] = "TEST_TEST_TOKEN!"

	session.Save(r, w)

	mw.next.ServeHTTP(w, r)

	// We can modify the response here
}
