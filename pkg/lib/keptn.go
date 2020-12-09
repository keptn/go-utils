package keptn

import (
	"encoding/json"
	"os"
	"strings"

	keptn "github.com/keptn/go-utils/pkg/lib/keptn"

	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"gopkg.in/yaml.v2"
)

type Keptn struct {
	keptn.KeptnBase
}

func NewKeptn(incomingEvent *cloudevents.Event, opts keptn.KeptnOpts) (*Keptn, error) {
	var shkeptncontext string
	_ = incomingEvent.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	// create a base KeptnBase Event
	keptnBase := &KeptnBaseEvent{}

	bytes, err := incomingEvent.DataBytes()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, keptnBase)
	if err != nil {
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

//
// replaces $ placeholders with actual values
// $CONTEXT, $EVENT, $SOURCE
// $PROJECT, $STAGE, $SERVICE, $DEPLOYMENT
// $TESTSTRATEGY
// $LABEL.XXXX  -> will replace that with a label called XXXX
// $ENV.XXXX    -> will replace that with an env variable called XXXX
//
func (k *Keptn) ReplaceKeptnPlaceholders(input string) string {
	result := input

	// first we do the regular keptn values
	result = strings.Replace(result, "$CONTEXT", k.KeptnContext, -1)
	result = strings.Replace(result, "$PROJECT", k.Event.GetProject(), -1)
	result = strings.Replace(result, "$STAGE", k.Event.GetStage(), -1)
	result = strings.Replace(result, "$SERVICE", k.Event.GetService(), -1)
	if k.Event.(KeptnBaseEvent).DeploymentStrategy != nil {
		result = strings.Replace(result, "$DEPLOYMENT", *k.Event.(KeptnBaseEvent).DeploymentStrategy, -1)
	}
	if k.Event.(KeptnBaseEvent).TestStrategy != nil {
		result = strings.Replace(result, "$TESTSTRATEGY", *k.Event.(KeptnBaseEvent).TestStrategy, -1)
	}

	// now we do the labels
	for key, value := range k.Event.GetLabels() {
		result = strings.Replace(result, "$LABEL."+key, value, -1)
	}

	// now we do all environment variables
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		result = strings.Replace(result, "$ENV."+pair[0], pair[1], -1)
	}

	return result
}
