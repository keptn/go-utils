package keptn

import (
	"encoding/json"
	"fmt"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"github.com/keptn/go-utils/pkg/api/models"
	api "github.com/keptn/go-utils/pkg/api/utils"
	"gopkg.in/yaml.v2"
	"log"
	"net/url"
	"os"
	"regexp"
	"strings"
)

type KeptnOpts struct {
	UseLocalFileSystem      bool
	ConfigurationServiceURL string
	EventBrokerURL          string
	IncomingEvent           *cloudevents.Event
	LoggingOptions          *LoggingOpts
}

type LoggingOpts struct {
	EnableWebsocket   bool
	WebsocketEndpoint *string
	ServiceName       *string
}

type Keptn struct {
	KeptnContext string

	KeptnBase *KeptnBase

	Logger LoggerInterface

	eventBrokerURL     string
	useLocalFileSystem bool
	resourceHandler    *api.ResourceHandler
	eventHandler       *api.EventHandler
}

// SLIConfig represents the struct of a SLI file
type SLIConfig struct {
	Indicators map[string]string `json:"indicators" yaml:"indicators"`
}

const configurationServiceURL = "configuration-service:8080"
const defaultEventBrokerURL = "http://event-broker.keptn.svc.cluster.local/keptn"
const defaultWebsocketEndpoint = "ws://api-service.keptn.svc.cluster.local:8080"
const defaultLoggingServiceName = "keptn"

func NewKeptn(incomingEvent *cloudevents.Event, opts KeptnOpts) (*Keptn, error) {
	var shkeptncontext string
	_ = incomingEvent.Context.ExtensionAs("shkeptncontext", &shkeptncontext)

	// create a base Keptn Event
	keptnBase := &KeptnBase{}

	bytes, err := incomingEvent.DataBytes()
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, keptnBase)
	if err != nil {
		return nil, err
	}

	k := &Keptn{
		KeptnBase:          keptnBase,
		KeptnContext:       shkeptncontext,
		useLocalFileSystem: opts.UseLocalFileSystem,
		resourceHandler:    nil,
	}
	csURL := configurationServiceURL
	if opts.ConfigurationServiceURL != "" {
		csURL = opts.ConfigurationServiceURL
	}

	if opts.EventBrokerURL != "" {
		k.eventBrokerURL = opts.EventBrokerURL
	} else {
		k.eventBrokerURL = defaultEventBrokerURL
	}

	k.resourceHandler = api.NewResourceHandler(csURL)
	k.eventHandler = api.NewEventHandler(csURL)

	loggingServiceName := defaultLoggingServiceName
	if opts.LoggingOptions != nil && opts.LoggingOptions.ServiceName != nil {
		loggingServiceName = *opts.LoggingOptions.ServiceName
	}
	k.Logger = NewLogger(k.KeptnContext, incomingEvent.Context.GetID(), loggingServiceName)

	if opts.LoggingOptions != nil && opts.LoggingOptions.EnableWebsocket {
		wsURL := defaultWebsocketEndpoint
		if opts.LoggingOptions.WebsocketEndpoint != nil && *opts.LoggingOptions.WebsocketEndpoint != "" {
			wsURL = *opts.LoggingOptions.WebsocketEndpoint
		}
		connData := ConnectionData{}
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
			ws, _, err := OpenWS(connData, *apiServiceURL)
			if err != nil {
				k.Logger.Error("Opening WebSocket connection failed:" + err.Error())
				return k, nil
			}
			stdLogger := NewLogger(shkeptncontext, incomingEvent.Context.GetID(), loggingServiceName)
			combinedLogger := NewCombinedLogger(stdLogger, ws, shkeptncontext)
			k.Logger = combinedLogger
		}
	}

	return k, nil
}

// GetShipyard returns the shipyard definition of a project
func (k *Keptn) GetShipyard() (*Shipyard, error) {
	shipyardResource, err := k.resourceHandler.GetProjectResource(k.KeptnBase.Project, "shipyard.yaml")
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

// GetSLIConfiguration retrieves the SLI configuration for a service considering SLI configuration on stage and project level.
// First, the configuration of project-level is retrieved, which is then overridden by configuration on stage level,
// overridden by configuration on service level.
func (k *Keptn) GetSLIConfiguration(project string, stage string, service string, resourceURI string) (map[string]string, error) {
	var res *models.Resource
	var err error
	SLIs := make(map[string]string)

	// get sli config from project
	if project != "" {
		res, err = k.resourceHandler.GetProjectResource(project, resourceURI)
		if err != nil {
			// return error except "resource not found" type
			if !strings.Contains(err.Error(), "resource not found") {
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
		res, err = k.resourceHandler.GetStageResource(project, stage, resourceURI)
		if err != nil {
			// return error except "resource not found" type
			if !strings.Contains(err.Error(), "resource not found") {
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
		res, err = k.resourceHandler.GetServiceResource(project, stage, service, resourceURI)
		if err != nil {
			// return error except "resource not found" type
			if !strings.Contains(err.Error(), "resource not found") {
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
	}
	return SLIs, nil
}

func (k *Keptn) GetKeptnResource(resource string) (string, error) {

	// if we run in a runlocal mode we are just getting the file from the local disk
	if k.useLocalFileSystem {
		return _getKeptnResourceFromLocal(resource)
	}

	// get it from Keptn
	requestedResource, err := k.resourceHandler.GetServiceResource(k.KeptnBase.Project, k.KeptnBase.Stage, k.KeptnBase.Service, resource)

	// return Nil in case resource couldnt be retrieved
	if err != nil || requestedResource.ResourceContent == "" {
		fmt.Printf("Keptn Resource not found: %s - %s", resource, err)
		return "", err
	}

	// now store that file on the same directory structure locally
	os.RemoveAll(resource)
	pathArr := strings.Split(resource, "/")
	directory := ""
	for _, pathItem := range pathArr[0 : len(pathArr)-1] {
		directory += pathItem + "/"
	}

	if directory != "" {
		err = os.MkdirAll(directory, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	resourceFile, err := os.Create(resource)
	if err != nil {
		fmt.Errorf(err.Error())
		return "", err
	}
	defer resourceFile.Close()

	_, err = resourceFile.Write([]byte(requestedResource.ResourceContent))

	if err != nil {
		fmt.Errorf(err.Error())
		return "", err
	}

	return strings.TrimSuffix(requestedResource.ResourceContent, "\n"), nil
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
	result = strings.Replace(result, "$PROJECT", k.KeptnBase.Project, -1)
	result = strings.Replace(result, "$STAGE", k.KeptnBase.Stage, -1)
	result = strings.Replace(result, "$SERVICE", k.KeptnBase.Service, -1)
	if k.KeptnBase.DeploymentStrategy != nil {
		result = strings.Replace(result, "$DEPLOYMENT", *k.KeptnBase.DeploymentStrategy, -1)
	}
	if k.KeptnBase.TestStrategy != nil {
		result = strings.Replace(result, "$TESTSTRATEGY", *k.KeptnBase.TestStrategy, -1)
	}

	// now we do the labels
	for key, value := range k.KeptnBase.Labels {
		result = strings.Replace(result, "$LABEL."+key, value, -1)
	}

	// now we do all environment variables
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		result = strings.Replace(result, "$ENV."+pair[0], pair[1], -1)
	}

	return result
}

/**
 * Retrieves a resource (=file) from the local file system. Basically checks if the file is available and if so returns it
 */
func _getKeptnResourceFromLocal(resource string) (string, error) {
	if _, err := os.Stat(resource); err == nil {
		return resource, nil
	} else {
		return "", err
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

// ValididateUnixDirectoryName checks whether the provided dirName contains
// any special character according to
// https://www.cyberciti.biz/faq/linuxunix-rules-for-naming-file-and-directory-names/
func ValididateUnixDirectoryName(dirName string) bool {
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
	url, err := url.Parse(os.Getenv(service))
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
