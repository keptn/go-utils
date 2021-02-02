package v0_2_0

import (
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	"gopkg.in/yaml.v2"
)

type Keptn struct {
	keptn.KeptnBase
}

func NewKeptn(incomingEvent *cloudevents.Event, opts keptn.KeptnOpts) (*Keptn, error) {
	extension, _ := incomingEvent.Context.GetExtension("shkeptncontext")
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

	if opts.EventBrokerURL != "" && opts.EventSender == nil {
		httpSender, err := NewHTTPEventSender(opts.EventBrokerURL)
		if err != nil {
			return nil, fmt.Errorf("could not initialize Keptn Handler: %s", err.Error())
		}
		k.KeptnBase.EventSender = httpSender
	} else if opts.EventSender != nil {
		k.KeptnBase.EventSender = opts.EventSender
	} else {
		httpSender, err := NewHTTPEventSender(DefaultHTTPEventEndpoint)
		if err != nil {
			return nil, fmt.Errorf("could not initialize Keptn Handler: %s", err.Error())
		}
		k.KeptnBase.EventSender = httpSender
	}

	datastoreURL := keptn.DatastoreURL
	if opts.DatastoreURL != "" {
		datastoreURL = opts.DatastoreURL
	}

	k.ResourceHandler = api.NewResourceHandler(csURL)
	k.EventHandler = api.NewEventHandler(datastoreURL)

	loggingServiceName := keptn.DefaultLoggingServiceName
	if opts.LoggingOptions != nil && opts.LoggingOptions.ServiceName != nil {
		loggingServiceName = *opts.LoggingOptions.ServiceName
	}
	k.Logger = keptn.NewLogger(k.KeptnContext, incomingEvent.Context.GetID(), loggingServiceName)

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
