package utils

import (
	"encoding/json"
	"fmt"
)

type keptnLogMessage struct {
	KeptnContext string `json:"keptnContext"`
	Message      string `json:"message"`
	KeptnService string `json:"keptnService"`
	LogLevel     string `json:"logLevel"`
}

// ServiceName stores the name of the service
var ServiceName = ""

// Info logs an info message
func Info(keptnContext string, message string) {
	printLogMessage(keptnLogMessage{KeptnContext: keptnContext, Message: message, KeptnService: ServiceName, LogLevel: "INFO"})
}

// Error logs an error message
func Error(keptnContext string, message string) {
	printLogMessage(keptnLogMessage{KeptnContext: keptnContext, Message: message, KeptnService: ServiceName, LogLevel: "ERROR"})
}

// Debug logs a debug message
func Debug(keptnContext string, message string) {
	printLogMessage(keptnLogMessage{KeptnContext: keptnContext, Message: message, KeptnService: ServiceName, LogLevel: "DEBUG"})
}

func printLogMessage(logMessage keptnLogMessage) {
	logString, err := json.Marshal(logMessage)

	if err != nil {
		fmt.Println("Could not log keptn log message")
		return
	}

	fmt.Println(string(logString))
}
