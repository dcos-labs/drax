package main

import (
	// "encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	// "io/ioutil"
	"net/http"
	"os"
	"strconv"
	// "strings"
)

type DestructionLevel int

const (
	// DRAX version
	VERSION string = "0.1.0"
	// The IP port DRAX is listening on
	DRAX_PORT int = 7777
)

const (
	// DL_BASIC means destroy random tasks
	DL_BASIC DestructionLevel = iota
	// DL_ADVANCED means destroy random apps
	DL_ADVANCED
	// DL_ALL means destroy random apps and services
	DL_ALL
)

var (
	mux              *http.ServeMux
	destructionLevel DestructionLevel = DL_BASIC
)

func main() {
	fmt.Printf("This is DRAX in version %s listening on port %d with default level %v\n", VERSION, DRAX_PORT, DL_BASIC)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{"handle": "/health"}).Info("health check")
		// w.Header().Set("Access-Control-Allow-Origin", "*")
		fmt.Fprint(w, "I am Groot")
	})
	mux.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{"handle": "/stats"}).Info("Reporting on runtime statistics")
		// extract $TIME parameter from /stats?t=$TIME in the following:

		if window := r.URL.Query().Get("t"); window != "" {
			log.WithFields(log.Fields{"handle": "/stats"}).Info("For the past ", window, " second(s)")
		} else {
			log.WithFields(log.Fields{"handle": "/stats"}).Info("From beginning of time")
		}
		fmt.Fprint(w, "not yet implemented")
	})
	p := strconv.Itoa(DRAX_PORT)
	log.Fatal(http.ListenAndServe(":"+p, mux))
}

func init() {
	mux = http.NewServeMux()
	if dl := os.Getenv("DESTRUCTION_LEVEL"); dl != "" {
		l, _ := strconv.Atoi(dl)
		destructionLevel = DestructionLevel(l)
	}
	log.WithFields(log.Fields{"main": "init"}).Info("On destruction level ", destructionLevel)
}
