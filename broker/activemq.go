package broker

import (
	"fmt"
	"time"

	logger "sms_gateway/logger"

	"github.com/go-stomp/stomp"
	"github.com/magiconair/properties"
)

type activemqConfig struct {
	brokerURL      string
	username       string
	password       string
	heartbeat      time.Duration
	heartbeatGrace time.Duration
}

// MessageBroker represents a simple message broker using the STOMP protocol.
type activemq struct {
	conn   *stomp.Conn
	config activemqConfig
	logger *logger.FileLogger
}

// NewMessageBroker creates a new instance of MessageBroker.
func NewMessageBroker(brokerURL string, username string, password string, heartbeat, heartbeatGrace time.Duration) *MessageBroker {

	prop := properties.MustLoadFile("main.properties", properties.UTF8)
	brokerURL = prop.GetString("activemq.broker.url", "localhost:61613")
	username = prop.GetString("activemq.broker.username", "admin")
	password = prop.GetString("activemq.broker.password", "admin")
	heartbeat = time.Second * time.Duration(prop.GetInt("activemq.broker.heartbeat", 10))
	heartbeatGrace = time.Second * time.Duration(prop.GetInt("activemq.broker.heartbeat.grace", 10))

	activemqConfig := activemqConfig{brokerURL, username, password, heartbeat, heartbeatGrace}

	return &activemq{
		conn:   nil,
		config: activemqConfig,
		logger: logger.GetLumberJack(),
	}
}

// Connect connects to the message broker.
func (mb *activemq) Connect() error {
	if mb.conn != nil {
		return fmt.Errorf("already connected")
	}
	options := []func(*stomp.Conn) error{
		stomp.ConnOpt.Login(mb.config.username, mb.config.password),
		stomp.ConnOpt.HeartBeat(mb.config.heartbeat, mb.config.heartbeatGrace),
	}

	conn, err := stomp.Dial("tcp", mb.config.brokerURL, options...)
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
		return nil, fmt.Errorf("not connected to queue: %s ", destination)
	}

	sub, err := mb.conn.Subscribe(destination, stomp.AckAuto)
	if err != nil {
		return nil, err
	}

	return sub, err

}
