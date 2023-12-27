package main

import (
	"fmt"
	"sms_gateway/broker"
	"time"
)

func main() {
	// Initialize and start the consumer application

	activemq := broker.NewMessageBroker()
	activemq.Subscribe("test_queue")
	for {

		activemq.Send("test_queue", "test message")
		fmt.Println(activemq.Read("test_queue"))

		time.Sleep(1 * time.Second)
	}
}
