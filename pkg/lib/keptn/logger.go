package keptn

import (
	"encoding/json"
	"fmt"
	"time"
)

type keptnLogMessage struct {
	Timestamp time.Time `json:"timestamp,string"`
	LogLevel  string    `json:"logLevel"`
	Message   string    `json:"message"`
}

// NewLogger creates a new Logger
func NewLogger(keptnContext string, eventID string, serviceName string) *Logger {
	return &Logger{
		KeptnContext: keptnContext,
		EventID:      eventID,
		ServiceName:  serviceName,
	}
}

// Logger contains data for logging
type Logger struct {
	KeptnContext string `json:"keptnContext"`
	EventID      string `json:"eventId"`
	ServiceName  string `json:"keptnService"`
}

// LoggerInterface collects signatures of the logger
type LoggerInterface interface {
	Info(message string)
	Infof(format string, v ...interface{})
	Error(message string)
	Errorf(format string, v ...interface{})
	Debug(message string)
	Debugf(format string, v ...interface{})
	Terminate(message string)
	Terminatef(format string, v ...interface{})
}

// Info logs an info message
func (l *Logger) Info(message string) {
	l.printLogMessage(keptnLogMessage{Timestamp: time.Now(), Message: message, LogLevel: "INFO"})
}

// Infof formats and logs an info message
func (l *Logger) Infof(format string, v ...interface{}) {
	l.Info(fmt.Sprintf(format, v...))
}

// Error logs an error message
func (l *Logger) Error(message string) {
	l.printLogMessage(keptnLogMessage{Timestamp: time.Now(), Message: message, LogLevel: "ERROR"})
}

// Errorf formats and logs an error message
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Error(fmt.Sprintf(format, v...))
}

// Debug logs a debug message
func (l *Logger) Debug(message string) {
	l.printLogMessage(keptnLogMessage{Timestamp: time.Now(), Message: message, LogLevel: "DEBUG"})
}

// Debugf formats and logs an debug message
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.Debug(fmt.Sprintf(format, v...))
}

// Terminate logs an info message
func (l *Logger) Terminate(message string) {
	l.printLogMessage(keptnLogMessage{Timestamp: time.Now(), Message: message, LogLevel: "INFO"})
}

// Terminatef formats and logs an debug message
func (l *Logger) Terminatef(format string, v ...interface{}) {
	l.Terminate(fmt.Sprintf(format, v...))
}

func (l *Logger) printLogMessage(logMessage keptnLogMessage) {
	logString, err := json.Marshal(logMessage)

	if err != nil {
		fmt.Println("Could not log keptn log message")
		return
	}

	fmt.Println(string(logString))
}
