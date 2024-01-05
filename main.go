package main

import (
	"fmt"
	"net/url"
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

	var err error
	var msg string

	for {

		data := "from=XtraCash&to=%2B254722000000&message=Hello%20World&test=false"

		err = sender.Send("http_calls", data)

		if err == nil {
			msg, err = receiver.Read("http_calls")
			receiver.FileLogger.WriteLog("|HTTP|%s", msg)
			if err == nil {

				// Parse query string using net/url package
				// and convert to map[string]string

				q, err := url.ParseQuery(msg)
				from := q.Get("from")
				to := q.Get("to")
				message, _ := url.QueryUnescape(q.Get("message"))
				test := q.Get("test")

				err = smpp.Send(from, to, message, test)
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
