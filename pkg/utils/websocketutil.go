package utils

import (
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/cloudevents/sdk-go/pkg/cloudevents/types"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

// MyCloudEvent represents a keptn cloud event
type MyCloudEvent struct {
	SpecVersoin    string          `json:"specversion"`
	ContentType    string          `json:"contentType"`
	Data           json.RawMessage `json:"data"`
	ID             string          `json:"id"`
	Time           string          `json:"time"`
	Type           string          `json:"type"`
	Source         string          `json:"source"`
	ShKeptnContext string          `json:"shkeptncontext"`
}

// LogData represents log data
type LogData struct {
	Message   string `json:"message"`
	Terminate bool   `json:"terminate"`
	LogLevel  string `json:"loglevel"`
}

// IncompleteCE is a helper type for unmarshalling the CE data
type IncompleteCE struct {
	ConnData ConnectionData `json:"data"`
}

// ConnectionData stores ChannelInfo and Success data
type ConnectionData struct {
	ChannelInfo ChannelInfo `json:"channelInfo"`
}

// ChannelInfo stores a token and a channelID used for opening the websocket
type ChannelInfo struct {
	Token     string `json:"token"`
	ChannelID string `json:"channelID"`
}

// OpenWS opens a websocket
func OpenWS(connData ConnectionData, apiEndPoint url.URL) (*websocket.Conn, *http.Response, error) {

	wsEndPoint := apiEndPoint
	wsEndPoint.Scheme = "ws"

	header := http.Header{}
	header.Add("Token", connData.ChannelInfo.Token)

	dialer := websocket.DefaultDialer
	dialer.HandshakeTimeout = 120 * time.Second

	return dialer.Dial(wsEndPoint.String(), header)
}

// WriteWSLog writes the log event to the websocket
func WriteWSLog(ws *websocket.Conn, logEvent cloudevents.Event, message string, terminate bool, logLevel string) error {
	logData := LogData{
		Message:   message,
		Terminate: terminate,
		LogLevel:  logLevel,
	}
	logDataRaw, _ := json.Marshal(logData)

	var shkeptncontext string
	logEvent.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	messageCE := MyCloudEvent{
		SpecVersoin:    logEvent.SpecVersion(),
		ContentType:    logEvent.DataContentType(),
		Data:           logDataRaw,
		ID:             logEvent.ID(),
		Time:           logEvent.Time().String(),
		Type:           "sh.keptn.events.log",
		Source:         logEvent.Source(),
		ShKeptnContext: shkeptncontext,
	}

	data, _ := json.Marshal(messageCE)
	return ws.WriteMessage(1, data) // websocket.TextMessage = 1; ws.WriteJSON not supported because keptn CLI does a ReadMessage
}

// WriteLog writes the logData to the websocket connection
func WriteLog(ws *websocket.Conn, logData LogData, shkeptnContext string) error {

	logDataRaw, _ := json.Marshal(logData)
	now := &types.Timestamp{Time: time.Now()}

	messageCE := MyCloudEvent{
		SpecVersoin:    "0.2",
		ContentType:    "application/json",
		Data:           logDataRaw,
		ID:             uuid.New().String(),
		Time:           now.String(),
		Type:           "sh.keptn.events.log",
		Source:         "https://github.com/keptn/keptn",
		ShKeptnContext: shkeptnContext,
	}

	data, _ := json.Marshal(messageCE)
	return ws.WriteMessage(1, data) // websocket.TextMessage = 1; ws.WriteJSON not supported because keptn CLI does a ReadMessage
}
