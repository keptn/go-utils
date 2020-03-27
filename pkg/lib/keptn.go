package keptn

import (
	"encoding/json"
	"fmt"
	"github.com/cloudevents/sdk-go/pkg/cloudevents"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"regexp"
	"strings"
)

type KeptnOpts struct {
	UseLocalFileSystem      bool
	ConfigurationServiceURL string
	EventBrokerURL          string
	IncomingEvent           *cloudevents.Event
}

type Keptn struct {
	KeptnContext string

	KeptnBase *KeptnBase

	eventBrokerURL     string
	useLocalFileSystem bool
	resourceHandler    *ResourceHandler
}

// Shipyard defines the name, deployment strategy and test strategy of each stage
type Shipyard struct {
	Stages []struct {
		Name                string `json:"name" yaml:"name"`
		DeploymentStrategy  string `json:"deployment_strategy" yaml:"deployment_strategy"`
		TestStrategy        string `json:"test_strategy,omitempty" yaml:"test_strategy"`
		RemediationStrategy string `json:"remediation_strategy,omitempty" yaml:"remediation_strategy"`
	} `json:"stages" yaml:"stages"`
}

const configurationServiceURL = "configuration-service:8080"
const defaultEventBrokerURL = "http://event-broker.keptn.svc.cluster.local/keptn"

func NewKeptn(incomingEvent *cloudevents.Event, opts KeptnOpts) (*Keptn, error) {
	shkeptncontext, _ := incomingEvent.Context.GetExtension("shkeptncontext")

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
		KeptnContext:       shkeptncontext.(string),
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

	k.resourceHandler = NewResourceHandler(csURL)

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
