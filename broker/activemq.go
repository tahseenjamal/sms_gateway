package broker

import (
	"fmt"
	"sync"
	"time"

	logger "sms_gateway/logger"

	"github.com/go-stomp/stomp"
	"github.com/magiconair/properties"
)

// Mutex for synchronizing access to the connection.
var mutex = &sync.Mutex{}

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
	subs   *stomp.Subscription
}

func getConfig() activemqConfig {

	prop := properties.MustLoadFile("main.properties", properties.UTF8)
	brokerURL := prop.GetString("activemq.broker.url", "localhost:61613")
	username := prop.GetString("activemq.broker.username", "admin")
	password := prop.GetString("activemq.broker.password", "admin")
	heartbeat := time.Second * time.Duration(prop.GetInt("activemq.broker.heartbeat", 10))
	heartbeatGrace := time.Second * time.Duration(prop.GetInt("activemq.broker.heartbeat.grace", 10))

	return activemqConfig{brokerURL, username, password, heartbeat, heartbeatGrace}
}

// NewMessageBroker creates a new instance of MessageBroker.
func NewMessageBroker() *activemq {

	activemqConfig := getConfig()

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

	mutex.Lock()
	defer mutex.Unlock()

	if mb.conn != nil {
		return fmt.Errorf("already connected")

	}

	for {

		conn, err := stomp.Dial("tcp", mb.config.brokerURL, options...)
		if err != nil {

			fmt.Println("already connected")
		} else {

			mb.conn = conn
			break
		}

	}

	return nil

}

// Disconnect disconnects from the message broker.
func (mb *activemq) Disconnect() error {
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
		fmt.Println("not connected")
		mb.Connect()
	}

	err := mb.conn.Send(destination, "text/plain", []byte(body))
	if err != nil {
		fmt.Printf("cannot send to queue: %s\n", destination)
		mb.Connect()
	}
	return nil
}

// Subscribe subscribes to a specified destination after checking if the connection is alive.
func (mb *activemq) Subscribe(destination string) {

	for {

		if mb.conn == nil {
			mb.Connect()
		}

		var err error
		mb.subs, err = mb.conn.Subscribe(destination, stomp.AckAuto)

		if err != nil {
			fmt.Println("Error subscribing to destination: ", destination)
		} else {

			break
		}
	}
}
