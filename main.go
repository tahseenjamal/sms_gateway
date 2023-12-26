package main

import (
	"fmt"
	"sms_gateway/broker"
)

func main() {
	// Initialize and start the consumer application

	activemq := broker.NewMessageBroker()
	fmt.Println(activemq.Connect())
	// fmt.Println(activemq.Subscribe("test_queue"))

}
