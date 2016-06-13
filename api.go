package main

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	// "io/ioutil"
	marathon "github.com/gambol99/go-marathon"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
)

// API nouns
type NOUN_Stats struct{}
type NOUN_Rampage struct{}

// JSON payloads
type RampageParams struct {
	Level string `json:"level"`
	AppID string `json:"app"`
}
type StatsResult struct {
	TasksKilled uint64 `json:"gone"`
}
type RampageResult struct {
	Success  bool     `json:"success"`
	TaskList []string `json:"goners"`
}

// Handles /stats API calls (GET only)
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
	sr := &StatsResult{}
	sr.TasksKilled = atomic.LoadUint64(&overallTasksKilled)
	jsonsr, _ := json.Marshal(sr)
	w.Header().Set("Content-Type", "application/javascript")
	fmt.Fprint(w, string(jsonsr))
}

// Handles /rampage API calls (POST only)
func (n NOUN_Rampage) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Can't parse rampage params", 500)
		} else {
			ok, rp := parseRampageParams(r)
			if !ok {
				http.Error(w, "Can't decode rampage params", 500)
			}
			levelParam := rp.Level
			if levelParam != "" {
				log.WithFields(log.Fields{"handle": "/rampage"}).Info("Got level param ", levelParam)
				if level, err := strconv.Atoi(levelParam); err == nil {
					destructionLevel = DestructionLevel(level)
				}
			}
			log.WithFields(log.Fields{"handle": "/rampage"}).Info("Starting rampage on destruction level ", destructionLevel)
			switch destructionLevel {
			case DL_BASIC:
				killTasks(w, r)
			case DL_ADVANCED:
				appParam := rp.AppID
				if appParam != "" {
					log.WithFields(log.Fields{"handle": "/rampage"}).Info("Got app param ", appParam)
					killTasksOfApp(w, r, appParam)
				} else {
					http.NotFound(w, r)
				}
			case DL_ALL:
				fmt.Fprint(w, "not yet implemented")
			default:
				http.NotFound(w, r)
			}
		}
	} else {
		log.WithFields(log.Fields{"handle": "/rampage"}).Error("Only POST method supported")
		http.NotFound(w, r)
	}
}

// parseRampageParams parses the parameters for a rampage from an HTTP request
func parseRampageParams(r *http.Request) (bool, *RampageParams) {
	decoder := json.NewDecoder(r.Body)
	rp := &RampageParams{}
	err := decoder.Decode(rp)
	if err != nil {
		return false, nil
	} else {
		return true, rp
	}
}

// killTasks will identify tasks of any apps (but not framework services)
// and randomly kill off a few of them
func killTasks(w http.ResponseWriter, r *http.Request) {
	if client, ok := getClient(); ok {
		nonFrameworkApps := 0
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
				nonFrameworkApps++
				if details.Tasks != nil && len(details.Tasks) > 0 {
					for _, task := range details.Tasks {
						log.WithFields(log.Fields{"handle": "/rampage"}).Debug("TASK ", task.ID)
						candidates = append(candidates, task.ID)
					}
				}
			}
		}
		rampage(w, client, nonFrameworkApps, candidates)
	} else {
		http.Error(w, "Can't connect to Marathon", 500)
	}
}

// killTasks will identify tasks of a specific app defined by targetAppID
// and randomly kill off a few of them
func killTasksOfApp(w http.ResponseWriter, r *http.Request, targetAppID string) {
	if client, ok := getClient(); ok {
		candidates := []string{}
		details, _ := client.Application(targetAppID)
		log.WithFields(log.Fields{"handle": "/rampage"}).Info("Found app ", details.ID, " running")
		if !myself(details) && !isFramework(details) {
			if details.Tasks != nil && len(details.Tasks) > 0 {
				for _, task := range details.Tasks {
					log.WithFields(log.Fields{"handle": "/rampage"}).Debug("TASK ", task.ID)
					candidates = append(candidates, task.ID)
				}
			}
		}
		rampage(w, client, 1, candidates)
	} else {
		http.Error(w, "Can't connect to Marathon", 500)
	}
}

// rampage kills random tasks from the candidates and returns a JSON result
func rampage(w http.ResponseWriter, c marathon.Marathon, numApps int, candidates []string) {
	rr := &RampageResult{}
	rr.TaskList = []string{}
	targets := []int{}
	if len(candidates) > 0 {
		log.WithFields(log.Fields{"handle": "/rampage"}).Info("Found ", len(candidates), " tasks in ", numApps, " apps to kill")
		// generates a list of random, non-repeating indices into the candidates:
		if len(candidates) > numTargets {
			targets = rand.Perm(len(candidates))[:numTargets]
		} else {
			targets = rand.Perm(len(candidates))
		}
		for _, t := range targets {
			candidate := candidates[t]
			rr.Success = killTask(c, candidate)
			if rr.Success {
				rr.TaskList = append(rr.TaskList, candidate)
			}
		}
		log.WithFields(log.Fields{"handle": "/rampage"}).Info("Killed tasks ", rr.TaskList)
		// at least killed some tasks, so consider it a success:
		rr.Success = true
	} else {
		rr.Success = false
	}
	jsonrr, _ := json.Marshal(rr)
	w.Header().Set("Content-Type", "application/javascript")
	fmt.Fprint(w, string(jsonrr))
}

// killTask kills a certain task and increments overall count if successful
func killTask(c marathon.Marathon, taskID string) bool {
	_, err := c.KillTask(taskID, nil)
	if err != nil {
		log.WithFields(log.Fields{"handle": "/rampage"}).Debug("Not able to kill task ", taskID)
		return false
	} else {
		log.WithFields(log.Fields{"handle": "/rampage"}).Debug("Killed task ", taskID)
		go incTasksKilled()
		return true
	}
}

// myself returns true if it is applied to DRAX Marathon app itself
func myself(app *marathon.Application) bool {
	if strings.Contains(app.ID, "drax") {
		return true
	} else {
		return false
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

// incTasksKilled increases the overall tasks killed counter in an atomic way
func incTasksKilled() {
	atomic.AddUint64(&overallTasksKilled, 1)
}
