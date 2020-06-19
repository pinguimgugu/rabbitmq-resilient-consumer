package queue

import (
	"fmt"
	"log"
	"time"

	"github.com/streadway/amqp"
)

type Rabbitmq struct {
	host            string
	user            string
	pwd             string
	port            int
	conn            *amqp.Connection
	connected       chan bool
	unavailableConn chan bool
}

func NewRabbitmq() *Rabbitmq {
	return &Rabbitmq{
		host:            "rabbitmq",
		user:            "guest",
		pwd:             "guest",
		port:            5672,
		connected:       make(chan bool),
		unavailableConn: make(chan bool),
	}
}

func (r *Rabbitmq) handerConnection() error {

	go func() error {
		conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/", r.user, r.pwd, r.host, r.port))

		r.conn = conn

		if err != nil {
			fmt.Println("Unable to rabbitmq connect... wait to retrieve connection")

			time.Sleep(time.Second * 2)
			return r.handerConnection()
		}

		go r.handlerErrorNotification()

		fmt.Println("> > > Connected :] ")
		r.connected <- true
		return nil
	}()

	return nil
}

func (r *Rabbitmq) handlerErrorNotification() error {

	chanError := make(chan *amqp.Error)

	r.conn.NotifyClose(chanError)

	<-chanError
	r.unavailableConn <- true
	return r.handerConnection()
}

func (r *Rabbitmq) Init() {
	r.handerConnection()
}

func (r *Rabbitmq) Consume() error {

	<-r.connected
	ch, _ := r.conn.Channel()

	defer ch.Close()

	q, _ := ch.QueueDeclare(
		"resilient_queue",
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)

	msgs, _ := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
		}
	}()

	<-r.unavailableConn
	fmt.Println("> > Lost Conection :[ ")

	return r.Consume()
}
