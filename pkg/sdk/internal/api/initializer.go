package api

import (
	"fmt"
	"github.com/benbjohnson/clock"
	keptnapi "github.com/keptn/go-utils/pkg/api/utils"
	keptnapiv2 "github.com/keptn/go-utils/pkg/api/utils/v2"
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
	KeptnAPIV2          keptnapiv2.KeptnInterface
	ControlPlane        *controlplane.ControlPlane
	EventSenderCallback controlplane.EventSender
}

// Initialize takes care of creating the API clients and initializing the cp-connector library based
// on environment variables
func Initialize(env config.EnvConfig, clientFactory HTTPClientGetter, logger logger.Logger) (*InitializationResult, error) {
	// initialize http client
	httpClient, err := clientFactory.Get()
	if err != nil {
		return nil, fmt.Errorf("could not initialize HTTP client: %w", err)
	}
	// fall back to uninitialized http client
	if httpClient == nil {
		httpClient = &http.Client{}
	}

	// initialize api
	api, err := apiSet(env, httpClient)
	if err != nil {
		return nil, fmt.Errorf("could not initialize control plane client api: %w", err)
	}

	apiV2, err := apiSetV2(env, httpClient)
	if err != nil {
		return nil, fmt.Errorf("could not initialize v2 control plane client api: %w", err)
	}

	// initialize api handlers and cp-connector components
	ss, es, lf := createCPComponents(api, logger, env)
	controlPlane := controlplane.New(ss, es, lf, controlplane.WithLogger(logger))

	return &InitializationResult{
		KeptnAPI:            api,
		KeptnAPIV2:          apiV2,
		ControlPlane:        controlPlane,
		EventSenderCallback: es.Sender(),
	}, nil

}

func apiSet(env config.EnvConfig, httpClient *http.Client) (keptnapi.KeptnInterface, error) {
	if env.PubSubConnectionType() == config.ConnectionTypeHTTP {
		scheme, err := getHttpScheme(env)
		if err != nil {
			return nil, err
		}
		return keptnapi.New(env.KeptnAPIEndpoint, keptnapi.WithScheme(scheme), keptnapi.WithHTTPClient(httpClient), keptnapi.WithAuthToken(env.KeptnAPIToken))

	}
	return keptnapi.NewInternal(httpClient)

}

func apiSetV2(env config.EnvConfig, httpClient *http.Client) (keptnapiv2.KeptnInterface, error) {
	if env.PubSubConnectionType() == config.ConnectionTypeHTTP {
		scheme, err := getHttpScheme(env)
		if err != nil {
			return nil, err
		}
		return keptnapiv2.New(env.KeptnAPIEndpoint, keptnapiv2.WithScheme(scheme), keptnapiv2.WithHTTPClient(httpClient), keptnapiv2.WithAuthToken(env.KeptnAPIToken))

	}
	return keptnapiv2.NewInternal(httpClient)

}

func getHttpScheme(env config.EnvConfig) (string, error) {
	parsed, err := url.ParseRequestURI(env.KeptnAPIEndpoint)
	if err != nil {
		return "", fmt.Errorf("could not parse given Keptn API endpoint: %w", err)
	}

	if !strings.HasPrefix(parsed.Scheme, "http") {
		return "", fmt.Errorf("invalid scheme for keptn endpoint, %s is not http or https", env.KeptnAPIEndpoint)
	}

	return parsed.Scheme, nil
}

func eventSource(apiSet keptnapi.KeptnInterface, logger logger.Logger, env config.EnvConfig) eventsource.EventSource {
	if env.PubSubConnectionType() == config.ConnectionTypeHTTP {
		return eventsourceHttp.New(clock.New(), eventsourceHttp.NewEventAPI(apiSet.ShipyardControlV1(), apiSet.APIV1()), eventsourceHttp.WithLogger(logger))
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
