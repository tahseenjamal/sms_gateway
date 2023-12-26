package main

import (
	"fmt"
	"sms_gateway/broker"
)

func main() {
	// Initialize and start the consumer application

	activemq := broker.NewMessageBroker()
	activemq.Send("test_queue", "test message")
	activemq.Subscribe("test_queue")
	fmt.Println(activemq.Read("test_queue"))

	c := make(chan bool)
	<-c
}
