package web

import (
	log "github.com/3fps/log2go"
	"github.com/qomun/shawty/data"
	"github.com/gorilla/mux"
	"net/http"
)

type ShortIDController struct {
	config map[string]string
	sh     data.Shawties
}

// NewShortIDController creates a new ShortIDController instance
func NewShortIDController(config map[string]string, sh data.Shawties) *ShortIDController {
	return &ShortIDController{config, sh}
}

func (ctrl *ShortIDController) Respond(shortID string) (res *ResPkg) {
	res = NewResPkg()

	// extract the ID
	id, random, err := data.FullID(shortID)

	if err != nil {
		res.HttpStatus = http.StatusNotFound
		res.Errors = append(res.Errors, err)
		log.Error("shortid_controller Respond error")
		log.Error(err)
		return
	}

	shawty, err := ctrl.sh.GetByID(id, random)
	if err != nil {
		res.HttpStatus = http.StatusNotFound
		res.Errors = append(res.Errors, err)
		log.Error("shortid_controller Respond error")
		log.Error(err)
		return
	}

	ctrl.sh.IncHits(shawty.ID) // increase hit

	res.Data["Domain"] = ctrl.config["SHAWTY_DOMAIN"]
	res.Data["Shawty"] = shawty
	res.HttpStatus = http.StatusMovedPermanently
	return
}

func (ctrl *ShortIDController) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer closeReqBody(r)

	vars := mux.Vars(r)
	shortID := vars["shortID"]
	res := ctrl.Respond(shortID)
	log.Infof("ServeHTTP for short ID %s", shortID)

	if res.HttpStatus == http.StatusMovedPermanently {
		s := res.Data["Shawty"].(*data.Shawty)
		log.Infof("Redirecting '%s' to '%s'", shortID, s.Url)
		http.Redirect(w, r, s.Url, http.StatusMovedPermanently)
	} else {
		http.NotFound(w, r)
		log.Infof("Cannot find shawty for '%s'", shortID)
	}
}
