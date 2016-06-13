package main

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"net/http"
	"os"
	"strconv"
)

// destruction level type (0 .. 2)
type DestructionLevel int

const (
	// DRAX version
	VERSION string = "0.2.0"
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

func init() {
	mux = http.NewServeMux()
	if dl := os.Getenv("DESTRUCTION_LEVEL"); dl != "" {
		l, _ := strconv.Atoi(dl)
		destructionLevel = DestructionLevel(l)
	}
	log.WithFields(log.Fields{"main": "init"}).Info("On destruction level ", destructionLevel)
}

func main() {
	fmt.Printf("This is DRAX in version %s listening on port %d with default destruction level %v\n", VERSION, DRAX_PORT, DL_BASIC)
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{"handle": "/health"}).Info("health check")
		fmt.Fprint(w, "I am Groot")
	})
	mux.Handle("/stats", new(NOUN_Stats))
	mux.Handle("/rampage", new(NOUN_Rampage))
	p := strconv.Itoa(DRAX_PORT)
	log.Fatal(http.ListenAndServe(":"+p, mux))
}
