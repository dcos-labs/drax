package main

import (
	// "encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	// "io/ioutil"
	marathon "github.com/gambol99/go-marathon"
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
			killTasks(w, r)
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

func killTasks(w http.ResponseWriter, r *http.Request) {
	marathonURL := "http://localhost:8080"
	config := marathon.NewDefaultConfig()
	config.URL = marathonURL
	client, err := marathon.NewClient(config)
	if err != nil {
		log.WithFields(log.Fields{"handle": "/rampage"}).Error("Failed to create Marathon client due to ", err)
		http.NotFound(w, r)
		return
	}
	applications, err := client.Applications(nil)
	if err != nil {
		log.Fatalf("Failed to list apps")
		log.WithFields(log.Fields{"handle": "/rampage"}).Info("Failed to list apps")
		fmt.Fprint(w, "Failed to list apps")
		return
	}
	log.WithFields(log.Fields{"handle": "/rampage"}).Info("Found ", len(applications.Apps), " applications running")
	b := ""
	for _, application := range applications.Apps {
		log.WithFields(log.Fields{"handle": "/rampage"}).Debug("APP ", application.ID)
		details, _ := client.Application(application.ID)
		if details.Tasks != nil && len(details.Tasks) > 0 {
			health, _ := client.ApplicationOK(details.ID)
			b += fmt.Sprintf("Application: %s is healthy: %t\n", details.ID, health)
			for _, task := range details.Tasks {
				log.WithFields(log.Fields{"handle": "/rampage"}).Debug("TASK ", task.ID)
				b += fmt.Sprintf(" Task: %s\n", task.ID)
			}
		}
	}
	fmt.Fprint(w, b)
}
