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
func HandleEvent(ctx context.Context, event cloudevents.Event) error {
    keptnHandler, err := keptn.NewKeptn(&event, keptn.KeptnOpts{
		Context: ctx,
	})
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

By default, the `SendCloudEvent` function of the `Keptn` struct will send events to the distributor sidecar that is running within the same pod of the Keptn service (see [the Keptn doc](https://keptn.sh/docs/0.8.x/integrations/custom_integration/#subscription-to-keptn-event) for more details).
This behavior can be overridden by passing an implementation of the `EventSender` interface via the `KeptnOpts` object that is passed to the `NewKeptn` function, e.g.:

```go
// custom EventSender

type MyCustomEventSender struct {
}

func (es *MyCustomEventSender) SendEvent(event cloudevents.Event) error {
    // custom implementation
    return nil
}
//...
keptnHandler, err := keptn.NewKeptn(&event, keptn.KeptnOpts{
    EventSender: &MyCustomEventSender{}
})
if err != nil {
    return nil, err
}
```

For unit testing purposes, we offer a mock implementation of the `EventSender` interface in the `keptn/v0_2_0/fake` package. This mock implementation can be used to check if your service sends the expected events. E.g.:

```go
func MyTest(t *testing.T) {
    fakeSender := &keptnfake.EventSender{}

    // optionally, you can add custom behavior for certain event types (e.g. returning an error for a certain event type):
    fakeSender.AddReactor("sh.keptn.event.deployment.finished", func(event cloudevents.Event) error {
        return errors.New("i throw an error if i should send a 'sh.keptn.event.deployment.finished' event")
    })

    keptnHandler, err := keptn.NewKeptn(&event, keptn.KeptnOpts{
        EventSender: fakeSender
    })

    myService := &MyKeptnService{
        KeptnHandler: keptnHandler
    }
    
    myService.handleEvent(myReceivedCloudEvent)
    
    // check if your service sends out events of the type "sh.keptn.event.deployment.started" and "sh.keptn.event.deployment.finished"
    if err := fakeSender.AssertSentEventTypes([]string{"sh.keptn.event.deployment.started", "sh.keptn.event.deployment.finished"}); err != nil {
        t.Errorf("%s", err.Error())
    }

    // to inspect the sent events in more detail, you can access them via fakeSender.SentEvents
    for _, event := range fakeSender.SentEvents {
        // do some validation here
    }

    if err != nil {
        return nil, err
    }
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
events, err := eventHandler.GetEventsWithContext(context.Background(), filter)
```

### Watching for events in event store
```
// Create a watcher
watcher := api.NewEventWatcher(eventhandler),
	api.WithEventFilter(api.EventFilter{  // use custom filter
           Project: "sockshop",
           KeptnContext: "..."}),
	api.WithInterval(time.NewTicker(5*time.Second)), // fetch every 5 seconds
	api.WithStartTime(time.Now()),                  // start fetching events newer than this timestamp
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


## Automation

A [GitHub Action](https://github.com/keptn/go-utils/actions?query=workflow%3A%22Auto+PR+to+keptn%2Fkeptn%22) is used
that creates a Pull Request to  [github.com/keptn/keptn](https://github.com/keptn/keptn) to update `go.mod`
files with an updated version of this  package (based on the commit hash).