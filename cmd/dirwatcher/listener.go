package main

import (
	"dirwatcher/internal/messaging"

	log "github.com/sirupsen/logrus"
)

// Listen for events from the server, mostly GET_FILELIST for when server restarts
func listen(path string) error {
	msgs, err := messaging.ClientListen()
	if err != nil {
		return err
	}

	go func() {
		for {
			for d := range msgs {
				// log.Printf("< Received message: %s", d.Type)

				if d.Type == "GET_FILELIST" {
					filelist, err := getFilelist(path)
					if err != nil {
						log.Error("Error getting list of files: ", err)
					}
					if len(filelist) > 0 {
						sendFilelist(filelist)
					}
				}
			}
		}
	}()

	return nil
}
