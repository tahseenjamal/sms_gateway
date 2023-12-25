package main

import (
	"sms_gateway/logger"
)

func main() {
	// Initialize and start the consumer application

	logging := logger.GetLumberJack()

	logging.WriteLog("Hello World!")

}
