package main

import (
	"fmt"
	"os"
	"os/signal"

	log "github.com/sirupsen/logrus"
)

func main() {
	setLogLevel()

	startupError(connect())
	startupError(listen())

	server := startRestAPI()

	fmt.Println("Dirserver started, press <Ctrl-C> to exit")
	requestFilelists()

	waitForExit()

	disconnect()
	shutdownRestAPI(server)
}

// We hook Ctrl-C for cleaner exit
func waitForExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

// Handle critical errors and quit
func startupError(err error) {
	if err != nil {
		log.Error("Failed to initialise: ", err)
		fmt.Println("Failed to start:", err)
		os.Exit(1)
	}
}

func setLogLevel() {
	switch os.Getenv("DIRWATCHER_LOGGING") {
	case "ERROR":
		log.SetLevel(log.ErrorLevel)
	case "WARN":
		log.SetLevel(log.WarnLevel)
	case "INFO":
		log.SetLevel(log.InfoLevel)
	default:
		log.SetLevel(log.ErrorLevel)
	}
}
