FROM golang:1.11

ENV RABBITMQ_HOSTNAME=172.12.1.2
ENV RABBITMQ_PORT=5672

ENV DIRWATCHER_LOGGING=INFO

WORKDIR /go/src/dirwatcher

COPY . .

RUN go get github.com/sirupsen/logrus
RUN go get github.com/streadway/amqp
RUN go get github.com/gorilla/mux
RUN go get github.com/fsnotify/fsnotify

RUN go install ./...

EXPOSE 8080

CMD ["dirserver"]
