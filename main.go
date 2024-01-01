package main

import (
	"fmt"
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
	// 		activemq.FileLogger.WriteLog("Received message: " + msg)
	// 	} else {
	// 		activemq.FileLogger.WriteLog("Error reading message: " + err.Error())
	// 	}

	// 	time.Sleep(1 * time.Second)
	// }

	// Initialize smpp
	smpp := smppconnector.NewSmpp()
	conn := smpp.Connect()

	time.Sleep(1 * time.Second)

	go func() {
		for conn_instance := range conn {
			fmt.Println("SMPP connection status:", conn_instance.Status().String())
		}
	}()

	for {
		err := smpp.Send("test", "919899334417", "Hello", "")
		if err != nil {
			fmt.Println("Error sending message:", err.Error())
		}
		time.Sleep(1 * time.Microsecond)
	}
}
