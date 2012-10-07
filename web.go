package main

import "html/template"
import "net/http"
import "regexp"
import "os"
import "code.google.com/p/gorilla/mux"
import "strings"
import "net/url"
import "time"

var (
	indexHtml = template.Must(template.ParseFiles("index.html"))
	shawtyJs  = template.Must(template.ParseFiles("shawty.js"))
)
var urlPattern = regexp.MustCompile("^(?i)(https?|ftp|file)://.+$")
var domain = os.Getenv("SHAWTY_DOMAIN")

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	if err := indexHtml.Execute(w, nil); err != nil {
		Lerror("Cannot execute index template")
		Lerror(err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	Linfo("Index requested")
}

func HandleShawtyJS(w http.ResponseWriter, r *http.Request) {

	var u = r.FormValue("url")
	if u == "" {
		u = r.Referer()
		if u != "" {
			// if it's in the referer then we don't want the browser to cache this
			w.Header().Set("Cache-Control", "no-cache")
			http.Redirect(w, r, "shawty.js?url=" + url.QueryEscape(u), http.StatusTemporaryRedirect)
			return 
		}
	}

	if !urlPattern.MatchString(u) {
		http.NotFound(w, r)
		return
	}

	data := map[string]interface{} {
		"Success": 1,
		"Message": "",
		"Long": "",
		"Short": "",
		"Timestamp": time.Now(),
		"Hits": 0,
	}
	code := http.StatusOK

	var sh, err = NewShawties()
	if err != nil {
		Lerror(err)
		data["Success"] = 0
		data["Message"] = "unknown error occured"
		code = http.StatusInternalServerError
	}
	defer sh.Close()

	if (code == http.StatusOK) {
		if strings.HasPrefix(strings.ToLower(u), strings.ToLower(domain)) {
			Lerrorf("Do not try to shorten itself: %s", u)
			data["Success"] = 0
			data["Message"] = "link is already shortened"
			code = http.StatusBadRequest
		} else {
			s, err := sh.GetOrCreate(u)
			if err != nil {
				Lerror(err)
				data["Success"] = 0
				data["Message"] = "cannot shorten this link"
				code = http.StatusNotFound
			} else {
				data["Long"] = s.Url
				data["Short"] = domain + ShortID(s.ID, s.Rand)
				data["Timestamp"] = s.CreatedOn
				data["Hits"] = s.Hits

				Linfo("ShawtyJS requested")
				Linfof("URL: %s", u)
				Linfof("Shawty: %v", s)
			}
		}
	}

	w.Header().Set("Content-Type", "application/javascript")
	// w.WriteHeader(code)
	if err := shawtyJs.Execute(w, data); err != nil {
		Lerror("Cannot execute shawty javascript template")
		Lerror(err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
}

func HandleShortID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	shortID := vars["shortID"]
	id, random, err := FullID(shortID)
	if err != nil {
		http.NotFound(w, r)
		return
	}

	sh, err := NewShawties()
	if err != nil {
		Lerror(err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	defer sh.Close()

	s, err := sh.GetByID(id, random)
	if err != nil {
		Linfof("Cannot find shawty for '%s'", shortID)
		http.NotFound(w, r)
		return
	}

	sh.IncHits(id)

	Linfof("Redirecting '%s' to '%s'", shortID, s.Url)
	http.Redirect(w, r, s.Url, http.StatusMovedPermanently)
}
