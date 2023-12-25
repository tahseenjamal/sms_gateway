package main

import (
	"sms_gateway/consumer"
)

func main() {
	// Initialize and start the consumer application

	logger := consumer.GetLumberJack()

	logger.WriteLog("Hello World!")

}
