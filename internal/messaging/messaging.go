package messaging

import (
	"errors"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

/* ***
 * Abstraction layer for sending events.
 * ***/

// FromServer RabbitMQ routing key to use for server's messages
const FromServer = "from_server"

// FromClient RabbitMQ routing key to use for client's messages
const FromClient = "from_client"

// Connect to RabbitMQ
func Connect() error {
	return connect()
}

// Disconnect from RabbitMQ
func Disconnect() {
	disconnect()
}

// ClientListen message listener for watchers
func ClientListen() (<-chan amqp.Delivery, error) {
	return listen(FromServer)
}

// ServerListen message listener for server
func ServerListen() (<-chan amqp.Delivery, error) {
	return listen(FromClient)
}

// ClientSend wrapper for sending messages
func ClientSend(action string, message string) error {
	return send(FromClient, action, message)
}

// ServerSend wrapper for sending messages
func ServerSend(action string, message string) error {
	return send(FromServer, action, message)
}

// RegisterClient register client topic/subject details
func RegisterClient(appID string, path string) error {
	clientdata = ClientData{appID, fixPath(path)}

	err := Connect()
	if err != nil {
		log.Error("Error registering client at connect: ", err)
		return errors.New("Unable to connect to RabbitMQ")
	}

	return nil
}

// UnregisterClient unregister client
func UnregisterClient() {
	disconnect()
}
