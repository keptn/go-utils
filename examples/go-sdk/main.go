package main

import (
	"github.com/keptn/go-utils/pkg/sdk"
	"github.com/sirupsen/logrus"
	"log"
)

const greetingsTriggeredEventType = "sh.keptn.event.greeting.triggered"
const serviceName = "greetings-service"

func main() {
	logrus.SetLevel(logrus.DebugLevel)
	log.Fatal(sdk.NewKeptn(
		// the name of your keptn service
		serviceName,
		// the task handler containing logic to handle the
		// "sh.keptn.event.greeting.triggered" event
		sdk.WithTaskHandler(
			greetingsTriggeredEventType,
			NewGreetingsHandler()),
		// using logrus library as a logger
		sdk.WithLogger(logrus.StandardLogger()),
	).Start())
}
