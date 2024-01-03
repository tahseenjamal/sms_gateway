package broker

import (
	"sync"
	"time"

	logger "sms_gateway/logger"

	"github.com/go-stomp/stomp"
	"github.com/magiconair/properties"
)

type unofficialStompConn struct {
	Conn                     interface{}   `json:"conn"`
	ReadCh                   chan string   `json:"readCh"`
	WriteCh                  chan int      `json:"writeCh"`
	Version                  float64       `json:"version"`
	Session                  string        `json:"session"`
	Server                   string        `json:"server"`
	ReadTimeout              time.Duration `json:"readTimeout"`
	WriteTimeout             time.Duration `json:"writeTimeout"`
	MsgSendTimeout           time.Duration `json:"msgSendTimeout"`
	RcvReceiptTimeout        time.Duration `json:"rcvReceiptTimeout"`
	DisconnectReceiptTimeout time.Duration `json:"disconnectReceiptTimeout"`
	HbGracePeriodMultiplier  float64       `json:"hbGracePeriodMultiplier"`
	Closed                   bool          `json:"closed"`
	CloseMutex               interface{}   `json:"closeMutex"`
	Options                  interface{}   `json:"options"`
	Log                      interface{}   `json:"log"`
}

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
type activemq struct {
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
			time.Sleep(1 * time.Second)
			break
		} else {
			mb.FileLogger.WriteLog("Error connecting to broker: %s", err.Error())
			time.Sleep(1 * time.Second)
		}

	}
}

func (mb *activemq) Reconnect(destination string) {

	mb.conn.Disconnect()
	mb.Connect()
	if destination != "" {
		mb.Subscribe(destination)
	}
}

// Sends a message to a specified destination.
func (mb *activemq) Send(destination, body string) error {

	err := mb.conn.Send(destination, "text/plain", []byte(body))
	if err != nil {
		mb.FileLogger.WriteLog("|BROKER_ERROR|Error sending to destination: %s", destination)
		go mb.Reconnect("")
	}

	return err
}

// Subscribe subscribes to a specified destination after checking if the connection is alive.
func (mb *activemq) Subscribe(destination string) {

	var err error
	mb.subs, err = mb.conn.Subscribe(destination, stomp.AckAuto)
	if err != nil {
		mb.FileLogger.WriteLog("|BROKER_ERROR|Error subscribing to destination: %s", destination)
		go mb.Reconnect(destination)
	}
}

func (mb *activemq) Read(destination string) (string, error) {

	var message *stomp.Message
	var err error
	err = nil
	message, err = mb.subs.Read()
	if err != nil {
		mb.FileLogger.WriteLog("Error reading destination: %s", err.Error())
		go mb.Reconnect(destination)
		return NULLSTRING, err
	} else {
		return string(message.Body), err
	}

}
