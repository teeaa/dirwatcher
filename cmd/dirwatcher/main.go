package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"

	log "github.com/sirupsen/logrus"
)

func main() {
	setLogLevel()
	checkParams() // Show help if no parameters
	path := os.Args[1]
	appID, err := os.Hostname()

	watcher, err := addWatcher(appID, path)

	if err != nil {
		log.Error("Failed to initialise: ", err)
		fmt.Println("Failed to start:", err)
		return
	}

	fmt.Println("Dirwatcher started, press <Ctrl-C> to exit")

	waitForExit()

	closeWatcher(watcher)
}

func waitForExit() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}

func checkParams() {
	if len(os.Args) < 2 || strings.EqualFold(os.Args[1], "help") {
		fmt.Print("Usage: ", os.Args[0], " <directory>\n\n")
		os.Exit(0)
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
