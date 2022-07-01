package main

import (
	"github.com/keptn/go-utils/pkg/sdk"
	"log"
)

const greetingsTriggeredEventType = "sh.keptn.event.greeting.triggered"
const serviceName = "greetings-service"

func main() {
	log.Fatal(sdk.NewKeptn(
		serviceName,
		sdk.WithTaskHandler(
			greetingsTriggeredEventType,
			NewGreetingsHandler()),
	).Start())
}
