## Dirwatcher

### This project consists of two executables:
* **dirserver** which servers as a data store and offers a REST endpoint at GET /files
* **dirwatcher** which watches a directory and sends notifications to the server via RabbitMQ

### Installing

#### With Docker

Go to the project directory and run `docker-compose up -d --build` to create containers for **dirserver** in port 8080 and **RabbitMQ** in port 25672.

Build the **dirwatcher** by running `make` in the project directory. This will create **dirwatcher** binary in to the project root (along with **dirserver** which we can ignore for now)

Initiate the watcher by running `./dirwatcher <directory>`, eq `./dirwatcher .`

Access http://localhost:8080/files to get a JSON list of files

### Partially with Docker

Go to the project directory and run `docker-compose up -d --build rabbitmq` - this will create just the RabbitMQ container to port 25672.

Build the executables by running `make` in the project directory

Create a copy of config/rabbitmq.json.example as config/rabbitmq.json and modify it to use your own RabbitMQ.

#### Without Docker (Not recommended)

You need to have RabbitMQ running in an accessible place.

Create a copy of config/rabbitmq.json.example as config/rabbitmq.json and modify it to use your own RabbitMQ.

Build the executables by running `make` in the project directory

Run `./dirserver` in project root so it has access to config file or set up environment variables (documented below)

Run `./dirwatcher` as documented above

### Environment variables

#### Dirwatcher / dirserver log levels

By default the logging is set to ERROR. If you want to set it to WARN or INFO instead, use **DIRWATCHER_LOGGING** environment variable. Eg `export DIRWATCHER_LOGGING=INFO` (in Dockerised dirserver the log level is automatically set to INFO)

#### RabbitMQ connection variables

- `RABBITMQ_USER` for username
- `RABBITMQ_PASSWORD` for password
- `RABBITMQ_HOSTNAME` for hostname/ip
- `RABBITMQ_PORT` for port

The defaults for those are `guest` for username and password, `localhost` for hostname and `25672` for port.
