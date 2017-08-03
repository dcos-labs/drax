package main

import (
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync/atomic"
	"time"

	log "github.com/Sirupsen/logrus"
	marathon "github.com/gambol99/go-marathon"
)

// StatsResult JSON payload
type StatsResult struct {
	KilledTasks uint64 `json:"gone"`
}

// RampageResult JSON payload
type RampageResult struct {
	Success     bool     `json:"success"`
	KilledTasks []string `json:"goners"`
}

// Handles /health API calls (GET only)
func getHealth(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{"handle": "/health"}).Info("health check")
	io.WriteString(w, "I am Groot")
}

// Handles /stats API calls (GET only)
func getStats(w http.ResponseWriter, r *http.Request) {
	log.WithFields(log.Fields{"handle": "/stats"}).Info("Overall tasks killed: ", overallTasksKilled)

	sr := &StatsResult{
		KilledTasks: overallTasksKilled}

	srJSON, _ := json.Marshal(sr)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(srJSON))
}

// Handles /rampage API calls (POST only)
func postRampage(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.WithFields(log.Fields{"handle": "/rampage"}).Info("Starting rampage on destruction level ", destructionLevel)

		switch destructionLevel {
		case DLBASIC:
			killTasks(w, r)
		case DLADVANCED:
			io.WriteString(w, "not yet implemented")
		case DLALL:
			io.WriteString(w, "not yet implemented")
		default:
			http.NotFound(w, r)
		}

	} else {
		log.WithFields(log.Fields{"handle": "/rampage"}).Error("Only POST method supported")
		http.NotFound(w, r)
	}

}

// killTasks will identify tasks of any apps (but not framework services)
// and randomly kill off a few of them
func killTasks(w http.ResponseWriter, r *http.Request) {
	if client, ok := getClient(); ok {
		apps, err := client.Applications(nil)
		if err != nil {
			log.WithFields(log.Fields{"handle": "/rampage"}).Info("Failed to list apps")
			http.Error(w, "Failed to list apps", 500)
			return
		}
		log.WithFields(log.Fields{"handle": "/rampage"}).Info("Found overall ", len(apps.Apps), " applications running")

		candidates := []string{}
		for _, app := range apps.Apps {
			log.WithFields(log.Fields{"handle": "/rampage"}).Debug("APP ", app.ID)
			details, _ := client.Application(app.ID)

			if !myself(details) && !isFramework(details) {
				if details.Tasks != nil && len(details.Tasks) > 0 {
					for _, task := range details.Tasks {
						log.WithFields(log.Fields{"handle": "/rampage"}).Debug("TASK ", task.ID)
						candidates = append(candidates, task.ID)
					}
				}
			}
		}

		if len(candidates) > 0 {
			log.WithFields(log.Fields{"handle": "/rampage"}).Info("Found ", len(candidates), " tasks to kill")
			rampage(w, client, candidates)
		}

	} else {
		http.Error(w, "Can't connect to Marathon", 500)
	}
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

// rampage kills random tasks from the candidates and returns a JSON result
func rampage(w http.ResponseWriter, c marathon.Marathon, candidates []string) {
	var targets []int

	// generates a list of random, non-repeating indices into the candidates:
	if len(candidates) > numTargets {
		targets = rand.Perm(len(candidates))[:numTargets]
	} else {
		targets = rand.Perm(len(candidates))
	}

	rr := &RampageResult{
		Success:     true,
		KilledTasks: []string{}}

	// kill the candidates
	for _, t := range targets {
		candidate := candidates[t]
		killTask(c, candidate)
		log.WithFields(log.Fields{"handle": "/rampage"}).Info("Killed tasks ", candidate)
		rr.KilledTasks = append(rr.KilledTasks, candidate)
		atomic.AddUint64(&overallTasksKilled, 1)
		log.WithFields(log.Fields{"handle": "/rampage"}).Info("Counter: ", overallTasksKilled)
		time.Sleep(time.Millisecond * time.Duration(sleepTime))
	}

	rrJSON, _ := json.Marshal(rr)
	w.Header().Set("Content-Type", "application/json")
	io.WriteString(w, string(rrJSON))
}

// killTask kills a certain task and increments overall count if successful
func killTask(c marathon.Marathon, taskID string) bool {
	_, err := c.KillTask(taskID, nil)
	if err != nil {
		log.WithFields(log.Fields{"handle": "/rampage"}).Debug("Not able to kill task ", taskID)
		return false
	}

	log.WithFields(log.Fields{"handle": "/rampage"}).Debug("Killed task ", taskID)
	return true
}

// myself returns true if it is applied to DRAX Marathon app itself
func myself(app *marathon.Application) bool {
	if strings.Contains(app.ID, "drax") {
		return true
	}
	return false
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
