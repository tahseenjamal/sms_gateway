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

// activemqConfig object for storing the configuration of the message broker.
type activemqConfig struct {
	brokerURL      string
	username       string
	password       string
	heartbeat      time.Duration
	heartbeatGrace time.Duration
}

// activemq object for connecting to the message broker.
type activemq struct {
	conn       *stomp.Conn
	config     activemqConfig
	FileLogger *logger.FileLogger
	subs       *stomp.Subscription
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

	connection := &activemq{
		conn:       nil,
		config:     activemqConfig,
		FileLogger: logger.GetLumberJack(),
	}

	connection.Connect()
	return connection

}

// Connect connects to the message broker.
func (mb *activemq) Connect() {
	if mb.conn != nil {
		time.Sleep(5 * time.Second)
	}
	options := []func(*stomp.Conn) error{
		stomp.ConnOpt.Login(mb.config.username, mb.config.password),
		stomp.ConnOpt.HeartBeat(mb.config.heartbeat, mb.config.heartbeatGrace),
	}

	mutex.Lock()
	defer mutex.Unlock()

	for {

		conn, err := stomp.Dial("tcp", mb.config.brokerURL, options...)
		if err == nil {
			mb.conn = conn
			mb.FileLogger.WriteLog("connected")
			break
		} else {
			fmt.Println("Not connected", err)
			time.Sleep(5 * time.Second)
		}

	}

}

func (mb *activemq) Reconnect(destination string) {
	mb.conn.Disconnect()
	mb.Connect()
	mb.Subscribe(destination)
}

// Sends a message to a specified destination.
func (mb *activemq) Send(destination, body string) {
	if mb.conn == nil {
		fmt.Println("not connected")
		mb.Reconnect()
	}

	err := mb.conn.Send(destination, "text/plain", []byte(body))
	if err != nil {
		mb.FileLogger.WriteLog(fmt.Sprintf("Error sending to destination: %s", destination))
		mb.Reconnect()
	}
}

// Subscribe subscribes to a specified destination after checking if the connection is alive.
func (mb *activemq) Subscribe(destination string) {

	for {

		var err error
		mb.subs, err = mb.conn.Subscribe(destination, stomp.AckAuto)
		if err != nil {
			mb.FileLogger.WriteLog(fmt.Sprintf("Error subscribing to destination: %s", destination))
			mb.Reconnect(destination)
			time.Sleep(1 * time.Second)
		} else {

			break
		}
	}
}

func (mb *activemq) Read(destination string) (string, error) {

	var message *stomp.Message
	var err error
	err = nil
	message, err = mb.subs.Read()
	if err != nil {
		mb.FileLogger.WriteLog(fmt.Sprintf("Error reading destination: %s", err.Error()))
		mb.Reconnect(destination)
	}

	return string(message.Body), err

}
