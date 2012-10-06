package main

import "html/template"
import "net/http"
import "regexp"
import "os"

var (
	indexHtml = template.Must(template.ParseFiles("index.html"))
	shawtyJs = template.Must(template.ParseFiles("shawty.js"))
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
	var url = r.FormValue("url")
	if !urlPattern.MatchString(url) {
		http.NotFound(w, r)
		return
	}

	var sh, err = NewShawties()
	if err != nil {
		Lerror(err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}
	defer sh.Close()

	s, err := sh.GetOrCreate(url)
	if err != nil {
		Lerror(err)
		http.NotFound(w, r)
	}

	var data = map[string]interface{}{
		"Short": domain + ShortID(s.ID, s.Rand),
		"Hits":  s.Hits,
		"Timestamp": s.CreatedOn,
	}

	if err := shawtyJs.Execute(w, data); err != nil {
		Lerror("Cannot execute shawty javascript template")
		Lerror(err)
		http.Error(w, "Error", http.StatusInternalServerError)
		return
	}

	Linfo("ShawtyJS requested")
	Linfof("URL: %s", url)
	Linfof("Shawty: %v", s)
}
