package main

import (
	"fmt"
	"net/url"
	"sms_gateway/broker"
	"sms_gateway/handler"
	"sms_gateway/smppconnector"
	"strings"
	"time"
)

func main() {
	// Initialize smpp
	smpp := smppconnector.NewSmpp().WithRateLimit(2000).Connect()
	time.Sleep(1 * time.Second)

	receiver := broker.NewMessageBroker()
	receiver.Subscribe("http_calls")

	var err error
	var msg string

	go func() {

		for {

			msg, err = receiver.Read("http_calls")
			if err == nil {

				q, err := url.ParseQuery(msg)
				if err == nil {

					from := q.Get("from")
					to := strings.Trim(q.Get("to"), " ")
					message := q.Get("message")
					test := q.Get("test")

					err = smpp.Send(from, to, message, test)
					if err != nil {
						fmt.Println("SMPP Error sending message: ", err.Error())
						time.Sleep(1 * time.Second)
					}
				}
			}
		}
	}()

	requestHandler := handler.NewRequestHandler()
	requestHandler.Listen()

}
