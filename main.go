package main

import (
	"rabbitmq-consumer/queue"
)

func main() {

	r := queue.NewRabbitmq()
	r.Init()
	r.Consume()

}
