package smppconnector

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/fiorix/go-smpp/smpp"
	"github.com/fiorix/go-smpp/smpp/encoding"
	"github.com/fiorix/go-smpp/smpp/pdu"
	"github.com/fiorix/go-smpp/smpp/pdu/pdufield"
	"github.com/fiorix/go-smpp/smpp/pdu/pdutext"
	"github.com/magiconair/properties"
	"golang.org/x/time/rate"
)

var (
	re        *regexp.Regexp
	startTime time.Time
	endTime   time.Time

	rate_limiter *rate.Limiter
)

// Create smmp configuration struct
type smppConfig struct {
	host       string
	port       int
	systemId   string
	password   string
	systemType string
	window     uint
	srcTON     uint8
	srcNPI     uint8
	dstTON     uint8
	dstNPI     uint8
}

type connection struct {
	// forix go smpp library
	conn   *smpp.Transceiver
	config smppConfig
}

func splitString(input string, delimiter string) (int, int) {
	// Receive a hour minute string separated by :
	// Split and return
	parts := strings.Split(input, delimiter)
	h, _ := strconv.Atoi(parts[0])
	m, _ := strconv.Atoi(parts[1])
	return h, m
}

func init() {

	prop := properties.MustLoadFile("main.properties", properties.UTF8)

	pattern := `id:(\w+) sub:(\d+) dlvrd:(\d+) submit date:(\d+) done date:(\d+) stat:(\w+) err:(\d+) [Tt]ext:(?i)(.+)`
	re = regexp.MustCompile(pattern)

	rate_limiter = rate.NewLimiter(rate.Every(time.Duration(1000/prop.GetUint("tps", 50))*time.Millisecond), 1)

	morningHour, morningMinute := splitString(prop.GetString("smpp.morning", "9:00"), ":")
	eveningHour, eveningMinute := splitString(prop.GetString("smpp.evening", "20:00"), ":")

	// start and end time basis black hour, messages outside this time would be dropped
	currentTime := time.Now()
	startTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), morningHour, morningMinute, 0, 0, currentTime.Location())
	endTime = time.Date(currentTime.Year(), currentTime.Month(), currentTime.Day(), eveningHour, eveningMinute, 0, 0, currentTime.Location())

}

func extract(message string) (map[string]string, error) {

	matches := re.FindStringSubmatch(message)

	var resultMap = make(map[string]string)

	if len(matches) > 0 {
		keys := []string{"id", "sub", "dlvrd", "submit_date", "done_date", "stat", "err", "text"}

		for i, key := range keys {
			resultMap[key] = matches[i+1]
		}

		return resultMap, nil
	}

	return nil, fmt.Errorf("invalid data length")

}

func getConfig() smppConfig {

	prop := properties.MustLoadFile("main.properties", properties.UTF8)
	host := prop.GetString("smpp.host", "localhost")
	port := prop.GetInt("smpp.port", 2775)
	systemId := prop.GetString("smpp.systemId", "systemId")
	password := prop.GetString("smpp.password", "password")
	systemType := prop.GetString("smpp.systemType", "systemType")
	window := prop.GetUint("smpp.window", 1)
	srcTON := uint8(prop.GetUint("smpp.srcTON", 5))
	srcNPI := uint8(prop.GetUint("smpp.srcNPI", 1))
	dstTON := uint8(prop.GetUint("smpp.dstTON", 1))
	dstNPI := uint8(prop.GetUint("smpp.dstNPI", 1))

	return smppConfig{host, port, systemId, password, systemType, window, srcTON, srcNPI, dstTON, dstNPI}
}

// New smpp connection
func NewSmpp() *connection {

	smppConn := &connection{
		conn:   nil,
		config: getConfig(),
	}

	smppConn.Connect()

	return smppConn
}

func (smppConn *connection) Connect() <-chan smpp.ConnStatus {

	smppConn.conn = &smpp.Transceiver{
		Addr:       smppConn.config.host + ":" + strconv.Itoa(smppConn.config.port),
		User:       smppConn.config.systemId,
		Passwd:     smppConn.config.password,
		SystemType: smppConn.config.systemType,
		//timeout:    10 * time.Second,
		EnquireLink:        1 * time.Minute,
		EnquireLinkTimeout: 10 * time.Second,
		RespTimeout:        2 * time.Second,
		BindInterval:       30 * time.Second,
		Handler:            smppConn.Receive,
		RateLimiter:        rate_limiter,
		WindowSize:         smppConn.config.window,
	}

	return smppConn.conn.Bind()

}

func (smppConn *connection) Send(sender string, dest string, message string, test string) {

	if smppConn.IsBlackHour() {
		fmt.Println("|BLACK_HOUR|%s|%s|%s|%s", sender, dest, message)
		return
	}

	if len(encoding.ValidateGSM7String(message)) > 0 || len(message) > 160 {
		sml, err := smppConn.submitLong(sender, dest, message)
		if err == nil {
			for _, sm := range sml {
				fmt.Println("|SUBMITTED|%s|%s|%s|%s", sender, dest, message, sm.RespID())
			}
		} else {
			fmt.Println("|SMPP_ERROR|%s|%s|%s|%s", sender, dest, message, err.Error())
		}
	} else {
		sm, err := smppConn.submitShort(sender, dest, message)

		if err == nil {
			fmt.Println("|SUBMITTED|%s|%s|%s|%s", sender, dest, message, sm.RespID())
		} else {
			fmt.Println("|SMPP_ERROR|%s|%s|%s|%s", sender, dest, message, err.Error())
		}
	}

}

func (smppConn *connection) submitShort(sender string, dest string, message string) (*smpp.ShortMessage, error) {
	pduMessage := pdutext.Raw(message)
	sm, err := smppConn.conn.Submit(&smpp.ShortMessage{
		Src:           sender,
		Dst:           dest,
		Text:          pduMessage,
		Register:      pdufield.FinalDeliveryReceipt,
		SourceAddrTON: smppConn.config.srcTON,
		SourceAddrNPI: smppConn.config.dstTON,
		DestAddrTON:   smppConn.config.dstTON,
		DestAddrNPI:   smppConn.config.dstTON,
	})

	return sm, err
}

func (smppConn *connection) submitLong(sender string, dest string, message string) ([]smpp.ShortMessage, error) {

	pduMessage := pdutext.UCS2(message)
	sml, err := smppConn.conn.SubmitLongMsg(&smpp.ShortMessage{
		Src:           sender,
		Dst:           dest,
		Text:          pduMessage,
		Register:      pdufield.FinalDeliveryReceipt,
		SourceAddrTON: smppConn.config.srcTON,
		SourceAddrNPI: smppConn.config.dstTON,
		DestAddrTON:   smppConn.config.dstTON,
		DestAddrNPI:   smppConn.config.dstTON,
		ESMClass:      8,
	})

	return sml, err

}

func (smppConn *connection) IsConnected() {

}

func (smppConn *connection) Receive(p pdu.Body) {
	switch p.Header().ID {
	case pdu.DeliverSMID:
		f := p.Fields()
		dst := f[pdufield.SourceAddr]
		src := f[pdufield.DestinationAddr]
		text := f[pdufield.ShortMessage].String()
		dlr, _ := extract(text)

		fmt.Sprintf("|SMPP_RESPONSE|%s|+%s|%s|%s|%s", src, dst, dlr["stat"], dlr["text"], dlr["id"])
	}

}

func (smppConn *connection) IsBlackHour() bool {

	return time.Now().Before(startTime) || time.Now().After(endTime)
}
