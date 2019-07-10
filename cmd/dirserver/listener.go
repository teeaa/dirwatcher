package main

import (
	"dirwatcher/internal/messaging"
	"encoding/json"
	"strings"
)

func connect() error {
	return messaging.Connect()
}

func requestFilelists() {
	messaging.ServerSend("GET_FILELIST", "")
}

// Check if string contains substring ignoring case
func compare(str string, substr string) bool {
	return strings.Contains(strings.ToLower(str), strings.ToLower(substr))
}

// Listener for watcher datas
func listen() error {
	msgs, err := messaging.ServerListen()
	if err != nil {
		return err
	}

	go func() {
		for {
			for d := range msgs {
				// log.Infof("< Received message: %s @ %s: %s", d.Type, d.AppId, d.Body)

				switch {
				case compare(d.Type, "FULL_LIST"):
					var files []string
					json.Unmarshal(d.Body, &files)
					addFilesToDir(d.AppId, files)
				case compare(d.Type, "CREATE"):
					addFile(d.AppId, string(d.Body))
				case compare(d.Type, "REMOVE"), compare(d.Type, "RENAME"):
					// Renaming a file causes it to be removed (removing shows up as rename)
					// and created with the new filename in another event
					rmFile(d.AppId, string(d.Body))
				}
			}
		}
	}()

	return nil
}

// Disconnect from RabbitMQ
func disconnect() {
	messaging.Disconnect()
}
