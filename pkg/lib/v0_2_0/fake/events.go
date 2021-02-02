package fake

import (
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
)

// EventSender fakes the sending of CloudEvents
type EventSender struct {
	SentEvents []cloudevents.Event
	Reactors   map[string]func(event cloudevents.Event) error
}

// SendEvent fakes the sending of CloudEvents
func (es *EventSender) SendEvent(event cloudevents.Event) error {
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

// AssertSentEventTypes checks if the given event types have been passed to the SendEvent function
func (es *EventSender) AssertSentEventTypes(eventTypes []string) error {
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

// AddReactor adds custom logic that should be applied when SendEvent is called for the given event type
func (es *EventSender) AddReactor(eventTypeSelector string, reactor func(event cloudevents.Event) error) {
	if es.Reactors == nil {
		es.Reactors = map[string]func(event cloudevents.Event) error{}
	}
	es.Reactors[eventTypeSelector] = reactor
}
