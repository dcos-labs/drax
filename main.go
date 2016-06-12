package main

import (
	// "encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	// "io/ioutil"
	"net/http"
	"strconv"
	// "strings"
)

const (
	VERSION   string = "0.1.0"
	DRAX_PORT int    = 7777
)

var (
	mux *http.ServeMux
)

func main() {
	mux = http.NewServeMux()
	fmt.Printf("This is DRAX in version %s listening on port %d\n", VERSION, DRAX_PORT)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{"handle": "/health"}).Info("health check")
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprint(w, "I am Groot")
	})
	p := strconv.Itoa(DRAX_PORT)
	log.Fatal(http.ListenAndServe(":"+p, mux))
}
