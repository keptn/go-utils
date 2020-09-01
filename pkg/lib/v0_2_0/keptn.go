package v0_2_0

import (
	"net/url"

	keptn "github.com/keptn/go-utils/pkg/lib/keptn"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"gopkg.in/yaml.v2"
)

type Keptn struct {
	keptn.KeptnBase
}

func NewKeptn(incomingEvent *cloudevents.Event, opts keptn.KeptnOpts) (*Keptn, error) {
	extension, err := incomingEvent.Context.GetExtension("shkeptncontext")
	if err != nil {
		return nil, err
	}
	shkeptncontext := extension.(string)

	// create a base KeptnBase Event
	keptnBase := &EventData{}

	if err := incomingEvent.DataAs(keptnBase); err != nil {
		return nil, err
	}

	k := &Keptn{
		KeptnBase: keptn.KeptnBase{
			Event:              keptnBase,
			KeptnContext:       shkeptncontext,
			UseLocalFileSystem: opts.UseLocalFileSystem,
			ResourceHandler:    nil,
		}}

	csURL := keptn.ConfigurationServiceURL
	if opts.ConfigurationServiceURL != "" {
		csURL = opts.ConfigurationServiceURL
	}

	if opts.EventBrokerURL != "" {
		k.EventBrokerURL = opts.EventBrokerURL
	} else {
		k.EventBrokerURL = keptn.DefaultEventBrokerURL
	}

	k.ResourceHandler = api.NewResourceHandler(csURL)
	k.EventHandler = api.NewEventHandler(csURL)

	loggingServiceName := keptn.DefaultLoggingServiceName
	if opts.LoggingOptions != nil && opts.LoggingOptions.ServiceName != nil {
		loggingServiceName = *opts.LoggingOptions.ServiceName
	}
	k.Logger = keptn.NewLogger(k.KeptnContext, incomingEvent.Context.GetID(), loggingServiceName)

	if opts.LoggingOptions != nil && opts.LoggingOptions.EnableWebsocket {
		wsURL := keptn.DefaultWebsocketEndpoint
		if opts.LoggingOptions.WebsocketEndpoint != nil && *opts.LoggingOptions.WebsocketEndpoint != "" {
			wsURL = *opts.LoggingOptions.WebsocketEndpoint
		}
		connData := keptn.ConnectionData{}
		if err := incomingEvent.DataAs(&connData); err != nil ||
			connData.EventContext.KeptnContext == nil || connData.EventContext.Token == nil ||
			*connData.EventContext.KeptnContext == "" || *connData.EventContext.Token == "" {
			k.Logger.Debug("No WebSocket connection data available")
		} else {
			apiServiceURL, err := url.Parse(wsURL)
			if err != nil {
				k.Logger.Error(err.Error())
				return k, nil
			}
			ws, _, err := keptn.OpenWS(connData, *apiServiceURL)
			if err != nil {
				k.Logger.Error("Opening WebSocket connection failed:" + err.Error())
				return k, nil
			}
			stdLogger := keptn.NewLogger(shkeptncontext, incomingEvent.Context.GetID(), loggingServiceName)
			combinedLogger := keptn.NewCombinedLogger(stdLogger, ws, shkeptncontext)
			k.Logger = combinedLogger
		}
	}

	return k, nil
}

// GetShipyard returns the shipyard definition of a project
func (k *Keptn) GetShipyard() (*Shipyard, error) {
	shipyardResource, err := k.ResourceHandler.GetProjectResource(k.Event.GetProject(), "shipyard.yaml")
	if err != nil {
		return nil, err
	}

	shipyard := Shipyard{}
	err = yaml.Unmarshal([]byte(shipyardResource.ResourceContent), &shipyard)
	if err != nil {
		return nil, err
	}
	return &shipyard, nil
}
