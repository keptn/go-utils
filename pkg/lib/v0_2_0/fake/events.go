package fake

import (
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type FakeEventSender struct {
	SentEvents []cloudevents.Event
	Reactors   map[string]func(event cloudevents.Event) error
}

func (es *FakeEventSender) SendEvent(event cloudevents.Event) error {
	if es.Reactors != nil {
		for eventTypeSelector, reactor := range es.Reactors {
			if eventTypeSelector == "*" || eventTypeSelector == event.Type() {
				if err := reactor(event); err != nil {
					return err
				}
			}
		}
	}
	es.SentEvents = append(es.SentEvents, event)
	return nil
}

func (es *FakeEventSender) AssertSentEventTypes(eventTypes []string) error {
	if len(es.SentEvents) != len(eventTypes) {
		return fmt.Errorf("expected %d event, got %d", len(es.SentEvents), len(eventTypes))
	}
	for index, event := range es.SentEvents {
		if event.Type() != eventTypes[index] {
			return fmt.Errorf("received event type '%s' != %s", event.Type(), eventTypes[index])
		}
	}
	return nil
}

func (es *FakeEventSender) AddReactor(eventTypeSelector string, reactor func(event cloudevents.Event) error) {
	if es.Reactors == nil {
		es.Reactors = map[string]func(event cloudevents.Event) error{}
	}
	es.Reactors[eventTypeSelector] = reactor
}
