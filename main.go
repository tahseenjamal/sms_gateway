package main

import (
	"sms_gateway/broker"
	"time"
)

func main() {
	// Initialize and start the consumer application

	activemq := broker.NewMessageBroker()
	activemq.Subscribe("test_queue")
	for {

		var msg string
		var err error

		activemq.Send("test_queue", "test message")
		msg, err = activemq.Read("test_queue")
		if err == nil {
			activemq.FileLogger.WriteLog("Received message: " + msg)
		} else {
			activemq.FileLogger.WriteLog("Error reading message: " + err.Error())
		}

		time.Sleep(1 * time.Second)
	}
}
