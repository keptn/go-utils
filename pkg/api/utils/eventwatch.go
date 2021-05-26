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
	ticker          *time.Ticker
	timeout         <-chan time.Time
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
	defer func() {
		cancel()
		ew.ticker.Stop()
	}()

	for {
		// We need to query immediately because a time.Ticker cannot be configured
		// to emmit a tick event immediately
		ch <- ew.queryEvents(filter)
		select {
		// Query again once we receive a next tick
		case <-ew.ticker.C:
			continue
		// Close the channel and break out once we reach a timeout
		case <-ew.timeout:
			close(ch)
			return
		// Close the channel and break out once the user cancels via the context
		case <-ctx.Done():
			close(ch)
			return
		}
	}
}

func (ew *EventWatcher) queryEvents(filter EventFilter) []*models.KeptnContextExtendedCE {

	filter.FromTime = ew.nextCEFetchTime.Format("2006-01-02T15:04:05.000Z")
	events, err := ew.eventHandler.GetEvents(&filter)
	if err != nil {
		log.Printf("Unable to fetch events: %s", *err.Message)
	}
	SortByTime(events)
	if len(events) > 0 {
		if events[len(events)-1].Time.After(ew.nextCEFetchTime) {
			ew.nextCEFetchTime = events[len(events)-1].Time
		}
	}

	return events
}

// NewEventWatcher creates a new event watcher with the given options
func NewEventWatcher(eventHandler EventHandlerInterface, opts ...EventWatcherOption) *EventWatcher {
	e := &EventWatcher{
		nextCEFetchTime: time.Now().UTC(),
		eventHandler:    eventHandler,
		eventFilter:     EventFilter{},
		ticker:          time.NewTicker(10 * time.Second),
		timeout:         nil,
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

// WithInterval configures the EventWatcher to use a custom delay between each query
// You can use this to overwrite the default which is 10 * time.Second
func WithInterval(ticker *time.Ticker) EventWatcherOption {
	return func(ew *EventWatcher) {
		ew.ticker = ticker
	}
}

// WithTimeout configures the EventWatcher to use a custom timeout specifying
// after which duration the watcher shall stop
func WithTimeout(duration time.Duration) EventWatcherOption {
	return func(ew *EventWatcher) {
		ew.timeout = time.After(duration)
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
		return events[i].Time.Before(events[j].Time)
	})
}
