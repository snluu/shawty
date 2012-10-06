package main

import "fmt"
import "code.google.com/p/gorilla/mux"
import "net"
import "os"
import "net/http"

func main() {

	var router = mux.NewRouter()
	router.HandleFunc("/", HandleIndex)
	router.HandleFunc("/shawty.js", HandleShawtyJS)
	router.HandleFunc("/{shortID:[A-Za-z0-9]+}", HandleShortID)

	var port = os.Getenv("PORT")
	if port == "" {
		port = "80"
	}

	var l, err = net.Listen("tcp", "0.0.0.0:"+port)
	if err != nil {
		Lerrorf("Cannot listen at %s", port)
		fmt.Println(err)
		return
	}
	defer l.Close()
	Linfof("Listening at %s", port)

	http.Serve(l, router)
}
