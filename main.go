package main

import (
	"sms_gateway/consumer"
)

func main() {
	// Initialize and start the consumer application

	logger := consumer.GetLumberJack()

	logger.ConsumerLogger.Println("Hello World!")

	logger.writeLog("Hello World!")

}
