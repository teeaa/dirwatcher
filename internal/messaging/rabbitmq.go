package messaging

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strconv"

	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
)

/* ***
 * RabbitMQ related methods and constructs
 * ***/

type rabbitMQConfig struct {
	Username string
	Password string
	Hostname string
	Port     int

	ExchangeName string
	ExchangeType string
}

// Rabbit handles rabbit messaging datas
type rabbit struct {
	conn   *amqp.Connection
	ch     *amqp.Channel
	q      amqp.Queue
	config rabbitMQConfig
}

func (c *rabbitMQConfig) getDsn() string {
	return "amqp://" + c.Username + ":" + c.Password + "@" + c.Hostname + ":" + strconv.Itoa(c.Port)
}

var r rabbit

// Init RabbitMQ variables
func init() {
	filename := "config/rabbitmq.json"

	if _, err := os.Stat(filename); err == nil {
		configJSON, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Warning("Found RabbitMQ config but was unable to read it, using defaults")
			return
		}
		err = json.Unmarshal(configJSON, &r.config)
		if err != nil {
			log.Warning("Found RabbitMQ config but was unable to parse it, using defaults")
		}
	} else if os.IsNotExist(err) {
		r.config = rabbitMQConfig{"guest", "guest", "localhost", 25672, "dirwatch", "topic"}

		if os.Getenv("RABBITMQ_USER") != "" && os.Getenv("RABBITMQ_PASSWORD") != "" {
			r.config.Username = os.Getenv("RABBITMQ_USER")
			r.config.Password = os.Getenv("RABBITMQ_PASSWORD")
		}

		if os.Getenv("RABBITMQ_HOSTNAME") != "" {
			r.config.Hostname = os.Getenv("RABBITMQ_HOSTNAME")
		}

		if os.Getenv("RABBITMQ_PORT") != "" {
			r.config.Port, err = strconv.Atoi(os.Getenv("RABBITMQ_PORT"))

			if err != nil {
				log.Warning("Unable to convert RABBIT_PORT environment variable to an integer, using default")
			}
		}
	}
}

func listen(routingKey string) (<-chan amqp.Delivery, error) {
	err := r.ch.QueueBind(
		r.q.Name,   // queue name
		routingKey, // routing key
		r.config.ExchangeName,
		false,
		nil)

	if err != nil {
		log.Error("Error binding to message queue: ", err)
		return nil, err
	}

	msgs, err := r.ch.Consume(
		r.q.Name, // queue
		"",       // consumer
		true,     // auto-ack
		false,    // exclusive
		false,    // no-local
		false,    // no-wait
		nil,      // args
	)
	if err != nil {
		log.Error("Error starting message consumer: ", err)
		return nil, err
	}

	return msgs, nil
}

func send(routingKey string, action string, filename string) error {
	filename = stripPath(filename)
	err := r.ch.Publish(
		r.config.ExchangeName,
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(filename),
			Type:        action,
			AppId:       clientdata.get(),
		},
	)
	if err != nil {
		log.Error("Failed to send to channel: ", err)
		return err
	}

	log.Info(clientdata.get(), " > Published message ", action, ": ", filename)
	return nil
}

// Connect to RabbitMQ
func connect() error {
	log.Info("Connecting to RabbitMQ")
	var err error
	r.conn, err = amqp.Dial(r.config.getDsn())
	if err != nil {
		log.Error("Failed to connect to RabbitMQ: ", err)
		return err
	}

	r.ch, err = r.conn.Channel()
	if err != nil {
		log.Error("Failed to open channel: ", err)
		return err
	}

	err = r.ch.ExchangeDeclare(
		r.config.ExchangeName,
		r.config.ExchangeType,
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Error("Declaring exchange: ", err)
		return err
	}

	r.q, err = r.ch.QueueDeclare(
		"",    // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Error("Failed to open channel: ", err)
		return err
	}

	return nil
}

// Disconnect from RabbitMQ
func disconnect() {
	log.Info("Disconnecting from RabbitMQ")
	if r.ch != nil {
		r.ch.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
}
