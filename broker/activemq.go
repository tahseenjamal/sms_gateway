package broker

import (
	"sync"
	"time"

	logger "sms_gateway/logger"

	"github.com/go-stomp/stomp"
	"github.com/magiconair/properties"
)

// Mutex for synchronizing access to the connection.
var (
	mutex *sync.Mutex
)

const (
	NULLSTRING = ""
)

// activemqConfig object for storing the configuration of the message broker.
type activemqConfig struct {
	brokerURL      string
	username       string
	password       string
	heartbeat      time.Duration
	heartbeatGrace time.Duration
}

// activemq object for connecting to the message broker.
type Activemq struct {
	conn       *stomp.Conn
	config     activemqConfig
	FileLogger *logger.FileLogger
	subs       *stomp.Subscription
}

func init() {
	mutex = &sync.Mutex{}
}

func getConfig() activemqConfig {

	prop := properties.MustLoadFile("main.properties", properties.UTF8)
	brokerURL := prop.GetString("activemq.broker.url", "localhost:61613")
	username := prop.GetString("activemq.broker.username", "admin")
	password := prop.GetString("activemq.broker.password", "admin")
	heartbeat := time.Millisecond * time.Duration(prop.GetInt("activemq.broker.heartbeat", 500))
	heartbeatGrace := time.Second * time.Duration(prop.GetInt("activemq.broker.heartbeat.grace", 5))

	return activemqConfig{brokerURL, username, password, heartbeat, heartbeatGrace}
}

// NewMessageBroker creates a new instance of MessageBroker.
func NewMessageBroker() *Activemq {

	activemqConfig := getConfig()

	connection := &Activemq{
		conn:       nil,
		config:     activemqConfig,
		FileLogger: logger.GetLumberJack(),
	}

	connection.Connect()
	return connection

}

// Connect connects to the message broker.
func (mb *Activemq) Connect() {

	options := []func(*stomp.Conn) error{
		stomp.ConnOpt.Login(mb.config.username, mb.config.password),
		stomp.ConnOpt.HeartBeat(mb.config.heartbeat, mb.config.heartbeatGrace),
	}

	for {

		conn, err := stomp.Dial("tcp", mb.config.brokerURL, options...)
		if err == nil {
			mb.conn = conn
			mb.FileLogger.WriteLog("connected")
			time.Sleep(1 * time.Second)
			break
		} else {
			mb.FileLogger.WriteLog("Error connecting to broker: %s", err.Error())
			time.Sleep(1 * time.Second)
		}

	}
}

// Can also launch this separately as a goroutine to check if the connection is alive.
// If not, it will reconnect. But if for any reasons the goroutine is killed, the connection status will not be checked.
func (mb *Activemq) Reconnect(destination string) {

	mutex.Lock()
	defer mutex.Unlock()

	if mb.conn != nil {
		t, e := mb.conn.Subscribe("/queue/heartbeat", stomp.AckAuto)
		if e != nil {

			mb.conn.Disconnect()
			mb.Connect()
			if destination != "" {
				mb.Subscribe(destination)
			}
		} else {
			t.Unsubscribe()
		}
	}
}

// Sends a message to a specified destination.
func (mb *Activemq) Send(destination, body string) error {

	err := mb.conn.Send(destination, "text/plain", []byte(body))
	if err != nil {
		mb.FileLogger.WriteLog("|BROKER_ERROR|Error sending to: %s", destination)
		time.Sleep(1 * time.Second)
		mb.Reconnect("")
	}

	return err
}

// Subscribe subscribes to a specified destination after checking if the connection is alive.
func (mb *Activemq) Subscribe(destination string) {

	var err error
	mb.subs, err = mb.conn.Subscribe(destination, stomp.AckAuto)
	if err != nil {
		mb.FileLogger.WriteLog("|BROKER_ERROR|Error subscribing to: %s", destination)
		time.Sleep(1 * time.Second)
		mb.Reconnect(destination)
	}
}

func (mb *Activemq) Read(destination string) (string, error) {

	var message *stomp.Message
	var err error
	message, err = mb.subs.Read()
	if err != nil {
		mb.FileLogger.WriteLog("Error reading from: %s", err.Error())
		time.Sleep(1 * time.Second)
		mb.Reconnect(destination)
		return NULLSTRING, err
	}
	return string(message.Body), err

}
