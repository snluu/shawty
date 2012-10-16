package main

import (
	"code.google.com/p/gorilla/mux"
	"fmt"
	"go.3fps.com/utils/log"
	"math/rand"
	"net"
	"net/http"
	"net/http/fcgi"
	"os"
	"strings"
	"time"
)

func main() {

	rand.Seed(time.Now().Unix())

	var router = mux.NewRouter()
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	router.HandleFunc("/", HandleIndex)
	router.HandleFunc("/shawty.js", HandleShawtyJS)
	router.HandleFunc("/{shortID:[A-Za-z0-9]+}", HandleShortID)

	var port = os.Getenv("SHAWTY_PORT")
	if port == "" {
		port = "80"
	}

	var l, err = net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		log.Errorf("Cannot listen at %s", port)
		fmt.Println(err)
		return
	}
	defer l.Close()
	log.Infof("Listening at %s", port)

	runMode := strings.ToLower(os.Getenv("SHAWTY_MODE"))
	switch runMode {
	case "fcgi":
		fcgi.Serve(l, router)
	default:
		http.Serve(l, router)
	}
}
