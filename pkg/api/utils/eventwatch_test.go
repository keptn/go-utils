package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-openapi/strfmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type fakeEventGetter struct {
}

var t0 = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

var fakeEventsDB = map[string][]*models.KeptnContextExtendedCE{
	"ctx1": {
		{
			ID:             "ID1",
			Shkeptncontext: "ctx1",
			Time:           strfmt.DateTime(t0.Add(time.Second)),
		},
		{
			ID:             "ID2",
			Shkeptncontext: "ctx1",
			Time:           strfmt.DateTime(t0.Add(time.Second * 2)),
		},
		{
			ID:             "ID3",
			Shkeptncontext: "ctx1",
			Time:           strfmt.DateTime(t0.Add(time.Second * 3)),
		},
	},
	"ctx2": {
		{
			ID:             "ID1",
			Shkeptncontext: "ctx2",
			Time:           strfmt.DateTime(t0.Add(time.Second * 30)),
		},
		{
			ID:             "ID2",
			Shkeptncontext: "ctx2",
			Time:           strfmt.DateTime(t0.Add(time.Second * 31)),
		},
	},
}

func (eg fakeEventGetter) Get(filter *EventFilter) ([]*models.KeptnContextExtendedCE, error) {
	events := fakeEventsDB[filter.KeptnContext]
	fakeEventsDB = map[string][]*models.KeptnContextExtendedCE{}
	return events, nil
}

func withFakeEventGetter() EventWatcherOption {
	return func(ew *EventWatcher) {
		ew.eventGetter = fakeEventGetter{}
	}
}

type failingEventGetter struct {
}

func (eg failingEventGetter) Get(filter *EventFilter) ([]*models.KeptnContextExtendedCE, error) {
	return nil, errors.New("FAILED")
}

func withFilingEventGetter() EventWatcherOption {
	return func(ew *EventWatcher) {
		ew.eventGetter = failingEventGetter{}
	}
}

func TestEventWatcher(t *testing.T) {
	watcher := NewEventWatcher(
		withFakeEventGetter(),
		WithEventFilter(EventFilter{KeptnContext: "ctx1"}),
		WithCustomInterval(NewFakeSleeper()),
	)

	stream, _ := watcher.Watch(context.Background())
	events, ok := <-stream
	if !ok {
		t.Fatalf("unexpected closed channel")
	}
	assert.Equal(t, 3, len(events))
}

func TestEventWatcherCancel(t *testing.T) {
	watcher := NewEventWatcher(
		withFakeEventGetter(),
		WithEventFilter(EventFilter{KeptnContext: "ctx1"}),
		WithCustomInterval(NewFakeSleeper()),
	)

	stream, cancel := watcher.Watch(context.Background())
	cancel()

	_, ok := <-stream
	if ok {
		t.Fatalf("unexpected opened channel")
	}
}

func TestSortedGetter(t *testing.T) {

	firstTime := strfmt.DateTime(t0.Add(-time.Second * 2))
	secondTime := strfmt.DateTime(t0.Add(-time.Second))
	thirdTime := strfmt.DateTime(t0)

	events := []*models.KeptnContextExtendedCE{
		{Time: strfmt.DateTime(t0.Add(-time.Second))},
		{Time: strfmt.DateTime(t0)},
		{Time: strfmt.DateTime(t0.Add(-time.Second * 2))},
	}

	sortByTime(events)
	assert.Equal(t, events[0].Time, firstTime)
	assert.Equal(t, events[1].Time, secondTime)
	assert.Equal(t, events[2].Time, thirdTime)

	for _, e := range events {
		fmt.Println(e.Time)
	}
}
