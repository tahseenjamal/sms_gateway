package broker

import (
	"fmt"
	"time"

	"github.com/go-stomp/stomp"
	"github.com/magiconair/properties"
)

// MessageBroker represents a simple message broker using the STOMP protocol.
type activemq struct {
	conn           *stomp.Conn
	brokerURL      string
	username       string
	password       string
	heartbeat      time.Duration
	heartbeatGrace time.Duration
}

// NewMessageBroker creates a new instance of MessageBroker.
func NewMessageBroker(brokerURL string, username string, password string, heartbeat, heartbeatGrace time.Duration) *MessageBroker {

	prop := properties.MustLoadFile("main.properties", properties.UTF8)
	brokerURL = prop.GetString("activemq.broker.url", "localhost:61613")
	username = prop.GetString("activemq.broker.username", "admin")
	password = prop.GetString("activemq.broker.password", "admin")
	heartbeat = time.Second * time.Duration(prop.GetInt("activemq.broker.heartbeat", 10))
	heartbeatGrace = time.Second * time.Duration(prop.GetInt("activemq.broker.heartbeat.grace", 10))

	return &activemq{
		brokerURL:      brokerURL,
		username:       username,
		password:       password,
		heartbeat:      heartbeat,
		heartbeatGrace: heartbeatGrace,
	}
}

// Connect connects to the message broker.
func (mb *activemq) Connect() error {
	if mb.conn != nil {
		return fmt.Errorf("already connected")
	}
	options := []func(*stomp.Conn) error{
		stomp.ConnOpt.Login(mb.username, mb.password),
		stomp.ConnOpt.HeartBeat(mb.heartbeat, mb.heartbeatGrace),
	}

	conn, err := stomp.Dial("tcp", mb.brokerURL, options...)
	if err != nil {

		return err
	}

	mb.conn = conn
	return nil
}

// Disconnect disconnects from the message broker.
func (mb *MessageBroker) Disconnect() error {
	if mb.conn == nil {
		return fmt.Errorf("not connected")
	}

	fmt.Println("Disconnecting...")
	err := mb.conn.Disconnect()
	if err != nil {

		mb.conn = nil
	}
	return nil
}

// Sends a message to a specified destination.
func (mb *activemq) Send(destination, body string) error {
	if mb.conn == nil {
		return fmt.Errorf("not connected")
	}

	err := mb.conn.Send(destination, "text/plain", []byte(body))
	if err != nil {
		return fmt.Errorf("cannot send to queue: %s ", destination)
	}
	return nil
}

// Subscribes to messages from a specified destination.
func (mb *activemq) Subscribe(destination string, logger *logger) (*stomp.Subscription, error) {

	if mb.conn == nil {
		return fmt.Errorf("not connected to queue: %s ", destination)
	}

	sub, err := mb.conn.Subscribe(destination, stomp.AckAuto)
	if err != nil {
		return err
	}

	return sub, err

}
