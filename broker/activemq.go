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
	conn        *stomp.Conn
	config      activemqConfig
	fileloggger *logger.FileLogger
	subs        *stomp.Subscription
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
		conn:        nil,
		config:      activemqConfig,
		fileloggger: logger.GetLumberJack(),
	}

	connection.Connect()
	return connection

}

// Connect connects to the message broker.
func (mb *activemq) Connect() {
	if mb.conn != nil {
		time.Sleep(5 * time.Second)
		mb.fileloggger.WriteLog("inside connection: already connected")
	}
	options := []func(*stomp.Conn) error{
		stomp.ConnOpt.Login(mb.config.username, mb.config.password),
		stomp.ConnOpt.HeartBeat(mb.config.heartbeat, mb.config.heartbeatGrace),
	}

	mutex.Lock()
	defer mutex.Unlock()

	if mb.conn != nil {
		time.Sleep(5 * time.Second)
		mb.fileloggger.WriteLog("second check inside connection: already connected")

	}

	for {

		conn, err := stomp.Dial("tcp", mb.config.brokerURL, options...)
		if err == nil {
			mb.conn = conn
			mb.fileloggger.WriteLog("connected")
			break
		} else {
			fmt.Println("Not connected", err)
			time.Sleep(5 * time.Second)
		}

	}

}

func (mb *activemq) Reconnect() {
	mb.conn.Disconnect()
	mb.Connect()
}

// Sends a message to a specified destination.
func (mb *activemq) Send(destination, body string) {
	if mb.conn == nil {
		fmt.Println("not connected")
		mb.Reconnect()
	}

	err := mb.conn.Send(destination, "text/plain", []byte(body))
	if err != nil {
		mb.fileloggger.WriteLog(fmt.Sprintf("Error sending to destination: %s", destination))
		mb.Reconnect()
	}
}

// Subscribe subscribes to a specified destination after checking if the connection is alive.
func (mb *activemq) Subscribe(destination string) {

	for {

		mb.fileloggger.WriteLog(fmt.Sprintf("Attempting to subscribe to destination: %s", destination))
		var err error
		mb.subs, err = mb.conn.Subscribe(destination, stomp.AckAuto)
		mb.fileloggger.WriteLog(fmt.Sprintf("Subscribe executed Successfully to destination: %s", destination))

		if err != nil {
			mb.fileloggger.WriteLog(fmt.Sprintf("Error subscribing to destination: %s", destination))
			mb.Reconnect()
			time.Sleep(1 * time.Second)
		} else {

			break
		}
	}
}

func (mb *activemq) Read(destination string) string {

	var message *stomp.Message
	var err error
	for {
		message, err = mb.subs.Read()
		if err != nil {
			mb.fileloggger.WriteLog(fmt.Sprintf("Error reading destination: %s", destination))
			mb.Reconnect()
			mb.fileloggger.WriteLog(fmt.Sprintf("Subscribing again to destination: %s", destination))
			mb.Subscribe(destination)
			mb.fileloggger.WriteLog("Back from subscription")
		} else {
			break
		}

	}
	return string(message.Body)

}
