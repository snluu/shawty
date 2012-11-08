package web

import (
	"errors"
	"fmt"
	"go.3fps.com/shawty/data"
	"go.3fps.com/utils/log"
	"go.3fps.com/shawty/utils"
	"html/template"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

var urlPattern = regexp.MustCompile("^(?i)(https?)://.+$")
var shawtyJs *template.Template

func getShawtyJs() *template.Template {
	if shawtyJs == nil {
		shawtyJs = template.Must(template.ParseFiles("templates/shawty.js"))
	}
	return shawtyJs
}

type ShawtyJSController struct {
	config map[string]string
	sh     data.Shawties
}

func NewShawtyJSController(config map[string]string, sh data.Shawties) *ShawtyJSController {
	return &ShawtyJSController{config, sh}
}

// GetJSResponse response to javascript request, based on the URL and the bookmarklet flag
func (controller *ShawtyJSController) GetJSResponse(url, creatorIP string, bm bool) (res *ResPkg) {
	domain := controller.config["SHAWTY_DOMAIN"]
	limit, _ := strconv.ParseUint(controller.config["SHAWTY_LPM"], 0, 32)

	res = NewResPkg()
	res.Data["Bookmarklet"] = bm
	res.Data["Domain"] = domain
	res.Data["Success"] = 1

	if url == "a long URL" {
		res.Data["Success"] = 0
		res.Errors = append(res.Errors, errors.New("look at you, the funniest person alive. you just won the internet."))
	} else if !urlPattern.MatchString(url) {
		res.Data["Success"] = 0
		res.Errors = append(res.Errors, errors.New("a valid url has to start with http:// or https://"))
	} else {
		dupDomainPattern := regexp.MustCompile(fmt.Sprintf("^(?i)(https?)://%s.*", domain))
		if dupDomainPattern.MatchString(url) {
			res.Data["Success"] = 0
			res.Errors = append(res.Errors, errors.New("this url is already shortened"))
		} else {

			shawty, err := controller.sh.GetByUrl(url)

			if err != nil { // need to create
				minuteAgo := time.Unix(time.Now().Unix()-60, 0)
				numLinks, err := controller.sh.NumLinks(creatorIP, minuteAgo)

				if err != nil {
					res.Data["Success"] = 0
					res.Errors = append(res.Errors, errors.New("unknown error occured"))
					log.Errorf("%v", err)
				} else {

					// check for rate limit
					if numLinks >= uint32(limit) {
						res.Data["Success"] = 0
						res.Errors = append(res.Errors, errors.New("you are shortening your links too quickly"))
						log.Errorf("%s has reached the rate limit (%d)", creatorIP, limit)
					} else {
						shawty, err = controller.sh.Create("", url, creatorIP)
						if err != nil {
							res.Data["Success"] = 0
							res.Errors = append(res.Errors, errors.New("cannot shorten this link"))
							log.Errorf("%v", err)
						}
					}
				}
			}

			if res.Data["Success"] == 1 {
				res.Data["Shawty"] = shawty
				res.Data["Short"] = data.ShortID(shawty.ID, shawty.Rand)
			}
		}
	}

	return
}

func (controller *ShawtyJSController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer closeReqBody(r)
	
	var u = r.FormValue("url")
	if u == "" {
		u = r.Referer()
		if u != "" {
			// if it's in the referer then we don't want the browser to cache this
			w.Header().Set("Cache-Control", "no-cache")
			http.Redirect(w, r, "shawty.js?url="+url.QueryEscape(u), http.StatusTemporaryRedirect)
			return
		} else {
			http.NotFound(w, r)
			return
		}
	}

	res := controller.GetJSResponse(u, utils.HttpRemoteIP(r), r.FormValue("bm") == "1")	
	tpl := getShawtyJs()
	w.Header().Set("Content-Type", "application/javascript")
	if err := tpl.Execute(w, res); err != nil {
		log.Error("Cannot execute shawty javascript template")
		log.Error(err)
		http.Error(w, "Error", http.StatusInternalServerError)
	}
}
