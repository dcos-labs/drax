package main

import (
	// "encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	// "io/ioutil"
	"net/http"
	"strconv"
)

// API nouns
type NOUN_Stats struct{}
type NOUN_Rampage struct{}

// Handles /stats API calls
func (n NOUN_Stats) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{"handle": "/stats"}).Info("Reporting on runtime statistics ...")
	// extract $RUNS parameter from /stats?runs=$RUNS in the following:
	runsParam := r.URL.Query().Get("runs")
	log.WithFields(log.Fields{"handle": "/stats"}).Info("Runs param = ", runsParam)
	if runs, err := strconv.Atoi(runsParam); err == nil && runs > 0 {
		log.WithFields(log.Fields{"handle": "/stats"}).Info("... for the past ", runs, " run(s)")
	} else {
		log.WithFields(log.Fields{"handle": "/stats"}).Info("... from beginning of time")
	}
	fmt.Fprint(w, "not yet implemented")
}

// Handles /rampage API calls
func (n NOUN_Rampage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		// extract $LEVEL parameter from /rampage?level=$LEVEL in the following:
		levelParam := r.URL.Query().Get("level")
		log.WithFields(log.Fields{"handle": "/rampage"}).Info("Level param = ", levelParam)
		if level, err := strconv.Atoi(levelParam); err == nil {
			destructionLevel = DestructionLevel(level)
		}
		log.WithFields(log.Fields{"handle": "/rampage"}).Info("Starting rampage on destruction level ", destructionLevel)

		switch destructionLevel {
		case DL_BASIC:
			fmt.Fprint(w, "killed some random tasks")
		case DL_ADVANCED:
			fmt.Fprint(w, "not yet implemented")
		case DL_ALL:
			fmt.Fprint(w, "not yet implemented")
		default:
			http.NotFound(w, r)
		}
	} else {
		log.WithFields(log.Fields{"handle": "/rampage"}).Error("Only POST method supported")
		http.NotFound(w, r)
	}
}
