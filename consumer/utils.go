package consumer

import (
	"regexp"
)

// Packet represents a data packet.
type Packet struct {
	Queue   string
	Message string
}

var (
	// Other global variables go here
	// ...

	re *regexp.Regexp
)

func extract(message string) (map[string]string, error) {
	// Implementation remains the same
	// ...
}

func splitString(input string, delimiter string) (int, int) {
	// Implementation remains the same
	// ...
}

func globalVariableInitialization() {
	// Implementation remains the same
	// ...
}

func blackHour() bool {
	// Implementation remains the same
	// ...
}

// Other utility functions go here
// ...
