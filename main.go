package main

import (
	consumer "sms_gateway/logger"
)

func main() {
	// Initialize and start the consumer application

	logger := consumer.GetLumberJack()

	logger.WriteLog("Hello World!")

}
