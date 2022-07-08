package main

import (
	"context"
	"github.com/benbjohnson/clock"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/sdk/connector/controlplane"
	"github.com/keptn/go-utils/pkg/sdk/connector/eventsource/http"
	"github.com/keptn/go-utils/pkg/sdk/connector/logforwarder"
	"github.com/keptn/go-utils/pkg/sdk/connector/subscriptionsource"
	"github.com/keptn/go-utils/pkg/sdk/connector/types"
	"github.com/sirupsen/logrus"
	"log"
	"time"
)

const (
	Endpoint = ""
	Token    = ""
)

func main() {
	if Endpoint == "" || Token == "" {
		log.Fatal("Please set Keptn API endpoint and API Token to use this example")
	}

	// Optional: Create your favorite logger (e.g. logrus)
	logger := logrus.New()
	logger.SetLevel(logrus.DebugLevel)

	// Create an APISet to be used to talk to the Keptn API
	// Provide it with the endpoint (e.g. http://my-keptn.nip.io/api")
	// and an API token
	keptnAPI, err := api.New(Endpoint, api.WithAuthToken(Token))
	if err != nil {
		log.Fatal(err)
	}

	// Create a subscription source that is responsible for registering your integration and managing
	// subscriptions
	subscriptionSource := subscriptionsource.New(keptnAPI.UniformV1(), subscriptionsource.WithLogger(logger))

	// Create an event source that is responsible for getting events from keptn
	eventSource := http.New(clock.New(), http.NewEventAPI(keptnAPI.ShipyardControlV1(), keptnAPI.APIV1()))

	// Optional: create a log forwarder that is responsible for sending error log events to Keptn
	// If you don't want/need that, you can simply pass nil
	logForwarder := logforwarder.New(keptnAPI.LogsV1())

	// Create a control plane component that is the main component of cp-connector and start it
	// using RunWithGraceFulShutdown
	controlPlane := controlplane.New(subscriptionSource, eventSource, logForwarder, controlplane.WithLogger(logger))
	if err := controlplane.RunWithGracefulShutdown(controlPlane, LocalService{}, time.Second*10); err != nil {
		log.Fatal(err)
	}
}

// LocalService is an implementation of the Integration interface
// and represents the entry point of your event handling logic
type LocalService struct{}

// OnEvent is called for every event that was received.
// This is the place to insert your event processing business logic.
//
// Note, that you are responsible for sending corresponding .started and .finished events
// on your own.
// Also note, that if you need to ensure that every incoming event is completely processed before the pod running your
// integration is shut down (e.g., due to an upgrade to a newer version), the OnEvent method should process the incoming events synchronously,
// i.e. not in a separate go routine. If you need to process events asynchronously, you need to implement your own synchronization mechanism to ensure all
// events have been completely processed before a shutdown
func (e LocalService) OnEvent(ctx context.Context, event models.KeptnContextExtendedCE) error {
	// You can grab handle the event and grab a sender to send back started / finished events to keptn
	// eventSender := ctx.Value(controlplane.EventSenderKeyType{}).(types.EventSender)
	return nil
}

// RegistrationData is used for initial registration to the Keptn control plane.
// usually this information is set with information coming from environment variables when the pod is started.
// In this example everything is hard-coded
func (e LocalService) RegistrationData() types.RegistrationData {
	return types.RegistrationData{
		// Name is the name of your service and is visible on the integrations page of the Keptn bridge
		Name: "local-service",
		MetaData: models.MetaData{
			// Hostname usually takes the value of the kubernetes node name
			Hostname: "localhost",
			// IntegrationVersion is the version of your service/integration that
			// is displayed on the integration page of the Keptn bridge
			IntegrationVersion: "dev",
			// For legacy reasons you must provide a distributor version.
			// The value does not really have an effect but must be something greater
			// than 0.9.0.
			DistributorVersion: "0.15.0",
			// Location is a hint whether this service is running on the control plane or remotely as a remote
			// execution plane service. Usually the values are just: "control-plane" or "remote-execution-plane"
			Location: "local",
			// KubernetesMetaData is important information used to register your service to the control plane
			KubernetesMetaData: models.KubernetesMetaData{
				// Namespace the service is running in
				Namespace: "keptn",
				// PodName is the K8S pod name
				PodName: "my-pod",
				// DeploymentName is the K8S deployment name
				DeploymentName: "my-deployment",
			},
		},
		// subscribe to sh.keptn.event.echo.triggered events, with no filter
		Subscriptions: []models.EventSubscription{
			{
				// Event represents the event type you want to subscribe to
				Event: "sh.keptn.event.echo.triggered",
				// Filter represents additional filters for your subscription
				Filter: models.EventSubscriptionFilter{},
			},
		},
	}
}
