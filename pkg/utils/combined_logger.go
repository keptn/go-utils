package utils

import (
	"time"

	"github.com/gorilla/websocket"
)

// CombinedLogger logs messages to the logger as well as to the websocket
type CombinedLogger struct {
	logger *Logger
	ws     *websocket.Conn
}

// NewCombinedLogger creates a new CombinedLogger which writes log messages
// to the logger as well as to the websocket
func NewCombinedLogger(logger *Logger, ws *websocket.Conn) *CombinedLogger {
	combinedLogger := CombinedLogger{
		logger: logger,
		ws:     ws,
	}
	return &combinedLogger
}

// Info logs an info message
func (l *CombinedLogger) Info(message string, terminate bool) error {
	if l.logger != nil {
		l.logger.printLogMessage(keptnLogMessage{Timestamp: time.Now(), Message: message, LogLevel: "INFO"})
	}
	if l.ws != nil {
		return WriteLog(l.ws, LogData{LogLevel: "INFO", Message: message, Terminate: terminate})
	}
	return nil
}

// Error logs an error message
func (l *CombinedLogger) Error(message string, terminate bool) error {
	if l.logger != nil {
		l.logger.printLogMessage(keptnLogMessage{Timestamp: time.Now(), Message: message, LogLevel: "ERROR"})
	}
	if l.ws != nil {
		return WriteLog(l.ws, LogData{LogLevel: "ERROR", Message: message, Terminate: terminate})
	}
	return nil
}

// Debug logs a debug message
func (l *CombinedLogger) Debug(message string, terminate bool) error {
	if l.logger != nil {
		l.logger.printLogMessage(keptnLogMessage{Timestamp: time.Now(), Message: message, LogLevel: "DEBUG"})
	}
	if l.ws != nil {
		return WriteLog(l.ws, LogData{LogLevel: "DEBUG", Message: message, Terminate: terminate})
	}
	return nil
}
