package api

import (
	"context"
	"github.com/keptn/go-utils/pkg/api/models"
	"log"
	"sort"
	"time"
)

// EventWatcher implements the logic to query for events and provide them to the client
type EventWatcher struct {
	nextCEFetchTime time.Time
	eventHandler    EventHandlerInterface
	eventFilter     EventFilter
	sleeper         Sleeper
	manipulator     EventManipulatorFunc
}

// Watch starts the watch loop and returns a channel to get the actual events as well as a context.CancelFunc in order
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
	filter.FromTime = ew.nextCEFetchTime.Format("2006-01-02T15:04:05.000Z")
	ew.nextCEFetchTime = time.Now().UTC()
	events, err := ew.eventHandler.GetEvents(&filter)
	if err != nil {
		log.Fatal("Unable to fetch events")
	}
	ew.manipulator(events)
	return events
}

// NewEventWatcher creates a new event watcher with the given options
func NewEventWatcher(eventHandler EventHandlerInterface, opts ...EventWatcherOption) *EventWatcher {
	e := &EventWatcher{
		nextCEFetchTime: time.Now(),
		eventHandler:    eventHandler,
		eventFilter:     EventFilter{},
		sleeper:         NewConfigurableSleeper(10 * time.Second),
		manipulator:     func(ces []*models.KeptnContextExtendedCE) {},
	}

	for _, opt := range opts {
		opt(e)
	}

	return e
}

// EventWatcherOption can be used to configure the EventWatcher
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
		ew.nextCEFetchTime = startTime
	}
}

// WithEventManipulator configures the EventWatcher to manipulate the fetched events using the given EventSorterFunc
func WithEventManipulator(sorter EventManipulatorFunc) EventWatcherOption {
	return func(ew *EventWatcher) {
		ew.manipulator = sorter
	}
}

// WithCustomInterval configures the EventWatcher to use a custom delay between each query
// You can use this to overwrite the default  which is 10 * time.Second
func WithCustomInterval(sleeper Sleeper) EventWatcherOption {
	return func(ew *EventWatcher) {
		ew.sleeper = sleeper
	}
}

// EventHandlerInterface is the api to fetch events from the event store
type EventHandlerInterface interface {
	GetEventsWithRetry(filter *EventFilter, maxRetries int, retrySleepTime time.Duration) ([]*models.KeptnContextExtendedCE, error)
	GetEvents(filter *EventFilter) ([]*models.KeptnContextExtendedCE, *models.Error)
}

// EventManipulatorFunc can be used to manipulate a slice of events
type EventManipulatorFunc func([]*models.KeptnContextExtendedCE)

// SortByTime sorts the event slice by time (oldest to newest)
func SortByTime(events []*models.KeptnContextExtendedCE) {
	sort.Slice(events, func(i, j int) bool {
		return time.Time(events[i].Time).Before(time.Time(events[j].Time))
	})
}
