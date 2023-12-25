package consumer

import (
	"fmt"
)

// StartConsumer initializes and starts the consumer application.
func StartConsumer() {
	// Global variable initialization
	globalVariableInitialization()

	// Setup logging
	Logging

	// Initialize Stomp client
	stompClient := NewStompClient("localhost:61613")
	subscriber, _ := stompClient.Subscribe(consumer_queue)
	defer stompClient.Unsubscribe(subscriber)

	// Initialize SMPP client
	smppConfig := SMPPConfig{
		// Configure SMPP client parameters
		// ...
	}
	smppClient := NewSMPPClient(smppConfig)

	// Start a goroutine for sending messages
	go sendMessage()

	// Initialize and bind to SMPP
	smppClient.smppBinding()

	// Start goroutines for SMPP delivery
	for i := 0; i < threads; i++ {
		go smppClient.smppDelivery()
	}

	// Start a goroutine for logging
	go Producer()

	// Start a goroutine for handling black hours
	go blackHourNotification()

	fmt.Printf("Starting %d thread(s)\n", threads)

	// Block the main goroutine to keep the application running
	select {}
}
