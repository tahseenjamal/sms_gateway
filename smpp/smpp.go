package smpp

import (
	"time"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/magiconair/properties"
)

var (
// smppConfig object for storing the configuration of the message broker.

)

// Create smmp configuration struct
type smppConfig struct {
	host       string
	port       int
	systemId   string
	password   string
	systemType string
	window     int
}

type connection struct {
	// forix go smpp library
	conn   *smpp
	config smppConfig
}

func init() {

}

func getConfig() smppConfig {

	prop := properties.MustLoadFile("main.properties", properties.UTF8)
	host := prop.GetString("smpp.host", "localhost")
	port := prop.GetInt("smpp.port", 2775)
	systemId := prop.GetString("smpp.systemId", "systemId")
	password := prop.GetString("smpp.password", "password")
	systemType := prop.GetString("smpp.systemType", "systemType")
	window := prop.GetInt("smpp.window", 1)

	return smppConfig{host, port, systemId, password, systemType, window}
}

// New smpp connection
func New() *connection {

	smppConn := &connection{
		conn:   smppConn.Connect(),
		config: getConfig(),
	}

	return smppConn
}

func (smppConn *connection) Connect() *smpp.Transceiver {

	return smpp.Transceiver{
		Addr:       smppConn.config.host,
		Port:       smppConn.config.port,
		User:       smppConn.config.systemId,
		Passwd:     smppConn.config.password,
		SystemType: smppConn.config.systemType,
		//timeout:    10 * time.Second,
		EnquireLink:        5 * time.Minute,
		EnquireLinkTimeout: 10 * time.Second,
		RespTimeout:        2 * time.Second,
		BindInterval:       10 * time.Second,
		Handler:            smppConn,
		Window:             smppConn.config.window,
	}

}

func (smppConn *connection) Send() {

}

func (smppConn *connection) Receive() {

}

func (smppConn *connection) Reconnect() {

}

func (smppConn *connection) IsConnected() {

}

func (smppConn *connection) Close() {

}
