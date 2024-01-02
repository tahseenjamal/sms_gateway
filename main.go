package main

import (
	"sms_gateway/smppconnector"
	"time"
)

func main() {
	// Initialize and start the consumer application

	// activemq := broker.NewMessageBroker()
	// activemq.Subscribe("test_queue")
	// for {

	// 	var msg string
	// 	var err error

	// 	activemq.Send("test_queue", "test message")
	// 	msg, err = activemq.Read("test_queue")
	// 	if err == nil {
	// 		activemq.FileLogger.WriteLog("Received message: %s ", msg)
	// 	}

	// 	time.Sleep(1 * time.Second)
	// }

	// Initialize smpp
	smpp := smppconnector.NewSmpp()
	smpp.Connect()

	time.Sleep(1 * time.Second)

	for {
		_ = smpp.Send("test", "919899334417", "Hello", "")
	}
}
