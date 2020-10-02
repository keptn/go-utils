# Keptn go-utils
[![Build Status](https://travis-ci.org/keptn/go-utils.svg?branch=master)](https://travis-ci.org/keptn/go-utils)
![GitHub release (latest by date)](https://img.shields.io/github/v/release/keptn/go-utils)
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

## Updating model definitions
After updating the model definitions in the `swagger.yaml` file, execute the command

```
swagger generate model --spec=swagger.yaml -t=./pkg/api/
```

to update the models located in `./pkg/api/models`

## Automation

Within [.travis.yml](.travis.yml) we have included an automation that creates a Pull Request to 
 [github.com/keptn/keptn](https://github.com/keptn/keptn) to update `go.mod` files with an updated version of this 
 package (based on the commit hash). To make this work, a `GITHUB_TOKEN` (personal access token) 
 needs to be added within the [travis-ci settings page](https://travis-ci.org/keptn/go-utils/settings).
 
## Upgrade to 0.7.2 from previous versions
This version introduces a couple of changes within the structure of the module. When upgrading from an earlier version, please follow the following steps:

The following exported types/funcs that have been imported from `github.com/keptn/go-utils/pkg/lib` have been moved to `github.com/keptn/go-utils/pkg/lib/keptn`
 
- `KeptnOpts`
- `LoggingOpts`
- `KeptnBase`
- `EventProperties`
- `SLIConfig`
- `CombinedLogger`
- `NewCombinedLogger()`
- `NewLogger()`
- `Logger`
- `MyCloudEvent`
- `LogData`
- `IncompleteCE`
- `ConnectionData`
- `OpenWS()`
- `WriteWSLog()`
- `WriteLog()`
- `LoggerInterface`
- `ValidateKeptnEntityName()`
- `ValididateUnixDirectoryName()`
- `GetServiceEndpoint()`

If you have used any of those, you will need to change the import from 

```go
import github.com/keptn/go-utils/pkg/lib
```

to 

```go
import github.com/keptn/go-utils/pkg/lib/keptn
```

