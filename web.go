package main

import (
	"code.google.com/p/gorilla/mux"
	"go.3fps.com/utils/log"
	"html/template"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"
)

var (
	indexHtml = template.Must(template.ParseFiles("templates/index.html"))
	shawtyJs  = template.Must(template.ParseFiles("templates/shawty.js"))
)
var urlPattern = regexp.MustCompile("^(?i)(https?|ftp|file)://.+$")
var domain = os.Getenv("SHAWTY_DOMAIN")

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	data := map[string]interface{}{"Domain": domain}

	if err := indexHtml.Execute(w, data); err != nil {
		log.Error("Cannot execute index template")
		log.Error(err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	log.Info("Index requested")
}

func HandleShawtyJS(w http.ResponseWriter, r *http.Request) {

	var u = r.FormValue("url")
	if u == "" {
		u = r.Referer()
		if u != "" {
			// if it's in the referer then we don't want the browser to cache this
			w.Header().Set("Cache-Control", "no-cache")
			http.Redirect(w, r, "shawty.js?url="+url.QueryEscape(u), http.StatusTemporaryRedirect)
			return
		}
	}

	data := map[string]interface{}{
		"Bookmarklet": false,
		"Success":     1,
		"Message":     "",
		"Long":        "",
		"Short":       "",
		"Timestamp":   time.Now(),
		"Hits":        0,
	}
	code := http.StatusOK

	// bookmarklet param
	data["Bookmarklet"] = r.FormValue("bm") == "1"

	if !urlPattern.MatchString(u) {
		data["Success"] = 0
		data["Message"] = "a valid url has to start with http:// or https://"
		code = http.StatusBadRequest
	}

	if code == http.StatusOK {
		var sh, err = NewShawties()
		if err != nil {
			log.Error(err)
			data["Success"] = 0
			data["Message"] = "unknown error occured"
			code = http.StatusInternalServerError
		}
		defer sh.Close()

		if code == http.StatusOK {
			if strings.HasPrefix(strings.ToLower(u), strings.ToLower(domain)) {
				log.Errorf("Do not try to shorten itself: %s", u)
				data["Success"] = 0
				data["Message"] = "link is already shortened"
				code = http.StatusBadRequest
			} else {
				s, err := sh.GetOrCreate(u)
				if err != nil {
					log.Error(err)
					data["Success"] = 0
					data["Message"] = "cannot shorten this link"
					code = http.StatusNotFound
				} else {
					data["Long"] = s.Url
					data["Short"] = domain + ShortID(s.ID, s.Rand)
					data["Timestamp"] = s.CreatedOn
					data["Hits"] = s.Hits

					log.Info("ShawtyJS requested")
					log.Infof("URL: %s", u)
					log.Infof("Shawty: %v", s)
				}
			}
		}
	}

	w.Header().Set("Content-Type", "application/javascript")
	if err := shawtyJs.Execute(w, data); err != nil {
		log.Error("Cannot execute shawty javascript template")
		log.Error(err)
		http.Error(w, "Error", http.StatusInternalServerError)
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
		log.Error(err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	defer sh.Close()

	s, err := sh.GetByID(id, random)
	if err != nil {
		log.Infof("Cannot find shawty for '%s'", shortID)
		http.NotFound(w, r)
		return
	}

	sh.IncHits(id)

	log.Infof("Redirecting '%s' to '%s'", shortID, s.Url)
	http.Redirect(w, r, s.Url, http.StatusMovedPermanently)
}
