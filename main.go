package main

import (
	"errors"
	"fmt"
	log "github.com/3fps/log2go"
	"shawty/data"
	"shawty/utils"
	"shawty/web"
	"github.com/gorilla/mux"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"runtime"
	"strings"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	// read configurations
	confKeys := []string{
		"SHAWTY_PORT", "SHAWTY_DB_TYPE", "SHAWTY_DB", "SHAWTY_DOMAIN",
		"SHAWTY_MODE", "SHAWTY_LPM", "SHAWTY_LOG_DIR",
	}
	config := make(map[string]string)
	for _, k := range confKeys {
		config[k] = os.Getenv(k)
	}

	// setup logger
	log.SetDir(config["SHAWTY_LOG_DIR"])

	// setup data
	random := utils.NewBestRand()
	dbType := strings.ToLower(config["SHAWTY_DB_TYPE"])
	var shawties data.Shawties
	var err error

	if dbType == "pg" || dbType == "postgres" || dbType == "postgresql" {
		shawties, err = data.NewPgSh(random, config["SHAWTY_DB"])
	} else if dbType == "mysql" {
		shawties, err = data.NewMySh(random, config["SHAWTY_DB"])
	} else {
		err = errors.New("Unknown database type. Please check $SHAWTY_DB_TYPE")
	}
	if err != nil {
		log.Error("Cannot create data source")
		log.Error(err)
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

	runMode := strings.ToLower(config["SHAWTY_MODE"])

	switch runMode {
	case "fcgi":
		log.Info("Serving FCGI mode")
		fcgi.Serve(l, router)
	default:
		log.Info("Serving HTTP mode")
		http.Serve(l, router)
	}
}
