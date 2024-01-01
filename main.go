package main

import (
	"fmt"
	"sms_gateway/smppconnector"
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

	for conn_instance := range conn {
		fmt.Println("SMPP connection status:", conn_instance.Status().String())
		smpp.Send("test", "123455", "test", "false")
	}

}
