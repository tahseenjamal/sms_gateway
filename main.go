package main

import "sms_gateway/logger"

func main() {
	// Initialize and start the consumer application

	logger := logger.GetLumberJack()

	logger.WriteLog("Hello World!")

}
