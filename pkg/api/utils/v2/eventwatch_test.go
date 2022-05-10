package v2

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/timeutils"
	"github.com/stretchr/testify/assert"
)

type fakeEventHandler struct {
	data map[string][]*models.KeptnContextExtendedCE
}

func (fh *fakeEventHandler) GetEvents(filter *EventFilter) ([]*models.KeptnContextExtendedCE, *models.Error) {
	events := fh.data[filter.KeptnContext]
	fh.data = map[string][]*models.KeptnContextExtendedCE{}
	return events, nil
}

func (fh *fakeEventHandler) GetEventsWithRetry(filter *EventFilter, maxRetries int, retrySleepTime time.Duration) ([]*models.KeptnContextExtendedCE, error) {
	panic("not implemented")
}

func newFakeEventHandler() *fakeEventHandler {
	return &fakeEventHandler{
		data: map[string][]*models.KeptnContextExtendedCE{
			"ctx1": {
				{
					ID:             "ID1",
					Shkeptncontext: "ctx1",
					Time:           t0.Add(time.Second),
				},
				{
					ID:             "ID2",
					Shkeptncontext: "ctx1",
					Time:           t0.Add(time.Second * 2),
				},
				{
					ID:             "ID3",
					Shkeptncontext: "ctx1",
					Time:           t0.Add(time.Second * 3),
				},
			},
			"ctx2": {
				{
					ID:             "ID1",
					Shkeptncontext: "ctx2",
					Time:           t0.Add(time.Second * 30),
				},
				{
					ID:             "ID2",
					Shkeptncontext: "ctx2",
					Time:           t0.Add(time.Second * 31),
				},
			},
		},
	}
}

var t0 = time.Date(2021, 1, 1, 0, 0, 0, 0, time.UTC)

func TestEventWatcher(t *testing.T) {
	watcher := NewEventWatcher(newFakeEventHandler(),
		WithEventFilter(EventFilter{KeptnContext: "ctx1"}),
		WithInterval(time.NewTicker(1)),
	)

	stream, _ := watcher.Watch(context.Background())
	events, ok := <-stream
	if !ok {
		t.Fatalf("unexpected closed channel")
	}
	assert.Equal(t, 3, len(events))
}

func TestEventWatcherCancel(t *testing.T) {
	watcher := NewEventWatcher(newFakeEventHandler(),
		WithEventFilter(EventFilter{KeptnContext: "ctx1"}),
		WithInterval(time.NewTicker(1)),
	)

	stream, cancel := watcher.Watch(context.Background())
	cancel()
	for ev := range stream {
		fmt.Println(ev)
	}

	_, ok := <-stream
	if ok {
		t.Fatalf("unexpected opened channel")
	}
}

func TestEventWatcherTimeout(t *testing.T) {
	watcher := NewEventWatcher(newFakeEventHandler(),
		WithEventFilter(EventFilter{KeptnContext: "ctx1"}),
		WithTimeout(10*time.Millisecond),
	)

	stream, _ := watcher.Watch(context.Background())
	for ev := range stream {
		fmt.Println(ev)
	}

	_, ok := <-stream
	if ok {
		t.Fatalf("unexpected opened channel")
	}

}

func TestSortedGetter(t *testing.T) {

	firstTime := timeutils.GetKeptnTimeStamp(t0.Add(-time.Second * 2))
	secondTime := timeutils.GetKeptnTimeStamp(t0.Add(-time.Second))
	thirdTime := timeutils.GetKeptnTimeStamp(t0)

	events := []*models.KeptnContextExtendedCE{
		{Time: t0.Add(-time.Second)},
		{Time: t0},
		{Time: t0.Add(-time.Second * 2)},
	}

	SortByTime(events)
	assert.Equal(t, timeutils.GetKeptnTimeStamp(events[0].Time), firstTime)
	assert.Equal(t, timeutils.GetKeptnTimeStamp(events[1].Time), secondTime)
	assert.Equal(t, timeutils.GetKeptnTimeStamp(events[2].Time), thirdTime)

	for _, e := range events {
		fmt.Println(e.Time)
	}
}
