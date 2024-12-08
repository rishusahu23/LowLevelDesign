package loger

import (
	"fmt"
	"os"
	"sync"
)

type LogLevel int

const (
	Info LogLevel = iota
	Debug
	Error
)

type LoggerConfig struct {
	LogLevel    LogLevel
	MaxFileSize int64
	LogFilepath string
}

type Logger interface {
	Log(level LogLevel, message string)
	SetLogLevel(level LogLevel)
}

type BaseLogger struct {
	config LoggerConfig
	mutex  sync.Mutex
}

func (bl *BaseLogger) SetLogLevel(level LogLevel) {
	bl.mutex.Lock()
	defer bl.mutex.Unlock()
	bl.config.LogLevel = level
}

func (bl *BaseLogger) ShouldLog(level LogLevel) bool {
	return level >= bl.config.LogLevel
}

func (bl *BaseLogger) formatLogMessage(level LogLevel, message string) string {
	var levelStr string
	switch level {
	case Info:
		levelStr = "info"
	case Error:
		levelStr = "error"
	case Debug:
		levelStr = "debug"
	}
	return levelStr + message
}

type ConsoleLogger struct {
	BaseLogger
}

func NewConsoleLogger(config LoggerConfig) *ConsoleLogger {
	return &ConsoleLogger{
		BaseLogger: BaseLogger{
			config: config,
		},
	}
}

func (cl *ConsoleLogger) Log(level LogLevel, message string) {
	if !cl.ShouldLog(level) {
		return
	}
	logMsg := cl.formatLogMessage(level, message)
	fmt.Println(logMsg)
}

type FileLogger struct {
	BaseLogger
	file     *os.File
	fileSize int64
}

func NewFileLogger(config LoggerConfig) (*FileLogger, error) {
	file, err := os.OpenFile(config.LogFilepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, err
	}
	return &FileLogger{
		BaseLogger: BaseLogger{config: config},
		file:       file,
	}, nil
}

func (fl *FileLogger) rotateFile() error {
	fl.file.Close()
	err := os.Rename(fl.config.LogFilepath, fl.config.LogFilepath+".old")
	if err != nil {
		return err
	}
	fl.file, err = os.OpenFile(fl.config.LogFilepath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	return err
}

func (fl *FileLogger) Log(level LogLevel, message string) {
	if !fl.ShouldLog(level) {
		return
	}
	logMsg := fl.formatLogMessage(level, message)
	fl.mutex.Lock()
	defer fl.mutex.Unlock()

	if fl.fileSize >= fl.config.MaxFileSize {
		if err := fl.rotateFile(); err != nil {
			fmt.Println("Error rotating file")
			return
		}
		fl.fileSize = 0
	}
	_, err := fl.file.WriteString(logMsg + "\n")
	if err != nil {
		return
	}
	fl.fileSize += int64(len(logMsg))
}

func (fl *FileLogger) Close() {
	fl.file.Close()
}

type CompositeLogger struct {
	loggers []Logger
}

func NewCompositeLogger(loggers ...Logger) *CompositeLogger {
	return &CompositeLogger{
		loggers: loggers,
	}
}

func (cl *CompositeLogger) Log(level LogLevel, message string) {
	for _, logger := range cl.loggers {
		logger.Log(level, message)
	}
}

func (cl *CompositeLogger) SetLogLevel(level LogLevel) {
	for _, logger := range cl.loggers {
		logger.SetLogLevel(level)
	}
}

func main() {
	consoleConfig := LoggerConfig{
		LogLevel: Info,
	}
	fileConfig := LoggerConfig{
		LogLevel:    Debug,
		MaxFileSize: 1024 * 10,
		LogFilepath: "app.log",
	}

	consoleLogger := NewConsoleLogger(consoleConfig)
	fileLogger, err := NewFileLogger(fileConfig)
	if err != nil {
		fmt.Printf("Error initializing file logger: %v\n", err)
		return
	}
	defer fileLogger.Close()

	compositeLogger := NewCompositeLogger(consoleLogger, fileLogger)

	// Log messages at different levels
	compositeLogger.Log(Info, "This is an info message.")
	compositeLogger.Log(Debug, "This is a debug message.")
	compositeLogger.Log(Error, "This is an error message.")

	// Changing log level to debug
	compositeLogger.SetLogLevel(Debug)
	compositeLogger.Log(Info, "This is another info message at debug level.")
	compositeLogger.Log(Debug, "This is another debug message.")
}
