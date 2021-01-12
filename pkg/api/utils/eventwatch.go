package api

import (
	"context"
	"errors"
	"github.com/keptn/go-utils/pkg/api/models"
	"log"
	"sort"
	"time"
)

// EventWatcher implements the logic to query for events and provide them to the client
type EventWatcher struct {
	nextCETime  time.Time
	eventGetter EventGetter
	eventFilter EventFilter
	sleeper     Sleeper
}

// Watch starts the watch loop.
// It returns a channel to get the actual events as well as a context.CancelFunc in order
// to stop the watch routine
func (ew *EventWatcher) Watch(ctx context.Context) (<-chan []*models.KeptnContextExtendedCE, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	ch := make(chan []*models.KeptnContextExtendedCE)
	go ew.fetch(ctx, cancel, ch, ew.eventFilter)
	return ch, cancel
}

func (ew *EventWatcher) fetch(ctx context.Context, cancel context.CancelFunc, ch chan<- []*models.KeptnContextExtendedCE, filter EventFilter) {
	defer cancel()
	for {
		select {
		case <-ctx.Done():
			close(ch)
			return
		case ch <- ew.queryEvents(filter):
			ew.sleeper.Sleep()
		}
	}
}

func (ew *EventWatcher) queryEvents(filter EventFilter) []*models.KeptnContextExtendedCE {
	filter.FromTime = ew.nextCETime.Format("2006-01-02T15:04:05.000Z")
	ew.nextCETime = time.Now().UTC()
	events, err := ew.eventGetter.Get(&filter)
	if err != nil {
		log.Fatal("Unable to fetch events")
	}
	return events
}

// NewEventWatcher creates a new event watcher with the given options
func NewEventWatcher(opts ...EventWatcherOption) *EventWatcher {
	e := &EventWatcher{
		nextCETime:  time.Now(),
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
		ew.eventGetter = newAuthenticatedEventGetter(baseUrl, token)
	}
}

func WithAuthenticatedSortingEventGetter(baseUrl, token string) EventWatcherOption {
	return func(ew *EventWatcher) {
		ew.eventGetter = newAuthenticatedSortingEventGetter(baseUrl, token)
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
func newAuthenticatedEventGetter(baseUrl, token string) *DefaultEventGetter {
	return &DefaultEventGetter{
		handler: NewAuthenticatedEventHandler(
			baseUrl,
			token,
			"x-token",
			nil,
			"http"),
	}
}

func newAuthenticatedSortingEventGetter(baseUrl, token string) *SortingEventGetter {
	return &SortingEventGetter{
		handler: NewAuthenticatedEventHandler(
			baseUrl,
			token,
			"x-token",
			nil,
			"http"),
	}
}

// NewDefaultEventGetter creates a new instance of an EventGetter
func NewDefaultEventGetter(baseUrl string) *DefaultEventGetter {
	return &DefaultEventGetter{
		handler: NewEventHandler(baseUrl),
	}
}

// DefaultEventGetter is an implementation of an EventGetter
type DefaultEventGetter struct {
	handler *EventHandler
}

// Get queries the event database for events matching the given filter
func (eg DefaultEventGetter) Get(filter *EventFilter) ([]*models.KeptnContextExtendedCE, error) {

	events, err := eg.handler.GetEvents(filter)
	if err != nil {
		return nil, errors.New(*err.Message)
	}
	return events, nil
}

// SortingEventGetter is an implementation of an EventGetter which returns the events in the chronological
// time order (oldest to newest)
type SortingEventGetter struct {
	handler *EventHandler
}

// Get queries the event database for events matching the given filter
// Moreover, it sorts the fetched events by time in increasing order (oldest to newest)
func (eg SortingEventGetter) Get(filter *EventFilter) ([]*models.KeptnContextExtendedCE, error) {
	events, err := eg.handler.GetEvents(filter)
	if err != nil {
		return nil, errors.New(*err.Message)
	}
	sortByTime(events)
	return events, nil
}

func sortByTime(events []*models.KeptnContextExtendedCE) {
	sort.Slice(events, func(i, j int) bool {
		return time.Time(events[i].Time).Before(time.Time(events[j].Time))
	})
}
