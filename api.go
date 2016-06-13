package main

import (
	// "encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	// "io/ioutil"
	marathon "github.com/gambol99/go-marathon"
	"math/rand"
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

// killTasks will identify tasks from apps (not framework service)
// to be killed and randomly kill off a few of them
func killTasks(w http.ResponseWriter, r *http.Request) {
	if client, ok := getClient(); ok {
		apps, err := client.Applications(nil)
		if err != nil {
			log.WithFields(log.Fields{"handle": "/rampage"}).Info("Failed to list apps")
			http.Error(w, "Failed to list apps", 500)
			return
		}
		log.WithFields(log.Fields{"handle": "/rampage"}).Info("Found ", len(apps.Apps), " applications running")
		b := ""
		candidates := []string{}
		for _, app := range apps.Apps {
			log.WithFields(log.Fields{"handle": "/rampage"}).Debug("APP ", app.ID)
			details, _ := client.Application(app.ID)
			if !isFramework(details) {
				if details.Tasks != nil && len(details.Tasks) > 0 {
					for _, task := range details.Tasks {
						log.WithFields(log.Fields{"handle": "/rampage"}).Debug("TASK ", task.ID)
						candidates = append(candidates, task.ID)
					}
				}
			}
		}

		if len(candidates) > 0 {
			// pick one random task to be killed
			candidate := candidates[rand.Intn(len(candidates))]
			ok := killTask(client, candidate)
			if ok {
				b += fmt.Sprintf("Killed task %s", candidate)
			} else {
				b += fmt.Sprintf("Failed to kill task %s", candidate)
			}
		} else {
			b = fmt.Sprintf("No task found to kill")
		}
		fmt.Fprint(w, b)
	} else {
		http.Error(w, "Can't connect to Marathon", 500)
	}
}

// killTask kills a certain task
func killTask(c marathon.Marathon, taskID string) bool {
	_, err := c.KillTask(taskID, nil)
	if err != nil {
		log.WithFields(log.Fields{"handle": "/rampage"}).Debug("Not able to kill task ", taskID)
		return false
	} else {
		log.WithFields(log.Fields{"handle": "/rampage"}).Debug("Killed task ", taskID)
		return true
	}
}

// isFramework returns true if the Marathon app is a service framework,
// and false otherwise (determined via the DCOS_PACKAGE_IS_FRAMEWORK label key)
func isFramework(app *marathon.Application) bool {
	for k, v := range *app.Labels {
		log.WithFields(log.Fields{"handle": "/rampage"}).Debug("LABEL ", k, ":", v)
		if k == "DCOS_PACKAGE_IS_FRAMEWORK" && v == "true" {
			return true
		}
	}
	return false
}

// getClient tries to get a connection to the DC/OS System Marathon
func getClient() (marathon.Marathon, bool) {
	config := marathon.NewDefaultConfig()
	config.URL = marathonURL
	client, err := marathon.NewClient(config)
	if err != nil {
		log.WithFields(log.Fields{"handle": "/rampage"}).Error("Failed to create Marathon client due to ", err)
		return nil, false
	}
	return client, true
}
