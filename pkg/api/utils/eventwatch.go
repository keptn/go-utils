package api

import (
	"errors"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"sync"
	"time"
)

// EventWatchInterface is the interface for watching events
type EventWatchInterface interface {
	// Stop stops the event watch
	Stop()
	// Events starts the internal fetching of events and returns a channel to
	// returns the channel to get the events
	Events() <-chan models.KeptnContextExtendedCE
}

// EventWatcher is the implementation of the EventWatchInterface
type EventWatcher struct {
	sync.Mutex
	events      chan models.KeptnContextExtendedCE
	nextCETime  time.Time
	stopped     bool
	eventGetter EventGetter
	eventFilter EventFilter
	sleeper     Sleeper
}

// Stop stops the event watcher
func (ew *EventWatcher) Stop() {
	ew.Lock()
	defer ew.Unlock()
	if !ew.stopped {
		ew.stopped = true
	}
}

// Events starts the query loop and returns the channel to get the events
func (ew *EventWatcher) Events() <-chan models.KeptnContextExtendedCE {
	ew.Lock()
	defer ew.Unlock()
	if !ew.stopped {
		go ew.queryPeriodically(ew.eventFilter)
	}
	return ew.events
}

// isStopped checks whether the event watcher is currently running
func (ew *EventWatcher) isStopped() bool {
	ew.Lock()
	defer ew.Unlock()
	return ew.stopped
}

// queryPeriodically quiries the database for events based on the given filter
// if events are found, they get sent down the channel to the client
func (ew *EventWatcher) queryPeriodically(filter EventFilter) {

	defer close(ew.events)
	defer ew.Stop()

	for {

		filter.FromTime = ew.nextCETime.Add(-(time.Second * 10)).Format("2006-01-02T15:04:05.000Z")
		if ew.isStopped() {
			return
		}
		events, err := ew.eventGetter.Get(&filter)

		//TODO: implement retry
		if err != nil {
			fmt.Errorf("Unable to fetch events")
			return
		}

		for _, e := range events {
			ew.events <- *e
		}
		ew.sleeper.Sleep()
		ew.nextCETime = time.Now().UTC()
	}
}

// NewEventWatcher creates a new event watcher with the given options
func NewEventWatcher(opts ...EventWatcherOption) *EventWatcher {
	e := &EventWatcher{
		Mutex:       sync.Mutex{},
		events:      make(chan models.KeptnContextExtendedCE),
		nextCETime:  time.Now(),
		stopped:     false,
		eventGetter: NewDefaultEventGetter("localhost/api"),
		eventFilter: EventFilter{},
		sleeper:     NewConfigurableSleeper(10 * time.Second),
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

// EventWatcherOptions can be used to provide functionality to configure the EventWatcher
type EventWatcherOption func(*EventWatcher)

// WithEventFilter configures the EventWatcher to use a filter
func WithEventFilter(filter EventFilter) EventWatcherOption {
	return func(ew *EventWatcher) {
		ew.eventFilter = filter
	}
}

// WithStartTime configures the EventWatcher to use a specific start timestamp for the first query
func WithStartTime(startTime time.Time) EventWatcherOption {
	return func(ew *EventWatcher) {
		ew.nextCETime = startTime
	}
}

// WithAuthenticatedEventGetter configures the EventWatcher to use a authenticated event handler
// for querying the events against the event database
func WithAuthenticatedEventGetter(baseUrl, token string) EventWatcherOption {
	return func(ew *EventWatcher) {
		ew.eventGetter = NewAuthenticatedEventGetter(baseUrl, token)
	}
}

// WithCustomInterval configures the EventWatcher to use a cusstom delay between each query
// You can use this to overwrite the default  which is 10 * time.Second

func WithCustomInterval(sleeper Sleeper) EventWatcherOption {
	return func(ew *EventWatcher) {
		ew.sleeper = sleeper
	}
}

// EventGetter defines the interface for getting the events from the event database
type EventGetter interface {
	// Get queries the event database for events matching the given filter
	Get(filter *EventFilter) ([]*models.KeptnContextExtendedCE, error)
}

// NewAuthenticatedEventGetter creates a new instance of an EventGetter which authenticates itself
// with a given token
func NewAuthenticatedEventGetter(baseUrl, token string) CloudEventGetter {
	return CloudEventGetter{
		handler: NewAuthenticatedEventHandler(
			baseUrl,
			token,
			"x-token",
			nil,
			"http"),
	}
}

// NewDefaultEventGetter creates a new instance of an EventGetter
func NewDefaultEventGetter(baseUrl string) CloudEventGetter {
	return CloudEventGetter{
		handler: NewEventHandler(baseUrl),
	}
}

// CloudEventGetter is an implementation of an EventGetter
type CloudEventGetter struct {
	handler *EventHandler
}

// Get queries the event database for events matching the given filter
func (eg CloudEventGetter) Get(filter *EventFilter) ([]*models.KeptnContextExtendedCE, error) {

	events, err := eg.handler.GetEvents(filter)
	if err != nil {
		return nil, errors.New(*err.Message)
	}
	return events, nil
}
