package main

import (
	"fmt"
	"sms_gateway/broker"
	"sms_gateway/smppconnector"
	"time"
)

func main() {
	// Initialize smpp
	smpp := smppconnector.NewSmpp().WithRateLimit(2000).Connect()

	time.Sleep(1 * time.Second)

	receiver := broker.NewMessageBroker()
	receiver.Subscribe("http_calls")
	sender := broker.NewMessageBroker()

	var msg string
	var err error

	for {

		err = sender.Send("http_calls", "test message")

		if err == nil {
			msg, err = receiver.Read("http_calls")
			receiver.FileLogger.WriteLog("|HTTP|%s", msg)
			if err == nil {
				err = smpp.Send("test", "919899334417", "Hello", "")
				if err != nil {
					fmt.Println("SMPP Error sending message: ", err.Error())
					time.Sleep(1 * time.Second)
				}
			}
		} else {
			fmt.Println("Main: Error sending message: ", err.Error())
			time.Sleep(1 * time.Second)
		}

	}

}
