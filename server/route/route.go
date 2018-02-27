package route

import (
	"net/http"
	"fmt"
	"os"
	"path/filepath"
	"html/template"
)

const (
	PassingHeader = "349fDdk340238dkfZp2191l6Jd60"
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

func Logout(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s\n", r.RequestURI)
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	fmt.Fprintln(w, "data:image/x-icon;base64,iVBORw0KGgoAAAANSUhEUgAAABAAAAAQEAYAAABPYyMiAAAABmJLR0T///////8JWPfcAAAACXBIWXMAAABIAAAASABGyWs+AAAAF0lEQVRIx2NgGAWjYBSMglEwCkbBSAcACBAAAeaR9cIAAAAASUVORK5CYII=\n")
}

func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s\n", r.RequestURI)

	switch r.Method {
	case http.MethodGet:
		break
	case http.MethodPost:
		break
	}

	w.Header().Set("Content-Type", "text/html; charset=UTF-8")

	renderTemplate(w, "login", nil)
}

func Favicon(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("%s\n", r.RequestURI)
	w.Header().Set("Content-Type", "image/x-icon")
	w.Header().Set("Cache-Control", "public, max-age=7776000")
	fmt.Fprintln(w, "data:image/png;base64,iVBORw0KGgoAAAANSUhEUgAAAEAAAABACAYAAACqaXHeAAADHElEQVR4Xu2ZTchNQRjHfy/yEbKwISILhYgVpWxElAULWdhYCSHJRooiXhshLCwsSUpYWSsryoJ8hqwoFPmM8tVfczWNc+49c86Zc8+5Z566m/POOzP/3zzPM/PMDNFyG2q5fiKA6AEtJxBDoOUOEJNgDIEYAi0nUDQExgATasLwN/AB+Okzn6IALgCrfAYM2FYAlgGPfcYoAmAG8BQY7TNg4LYLgPs+YxQBsAc45jNYBW3nAw98xskLYARwB1hkDfYVOAF885lAwbYLgfVWH5UBWAzccib/GpgKKBarsnXA1X4AOA3scFS+BWabTFwVgI2AEnHHKvGAScAjs9q20NYAmAs8TFji1gCYYwC4CTQCaEsOKOoBK4A1JoS0dR4FPjshNQrYC0w2328CVxLCri9JsCiAg8ABI+aXSaZvHHGqMV5aAM4BmwcFwD7giBGj4kVbp/KHbQLwBJhpPp4BdkYA/xNoZAhED4ghEHNATIKt3gX2A4dMQtfdgba6pHPAC6vgOgtsG5RtcAuw24j5CKwG3jnidM12HZhuvp8HDg8KAN0m6dexHymXByPh38uVToz6ubYBuGR9rOQ+oOhROEVvrs+7gJNtBaCC6iIw0QKgO8J7PijzXIrWwQMk/jIw1hL7DFiaUFd05dFEAEnin5tkKghe1jQAaeL1OiUI3lYmAO3l04C0rO49OecfShev/ssEoD19K/AlRel7QDc7to0HlmeYh8AeT4h5nSFyrXxnEmUC6LXCd52XJL0q61FDV2S+plgvLL5sD+gl4jawxDSqhfh+AaiN+KoBKP7ltteAlb3cJeHvpbm93XeeHDDP9wnaDKgT2isDwVd/EPF5PUDlq6ozFSvdbBYwJYNSvQ0oQaaZdpftRbN9Wud5PCCDpr9NtG11yt5u4tYCN7J2Wna7kABUpalaq634vCGQdRG6AZBb93XlOyJCesCplNec2ogP7QHD5mhse8wnYFM/Y95135AeoAPPOGfA74A8oDYWEkBtRHabSATQiGUKOMnoAQHhNqLr6AGNWKaAk4weEBBuI7qOHtCIZQo4yT+PGvxBK37qMwAAAABJRU5ErkJggg==")
	}

func Index(w http.ResponseWriter, r *http.Request) {
	if ph := r.Header.Get("tspy-pass"); ph != PassingHeader {
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

}