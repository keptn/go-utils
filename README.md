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
  keptnutils "github.com/keptn/go-utils/pkg/utils"
}
```

Logging Example:
```go
keptnutils.Debug(keptncontext, message)
keptnutils.Info(keptncontext, message)
keptnutils.Error(keptncontext, message)
```

### CloudEvent Data
If you need to access data within CloudEvents:

```go
import {
	keptnevents "github.com/keptn/go-utils/pkg/events"
)
```

Example:

```go
func parseCloudEvent(event cloudevents.Event) (keptnevents.TestFinishedEventData, error) {
	eventData := &keptnevents.TestsFinishedEventData{}
	err := event.DataAs(eventData)
    
    return eventData, err
}
```

### Models
If you need to access Models for YAML files:

```go
import {
	keptnmodels "github.com/keptn/go-utils/pkg/models"
)
```
