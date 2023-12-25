package main

import (
	"consumer/logger"
)

func main() {
	// Initialize and start the consumer application

	logger := logger.GetLumberJack()

	logger.ConsumerLogger.Println("Hello World!")

	logger.writeLog("Hello World!")

}
