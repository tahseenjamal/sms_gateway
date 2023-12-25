package logger

import (
	"log"

	"github.com/magiconair/properties"
	"github.com/natefinch/lumberjack"
)

type configproperties struct {

	// Fetch properties
	filename  string
	maxsize   int
	maxbackup int
	maxage    int
	compress  bool
}

type FileLogger struct {
	configProp configproperties
	logger     *log.Logger
}

// Fetch properties from the properties file
func getFetchProperties() configproperties {

	mainProp := properties.MustLoadFile("main.properties", properties.UTF8)
	prop := properties.MustLoadFile(mainProp.GetString("consumer.properties.filename", ""), properties.UTF8)

	return configproperties{
		filename:  prop.GetString("consumer.log.filename", "consumer.log"),
		maxsize:   prop.GetInt("consumer.log.maxsize", 1024),
		maxbackup: prop.GetInt("consumer.log.maxbackup", 30),
		maxage:    prop.GetInt("consumer.log.maxage", 30),
		compress:  prop.GetBool("consumer.log.compress", true),
	}
}

// GetLumberJack returns a new logger for the consumer package.
func GetLumberJack() *FileLogger {

	var fileLogger *FileLogger
	prop := getFetchProperties()
	fileLogger = &FileLogger{configProp: prop, logger: log.Default()}
	fileLogger.logger.SetOutput(&lumberjack.Logger{
		Filename:   fileLogger.configProp.filename,
		MaxSize:    fileLogger.configProp.maxsize,
		MaxBackups: fileLogger.configProp.maxbackup,
		MaxAge:     fileLogger.configProp.maxage,
		Compress:   fileLogger.configProp.compress,
	})

	// Set the logger for the consumer package
	return fileLogger
}

func (l *FileLogger) WriteLog(data string) {

	l.logger.Println(data)

}
