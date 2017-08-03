package main

import (
	"net/http"
	"os"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
)

// DestructionLevel of type (0 .. 2)
type DestructionLevel int

const (
	// VERSION of DRAX
	VERSION string = "0.4.0"
	// DEFAULTPORT is the port DRAX is listening on
	DEFAULTPORT string = "7777"
	// MARATHONURL for connection to DC/OS
	MARATHONURL string = "http://marathon.mesos:8080"
	// DEFAULTNUMTARGET is the number of tasks to kill
	DEFAULTNUMTARGET int = 2
	// DEFAULTSLEEPTIME is the time in ms to wait between the killing of tasks
	DEFAULTSLEEPTIME int = 100
)

const (
	// DLBASIC means destroy random tasks
	DLBASIC DestructionLevel = iota
	// DLADVANCED means destroy random apps
	DLADVANCED
	// DLALL means destroy random apps and services
	DLALL
)

var (
	port               string
	marathonURL        string
	destructionLevel   = DestructionLevel(DLBASIC)
	numTargets         = int(DEFAULTNUMTARGET)
	sleepTime          = int(DEFAULTSLEEPTIME)
	overallTasksKilled uint64
)

func init() {

	if ll := os.Getenv("LOG_LEVEL"); ll != "" {
		switch strings.ToUpper(ll) {
		case "DEBUG":
			log.SetLevel(log.DebugLevel)
		case "INFO":
			log.SetLevel(log.InfoLevel)
		default:
			log.SetLevel(log.ErrorLevel)
		}

	}

	log.WithFields(log.Fields{"main": "init"}).Info("This is DRAX in version ", VERSION)

	// set port for http server
	if port = os.Getenv("PORT"); len(port) == 0 {
		port = DEFAULTPORT
	}
	log.WithFields(log.Fields{"main": "init"}).Info("Listening on port ", port)

	// set destruction level
	if dl := os.Getenv("DESTRUCTION_LEVEL"); dl != "" {
		l, _ := strconv.Atoi(dl)
		destructionLevel = DestructionLevel(l)
	}
	log.WithFields(log.Fields{"main": "init"}).Info("On destruction level ", destructionLevel)

	// per default, use the cluster-internal, non-auth endpoint:
	if marathonURL = os.Getenv("MARATHON_URL"); marathonURL == "" {
		marathonURL = MARATHONURL
	}
	log.WithFields(log.Fields{"main": "init"}).Info("Using Marathon at  ", marathonURL)

	if nt := os.Getenv("NUM_TARGETS"); nt != "" {
		n, _ := strconv.Atoi(nt)
		numTargets = n
	}
	log.WithFields(log.Fields{"main": "init"}).Info("I will destroy ", numTargets, " tasks on a rampage")

	if st := os.Getenv("SLEEP_TIME"); st != "" {
		s, _ := strconv.Atoi(st)
		sleepTime = s
	}
	log.WithFields(log.Fields{"main": "init"}).Info("I will wait ", sleepTime, "ms between the killing of tasks")

}

func main() {

	http.HandleFunc("/health", getHealth)
	http.HandleFunc("/stats", getStats)
	http.HandleFunc("/rampage", postRampage)

	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

}
