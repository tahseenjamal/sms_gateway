package main

import (
	"fmt"
	"sms_gateway/logger"
)

func main() {
	// Initialize and start the consumer application

	logging := logger.GetLumberJack()

	fmt.Println(logger.GetLumberJack())
	fmt.Println(logger.GetLumberJack())

	logging.WriteLog("Hello World!")

}
