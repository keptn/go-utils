package keptn

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	cloudevents "github.com/cloudevents/sdk-go/v2"

	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"gopkg.in/yaml.v3"
)

type KeptnOpts struct {
	UseLocalFileSystem      bool
	ConfigurationServiceURL string
	EventBrokerURL          string // Deprecated: use EventSender instead
	DatastoreURL            string
	IncomingEvent           *cloudevents.Event
	LoggingOptions          *LoggingOpts
	EventSender             EventSender
	Context                 context.Context
}

type LoggingOpts struct {
	EnableWebsocket   bool
	WebsocketEndpoint *string
	ServiceName       *string
}

type KeptnBase struct {
	KeptnContext string

	Event      EventProperties
	CloudEvent *cloudevents.Event
	Logger     LoggerInterface

	// EventSender object that is responsible for sending events
	EventSender EventSender

	EventBrokerURL     string // Deprecated: use EventSender instead
	UseLocalFileSystem bool
	ResourceHandler    *api.ResourceHandler
	EventHandler       *api.EventHandler

	// Context is used during the lifetime of the Keptn instance
	// in order to propagate cross-cutting concerns, such as tracecontext and cancellation
	Context context.Context
}

type EventProperties interface {
	GetProject() string
	GetStage() string
	GetService() string
	GetLabels() map[string]string
	SetProject(string)
	SetStage(string)
	SetService(string)
	SetLabels(map[string]string)
}

type EventSender interface {
	// SendEvent sends a cloud event
	// Deprecated: use EventSender.Send instead
	SendEvent(event cloudevents.Event) error

	// Send sends a cloud event
	Send(ctx context.Context, event cloudevents.Event) error
}

// SLIConfig represents the struct of a SLI file
type SLIConfig struct {
	Indicators map[string]string `json:"indicators" yaml:"indicators"`
}

const ConfigurationServiceURL = "configuration-service:8080"
const DatastoreURL = "mongodb-datastore:8080"
const DefaultLoggingServiceName = "keptn"

// GetSLIConfiguration retrieves the SLI configuration for a service considering SLI configuration on stage and project level.
// First, the configuration of project-level is retrieved, which is then overridden by configuration on stage level,
// overridden by configuration on service level.
func (k *KeptnBase) GetSLIConfiguration(project string, stage string, service string, resourceURI string) (map[string]string, error) {
	var res *models.Resource
	var err error
	SLIs := make(map[string]string)

	// get sli config from project
	if project != "" {
		res, err = k.ResourceHandler.GetProjectResource(project, resourceURI)
		if err != nil {
			// return error except "resource not found" type
			if !strings.Contains(strings.ToLower(err.Error()), "resource not found") {
				return nil, err
			}
		}
		SLIs, err = addResourceContentToSLIMap(SLIs, res)
		if err != nil {
			return nil, err
		}
	}

	// get sli config from stage
	if project != "" && stage != "" {
		res, err = k.ResourceHandler.GetStageResource(project, stage, resourceURI)
		if err != nil {
			// return error except "resource not found" type
			if !strings.Contains(strings.ToLower(err.Error()), "resource not found") {
				return nil, err
			}
		}
		SLIs, err = addResourceContentToSLIMap(SLIs, res)
		if err != nil {
			return nil, err
		}
	}

	// get sli config from service
	if project != "" && stage != "" && service != "" {
		res, err = k.ResourceHandler.GetServiceResource(project, stage, service, resourceURI)
		if err != nil {
			// return error except "resource not found" type
			if !strings.Contains(strings.ToLower(err.Error()), "resource not found") {
				return nil, err
			}
		}
		SLIs, err = addResourceContentToSLIMap(SLIs, res)
		if err != nil {
			return nil, err
		}
	}

	return SLIs, nil
}

func addResourceContentToSLIMap(SLIs map[string]string, resource *models.Resource) (map[string]string, error) {
	if resource != nil {
		sliConfig := SLIConfig{}
		err := yaml.Unmarshal([]byte(resource.ResourceContent), &sliConfig)
		if err != nil {
			return nil, err
		}

		for key, value := range sliConfig.Indicators {
			SLIs[key] = value
		}

		if len(SLIs) == 0 {
			return nil, errors.New("missing required field: indicators")
		}
	}
	return SLIs, nil
}

// GetKeptnResource returns a resource from the configuration repo based on the incoming cloud events project, service and stage
func (k *KeptnBase) GetKeptnResource(resource string) ([]byte, error) {

	// if we run in a runlocal mode we are just getting the file from the local disk
	if k.UseLocalFileSystem {
		return getKeptnResourceFromLocal(resource)
	}

	// get it from KeptnBase
	requestedResource, err := k.ResourceHandler.GetServiceResource(k.Event.GetProject(), k.Event.GetStage(), k.Event.GetService(), resource)

	// return Nil in case resource couldn't be retrieved
	if err != nil || requestedResource.ResourceContent == "" {
		return nil, fmt.Errorf("resource not found: %s - %s", resource, err)
	}

	return []byte(requestedResource.ResourceContent), nil
}

/**
 * Retrieves a resource (=file) from the local file system. Basically checks if the file is available and if so returns it
 */
func getKeptnResourceFromLocal(resource string) ([]byte, error) {
	if _, err := os.Stat(resource); err == nil {
		return []byte(resource), nil
	} else {
		return nil, err
	}
}

// ValidateKeptnEntityName checks whether the provided name represents a valid
// project, service, or stage name
func ValidateKeptnEntityName(name string) bool {
	if len(name) == 0 {
		return false
	}
	reg, err := regexp.Compile(`(^[a-z][a-z0-9-]*[a-z0-9]$)|(^[a-z][a-z0-9]*)`)
	if err != nil {
		log.Fatal(err)
	}
	processedString := reg.FindString(name)
	return len(processedString) == len(name)
}

// ValidateUnixDirectoryName checks whether the provided dirName contains
// any special character according to
// https://www.cyberciti.biz/faq/linuxunix-rules-for-naming-file-and-directory-names/
func ValidateUnixDirectoryName(dirName string) bool {
	return !(dirName == "." ||
		dirName == ".." ||
		strings.Contains(dirName, "/") ||
		strings.Contains(dirName, ">") ||
		strings.Contains(dirName, "<") ||
		strings.Contains(dirName, "|") ||
		strings.Contains(dirName, ":") ||
		strings.Contains(dirName, "&"))
}

// getServiceEndpoint gets an endpoint stored in an environment variable and sets http as default scheme
func GetServiceEndpoint(service string) (url.URL, error) {
	envVal := os.Getenv(service)
	if envVal == "" {
		return url.URL{}, fmt.Errorf("Provided environment variable %s has no valid value", service)
	}

	url, err := url.Parse(envVal)
	if err != nil {
		return *url, fmt.Errorf("Failed to retrieve value from ENVIRONMENT_VARIABLE: %s", service)
	}

	if url.Scheme == "" {
		url.Scheme = "http"
	}

	// check if only a service name has been provided, e.g. 'configuration-service'
	if url.Host == "" && url.Path != "" {
		url.Host = url.Path
		url.Path = ""
	}

	return *url, nil
}

func GetExpBackoffTime(retryNr int) time.Duration {
	f := 1.5 * float64(retryNr)
	if retryNr <= 1 {
		f = 1.5
	}
	currentInterval := float64(500*time.Millisecond) * f
	randomizationFactor := 0.5
	random := rand.Float64()

	var delta = randomizationFactor * currentInterval
	minInterval := float64(currentInterval) - delta
	maxInterval := float64(currentInterval) + delta

	return time.Duration(minInterval + (random * (maxInterval - minInterval + 1)))
}
