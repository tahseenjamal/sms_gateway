package consumer

import (
	"github.com/go-stomp/stomp"
)

// StompClient represents a Stomp client.
type StompClient struct {
	conn *stomp.Conn
}

// NewStompClient creates a new StompClient.
func NewStompClient(ipPort string) *StompClient {
	return &StompClient{
		conn: getStompSession(ipPort),
	}
}

func (s *StompClient) Subscribe(queueName string) (*stomp.Subscription, error) {
	return stomp_subscribe(s.conn, queueName)
}

func (s *StompClient) Unsubscribe(sub *stomp.Subscription) {
	sub.Unsubscribe()
	s.conn.Disconnect()
}

func getStompSession(ipPort string) *stomp.Conn {
	// Implementation remains the same
	// ...
}

func stomp_subscribe(conn *stomp.Conn, queueName string) (*stomp.Subscription, error) {
	// Implementation remains the same
	// ...
}
