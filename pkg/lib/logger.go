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
	Error(message string)
	Debug(message string)
}

// Info logs an info message
func (l *Logger) Info(message string) {
	l.printLogMessage(keptnLogMessage{Timestamp: time.Now(), Message: message, LogLevel: "INFO"})
}

// Error logs an error message
func (l *Logger) Error(message string) {
	l.printLogMessage(keptnLogMessage{Timestamp: time.Now(), Message: message, LogLevel: "ERROR"})
}

// Debug logs a debug message
func (l *Logger) Debug(message string) {
	l.printLogMessage(keptnLogMessage{Timestamp: time.Now(), Message: message, LogLevel: "DEBUG"})
}

func (l *Logger) printLogMessage(logMessage keptnLogMessage) {
	logString, err := json.Marshal(logMessage)

	if err != nil {
		fmt.Println("Could not log keptn log message")
		return
	}

	fmt.Println(string(logString))
}
