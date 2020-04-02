package keptn

import (
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

// CombinedLogger logs messages to the logger as well as to the websocket
type CombinedLogger struct {
	logger         *Logger
	ws             *websocket.Conn
	shKeptnContext string
}

// NewCombinedLogger creates a new CombinedLogger which writes log messages
// to the logger as well as to the websocket
func NewCombinedLogger(logger *Logger, ws *websocket.Conn, shKeptnContext string) *CombinedLogger {
	combinedLogger := CombinedLogger{
		logger:         logger,
		ws:             ws,
		shKeptnContext: shKeptnContext,
	}
	return &combinedLogger
}

// Info logs an info message
func (l *CombinedLogger) Info(message string) {
	l.logger.printLogMessage(keptnLogMessage{Timestamp: time.Now(), Message: message, LogLevel: "INFO"})
	if err := WriteLog(l.ws, LogData{LogLevel: "INFO", Message: message, Terminate: false}, l.shKeptnContext); err != nil {
		l.logWebsocketError(err)
	}
}

// Error logs an error message
func (l *CombinedLogger) Error(message string) {
	l.logger.printLogMessage(keptnLogMessage{Timestamp: time.Now(), Message: message, LogLevel: "ERROR"})
	if err := WriteLog(l.ws, LogData{LogLevel: "ERROR", Message: message, Terminate: false}, l.shKeptnContext); err != nil {
		l.logWebsocketError(err)
	}
}

// Debug logs a debug message
func (l *CombinedLogger) Debug(message string) {
	l.logger.printLogMessage(keptnLogMessage{Timestamp: time.Now(), Message: message, LogLevel: "DEBUG"})
	if err := WriteLog(l.ws, LogData{LogLevel: "DEBUG", Message: message, Terminate: false}, l.shKeptnContext); err != nil {
		l.logWebsocketError(err)
	}
}

// Terminate sends a terminate message to the websocket
func (l *CombinedLogger) Terminate() {
	if err := WriteLog(l.ws, LogData{LogLevel: "INFO", Message: "", Terminate: true}, l.shKeptnContext); err != nil {
		l.logWebsocketError(err)
	}
}

func (l *CombinedLogger) logWebsocketError(err error) {
	l.logger.printLogMessage(keptnLogMessage{Timestamp: time.Now(),
		Message:  fmt.Sprintf("Websocket error when writing message: %s", err.Error()),
		LogLevel: "ERROR"})
}
