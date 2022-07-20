package api

import (
	"fmt"
	"github.com/benbjohnson/clock"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	"github.com/keptn/go-utils/pkg/sdk/connector/controlplane"
	"github.com/keptn/go-utils/pkg/sdk/connector/eventsource"
	eventsourceHttp "github.com/keptn/go-utils/pkg/sdk/connector/eventsource/http"
	eventsourceNats "github.com/keptn/go-utils/pkg/sdk/connector/eventsource/nats"
	"github.com/keptn/go-utils/pkg/sdk/connector/logforwarder"
	"github.com/keptn/go-utils/pkg/sdk/connector/logger"
	"github.com/keptn/go-utils/pkg/sdk/connector/nats"
	"github.com/keptn/go-utils/pkg/sdk/connector/subscriptionsource"
	"github.com/keptn/go-utils/pkg/sdk/internal/config"
	"net/http"
	"net/url"
	"strings"
)

type InitializationResult struct {
	KeptnAPI            keptnapi.KeptnInterface
	ControlPlane        *controlplane.ControlPlane
	EventSenderCallback controlplane.EventSender
	ResourceHandler     *keptnapi.ResourceHandler
	Error               error
}

// Initialize takes care of creating the API clients and initializing the cp-connector library based
// on environment variables
func Initialize(env config.EnvConfig, clientFactory HTTPClientGetter, logger logger.Logger) *InitializationResult {
	// initialize http client
	httpClient, err := clientFactory.Get()
	if err != nil {
		return &InitializationResult{
			Error: fmt.Errorf("could not initialize HTTP client: %w", err),
		}
	}
	// fall back to uninitialized http client
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	// initialize api set
	apiSet, err := apiSet(env, httpClient)
	if err != nil {
		return &InitializationResult{
			Error: fmt.Errorf("could not initialize control plane client api: %w", err),
		}
	}
	// initialize api handlers and cp-connector components
	resourceHandler := resourceHandler(env)
	ss, es, lf := createCPComponents(apiSet, logger, env)
	controlPlane := controlplane.New(ss, es, lf, controlplane.WithLogger(logger))

	return &InitializationResult{
		KeptnAPI:            apiSet,
		ControlPlane:        controlPlane,
		EventSenderCallback: es.Sender(),
		ResourceHandler:     resourceHandler,
	}

}

func apiSet(env config.EnvConfig, httpClient *http.Client) (keptnapi.KeptnInterface, error) {

	if env.PubSubConnectionType() == config.ConnectionTypeHTTP {
		scheme := "http"
		parsed, err := url.ParseRequestURI(env.KeptnAPIEndpoint)
		if err != nil {
			return nil, fmt.Errorf("could not parse given Keptn API endpoint: %w", err)
		}

		if parsed.Scheme == "" || !strings.HasPrefix(parsed.Scheme, "http") {
			return nil, fmt.Errorf("invalid scheme for keptn endpoint, %s is not http or https", env.KeptnAPIEndpoint)
		}

		if strings.HasPrefix(parsed.Scheme, "http") {
			scheme = parsed.Scheme
		}
		return keptnapi.New(env.KeptnAPIEndpoint, keptnapi.WithScheme(scheme), keptnapi.WithHTTPClient(httpClient), keptnapi.WithAuthToken(env.KeptnAPIToken))

	}
	return keptnapi.NewInternal(httpClient)

}

func eventSource(apiSet keptnapi.KeptnInterface, logger logger.Logger, env config.EnvConfig) eventsource.EventSource {
	if env.PubSubConnectionType() == config.ConnectionTypeHTTP {
		return eventsourceHttp.New(clock.New(), eventsourceHttp.NewEventAPI(apiSet.ShipyardControlV1(), apiSet.APIV1()))
	}
	natsConnector := nats.New(env.EventBrokerURL, nats.WithLogger(logger))
	return eventsourceNats.New(natsConnector, eventsourceNats.WithLogger(logger))
}

func subscriptionSource(apiSet keptnapi.KeptnInterface, logger logger.Logger) subscriptionsource.SubscriptionSource {
	return subscriptionsource.New(apiSet.UniformV1(), subscriptionsource.WithLogger(logger))
}

func logForwarder(apiSet keptnapi.KeptnInterface, logger logger.Logger) logforwarder.LogForwarder {
	return logforwarder.New(apiSet.LogsV1(), logforwarder.WithLogger(logger))
}

func createCPComponents(apiSet keptnapi.KeptnInterface, logger logger.Logger, env config.EnvConfig) (subscriptionsource.SubscriptionSource, eventsource.EventSource, logforwarder.LogForwarder) {
	return subscriptionSource(apiSet, logger), eventSource(apiSet, logger, env), logForwarder(apiSet, logger)
}
func resourceHandler(env config.EnvConfig) *keptnapi.ResourceHandler {
	return keptnapi.NewResourceHandler(env.ConfigurationServiceURL)
}
