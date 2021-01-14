# Keptn go-utils
![GitHub release (latest by date)](https://img.shields.io/github/v/release/keptn/go-utils)
![tests](https://github.com/keptn/go-utils/workflows/tests/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/keptn/go-utils)](https://goreportcard.com/report/github.com/keptn/go-utils)

This repo serves as a util package for common functionalities such as logging of the [Keptn Project](https://github.com/keptn).

Please post any issues with this package to the [keptn/keptn repository](https://github.com/keptn/keptn/issues) and label them with `area:go-utils`.

## Installation

Get the latest version using
```console
go get github.com/keptn/go-utils
```
Also consider browsing our [Releases Page](https://github.com/keptn/go-utils/releases) to find out about all releases.


## Contributing

If you want to contribute, just create a PR on the **master** branch.

Please also see [CONTRIBUTING.md](CONTRIBUTING.md) instructions on how to contribute.

## Usage

Below we have listed some basic examples. You can find more information about the usage within the [docs/](docs/) folder.

### Utils
If you need to access several utility functions:

```go
import {
  "github.com/keptn/go-utils/pkg/lib"
}
```

This module provides you with a convenient `Keptn` helper struct that allows you to access several resources that are relevant within the context of a 
Keptn event. The helper struct can be initialized by passing a CloudEvent to the `NewKeptn` function. Example:

```go
func HandleEvent(event cloudevents.Event) error {
	keptnHandler, err := keptn.NewKeptn(&event, keptn.KeptnOpts{})
	if err != nil {
		return nil, err
	}
	
    // get the shipyard file of the project
    shipyard, _ := keptnHandler.GetShipyard()
    
    // get a resource within the current context (i.e., project, stage, service) of the event
    resourceContent, _ := keptnHandler.GetKeptnResource("resource.yaml")

    // send a cloud event
    _ = keptnHandler.SendCloudEvent(event)
    // ...
}
```

### CloudEvent Data
If you need to access data within CloudEvents:

```go
import {
	"github.com/keptn/go-utils/pkg/lib"
)
```

Example:

```go
func parseCloudEvent(event cloudevents.Event) (keptnevents.TestFinishedEventData, error) {
	eventData := &keptn.TestsFinishedEventData{}
	err := event.DataAs(eventData)
    
    return eventData, err
}
```

### Models
If you need to access Models for YAML files:

```go
import {
	"github.com/keptn/go-utils/pkg/lib"
)
```

### Querying Events from event store
```
// Create an event Handler
eventHandler := apiutils.NewAuthenticatedEventHandler("1.2.3.4/api", token, "x-token", nil, "http")

// Create a filter
filter := &apiutils.EventFilter{KeptnContext:  "03b2b951-9835-4e87-b0b0-0ad0bc288214"}

// Query for event(s)
events, err := eventHandler.GetEvents(filter)
```

### Watching for events in event store
```
// Create a watcher
watcher := api.NewEventWatcher(eventhandler),
	api.WithEventFilter(api.EventFilter{  // use custom filter
           Project: "sockshop",
           KeptnContext: "..."}),
	api.WithEventManipulator(api.SortByTime),         // apply custom logic to fetched events
	api.WithInterval(time.NewTicker(5*time.Second)), // fetch every 5 seconds
	api.WithStartTime(time.Now()),                  // fetch events new than this timestamp
	api.WithTimeout(time.Second * 15),             // stop fetching events after 15 secs
)

    // start watcher and consume events
	allEvents, _ := watcher.Watch(context.Background())
	for events := range allEvents {
		for _, e := range events {
			fmt.Println(*e.Type)
		}
	}
``` 


## Updating model definitions
After updating the model definitions in the `swagger.yaml` file, execute the command

```
swagger generate model --spec=swagger.yaml -t=./pkg/api/
```

to update the models located in `./pkg/api/models`


## Automation

A [Github Action](https://github.com/keptn/go-utils/actions?query=workflow%3A%22Auto+PR+to+keptn%2Fkeptn%22) is used
that creates a Pull Request to  [github.com/keptn/keptn](https://github.com/keptn/keptn) to update `go.mod`
files with an updated version of this  package (based on the commit hash).