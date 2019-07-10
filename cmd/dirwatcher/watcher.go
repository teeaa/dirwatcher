package main

import (
	"dirwatcher/internal/messaging"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
)

// fs watching related methods

// Check that given path exists, is a directory and return it in absolute form
func checkPath(path string) (string, error) {
	path, err := filepath.Abs(path)
	if err != nil {
		log.Error("Failed to resolve absolute path: ", err)
		return "", err
	}
	fi, err := os.Stat(path)
	if err != nil {
		log.Error("Error checking directory: ", err)
		return "", err
	}
	if !fi.IsDir() {
		log.Error("Path is not a directory")
		return "", errors.New("Not a directory: " + path)
	}
	return path, nil
}

// Send message to clear lists (when exiting)
func clearFilelist() error {
	return messaging.ClientSend("FULL_LIST", "[]")
}

// SendFilelist send full list of files
func sendFilelist(filelist []os.FileInfo) error {
	var filelistArr []string

	// Remove directories from list and make it a slice of filenames
	for _, file := range filelist {
		if !file.IsDir() {
			filelistArr = append(filelistArr, file.Name())
		}
	}

	filelistJSON, err := json.Marshal(filelistArr)
	if err != nil {
		log.Error("Error marshalling filelist to JSON: ", err)
		return err
	}

	return messaging.ClientSend("FULL_LIST", string(filelistJSON))
}

// Creates a watcher for directory
func addWatcher(appID string, path string) (*fsnotify.Watcher, error) {
	path, err := checkPath(path)
	if err != nil {
		return nil, err
	}

	filelist, err := getFilelist(path)
	if err != nil {
		return nil, err
	}

	err = messaging.RegisterClient(appID, path)
	if err != nil {
		return nil, err
	}

	err = sendFilelist(filelist)
	if err != nil {
		return nil, err
	}

	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		log.Error("Failed to initalise watcher: ", err)
		return nil, err
	}

	go dirWatcher(watcher)

	err = watcher.Add(path)

	if err != nil {
		log.Error("Error adding watcher for path ", path, ": ", err)
		return nil, err
	}

	log.Info("Started watching ", path, " for changes")
	listen(path)
	return watcher, err
}

// Run cleanup
func closeWatcher(watcher *fsnotify.Watcher) {
	err := clearFilelist()
	if err != nil {
		log.Error("Error clearing filelist: ", err)
	}
	watcher.Close()
	messaging.UnregisterClient()
}

// Get full filelist for path
func getFilelist(path string) ([]os.FileInfo, error) {
	f, err := os.Open(path)
	defer f.Close()
	if err != nil {
		log.Error("Error opening directory: ", err)
		return nil, err
	}

	filelist, err := f.Readdir(-1)
	if err != nil {
		log.Error("Error reading directory: ", err)
		return nil, err
	}

	return filelist, nil
}

// Send notifications on events in watched directory to server
func dirWatcher(watcher *fsnotify.Watcher) {
	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			log.Debug("Event: ", event)
			if arrContains(event.Op.String()) && !creatingDirectory(event.Op.String(), event.Name) {
				err := messaging.ClientSend(event.Op.String(), event.Name)
				if err != nil {
					log.Error("Error sending event to server: ", err)
					return
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Error("Error: ", err)
		}
	}
}

// Are we creating a directory?
func creatingDirectory(action string, file string) bool {
	// Creation event? Otherwise we can't stat file is it was removed
	if !strings.Contains(action, "CREATE") {
		return false
	}
	fi, err := os.Stat(file)
	if err != nil {
		log.Error("Error statting file "+file+": ", err)
	}

	return fi.IsDir()
}

// Check if action contains an event we want to broadcast
func arrContains(str string) bool {
	for _, compare := range []string{"CREATE", "REMOVE", "RENAME"} {
		if strings.Contains(strings.ToUpper(str), compare) {
			return true
		}
	}
	return false
}
