package api

import (
	"errors"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

type fakeEventGetter struct {
}

func (eg fakeEventGetter) Get(filter *EventFilter) ([]*models.KeptnContextExtendedCE, error) {

	ceCtx1_1 := &models.KeptnContextExtendedCE{
		ID:             "ID1",
		Shkeptncontext: "ctx1",
	}

	ceCtx1_2 := &models.KeptnContextExtendedCE{
		ID:             "ID2",
		Shkeptncontext: "ctx1",
	}

	ceCtx2_1 := &models.KeptnContextExtendedCE{
		ID:             "ID3",
		Shkeptncontext: "ctx2",
	}

	ceCtx2_2 := &models.KeptnContextExtendedCE{
		ID:             "ID4",
		Shkeptncontext: "ctx2",
	}

	ceCtx2_3 := &models.KeptnContextExtendedCE{
		ID:             "ID5",
		Shkeptncontext: "ctx2",
	}

	if filter.KeptnContext == "ctx1" {

		return []*models.KeptnContextExtendedCE{ceCtx1_1, ceCtx1_2}, nil
	}

	if filter.KeptnContext == "ctx2" {

		return []*models.KeptnContextExtendedCE{ceCtx2_1, ceCtx2_2, ceCtx2_3}, nil
	}
	return []*models.KeptnContextExtendedCE{ceCtx1_1, ceCtx1_2, ceCtx2_1, ceCtx2_2, ceCtx2_3}, nil
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
	var event models.KeptnContextExtendedCE
	watcher := NewEventWatcher(
		withFakeEventGetter(),
		WithEventFilter(EventFilter{KeptnContext: "ctx1"}),
		WithCustomInterval(NewFakeSleeper()),
	)

	eventStream := watcher.Events()
	event, ok := <-eventStream
	if !ok {
		t.Fatalf("unexpected closed channel")
	}
	assert.Equal(t, "ID1", event.ID)

	event = <-eventStream
	if !ok {
		t.Fatalf("unexpected closed channel")
	}
	assert.Equal(t, "ID2", event.ID)

	watcher.Stop()
	_, ok = <-eventStream
	if ok {
		t.Fatalf("unexpected open channel")
	}
}

func TestEventWatcherWithError(t *testing.T) {
	watcher := NewEventWatcher(
		withFilingEventGetter())

	eventStream := watcher.Events()

	_, ok := <-eventStream
	if ok {
		t.Fatalf("unexpected open channel")
	}

}
