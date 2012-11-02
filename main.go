package main

import (
	"runtime"
	"code.google.com/p/gorilla/mux"
	"fmt"
	"go.3fps.com/shawty/data"
	"go.3fps.com/shawty/utils"
	"go.3fps.com/shawty/web"
	"go.3fps.com/utils/log"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// read configurations
	confKeys := []string{"SHAWTY_PORT", "SHAWTY_DB", "SHAWTY_DOMAIN", "SHAWTY_MODE", "SHAWTY_LPM", "SHAWTY_LOG_DIR"}
	config := make(map[string]string)
	for _, k := range confKeys {
		config[k] = os.Getenv(k)
	}

	// setup logger
	log.SetDir(config["SHAWTY_LOG_DIR"])

	// setup data
	random := utils.NewBestRand()
	shawties, err := data.NewMySh(random, config["SHAWTY_DB"])
	if err != nil {
		log.Error("Cannot create MySh")
		return
	}
	defer shawties.Close()

	// register routes
	home := web.NewHomeController(config)
	shawtyjs := web.NewShawtyJSController(config, shawties)
	shortID := web.NewShortIDController(config, shawties)

	// setup HTTP server
	router := mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	router.Handle("/", home)
	router.Handle("/shawty.js", shawtyjs)
	router.Handle("/{shortID:[A-Za-z0-9]+}", shortID)

	var port = config["SHAWTY_PORT"]
	if port == "" {
		port = "80"
	}

	l, err := net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Errorf("Cannot listen at %s", port)
		fmt.Println(err)
		return
	}
	defer l.Close()
	log.Infof("Listening at %s", port)

	runMode := config["SHAWTY_MODE"]
	switch runMode {
	case "fcgi":
		fcgi.Serve(l, router)
	default:
		http.Serve(l, router)
	}
}
