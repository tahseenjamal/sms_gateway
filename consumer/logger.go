package consumer

import (
	"log"

	"github.com/magiconair/properties"
	"github.com/natefinch/lumberjack"
)

// ConsumerLogger represents a logger for the consumer package.
type logger struct {
	loggerProperties loggerproperties
	ConsumerLogger   *log.Logger
	loggerQueue      chan string
}

type loggerproperties struct {

	// Fetch properties
	filename  string
	maxsize   int
	maxbackup int
	maxage    int
	compress  bool
	queuesize int
}

// Fetch properties from the properties file
func getFetchProperties() loggerproperties {

	prop := properties.MustLoadFile("consumer.properties", properties.UTF8)

	return loggerproperties{
		filename:  prop.GetString("consumer.log.filename", "consumer.log"),
		maxsize:   prop.GetInt("consumer.log.maxsize", 1024),
		maxbackup: prop.GetInt("consumer.log.maxbackup", 30),
		maxage:    prop.GetInt("consumer.log.maxage", 30),
		compress:  prop.GetBool("consumer.log.compress", true),
	}
}

// GetLumberJack returns a new logger for the consumer package.
func GetLumberJack() logger {

	var fileLogger logger
	fileLogger.loggerProperties = getFetchProperties()

	fileLogger.ConsumerLogger = log.Default()
	fileLogger.ConsumerLogger.SetOutput(&lumberjack.Logger{
		Filename:   fileLogger.loggerProperties.filename,
		MaxSize:    fileLogger.loggerProperties.maxsize,
		MaxBackups: fileLogger.loggerProperties.maxbackup,
		MaxAge:     fileLogger.loggerProperties.maxage,
		Compress:   fileLogger.loggerProperties.compress,
	})

	fileLogger.loggerQueue = make(chan string, fileLogger.loggerProperties.queuesize)

	// Set the logger for the consumer package
	return fileLogger
}

func (l *logger) WriteLog(data string) {

	l.ConsumerLogger.Printf("Message received: %s", data)

}
