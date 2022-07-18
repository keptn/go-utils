package nats

import (
	"context"
	"errors"
	"fmt"
	"github.com/keptn/go-utils/pkg/sdk/connector/types"
	"sync"
	"testing"
	"time"

	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/strutils"
	nats2 "github.com/keptn/go-utils/pkg/sdk/connector/nats"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/require"
)

type NATSConnectorMock struct {
	SubscribeFn                 func(string, nats2.ProcessEventFn) error
	QueueSubscribeFn            func(string, string, nats2.ProcessEventFn) error
	SubscribeMultipleFn         func([]string, nats2.ProcessEventFn) error
	QueueSubscribeMultipleFn    func([]string, string, nats2.ProcessEventFn) error
	queueSubscribeMultipleCalls int
	PublishFn                   func(ce models.KeptnContextExtendedCE) error
	publishCalls                int
	DisconnectFn                func() error
	disconnectCalls             int
	UnsubscribeAllFn            func() error
	unsubscribeAllCalls         int
	QueueGroup                  string
	ProcessEventFn              nats2.ProcessEventFn
	mtx                         sync.RWMutex
}

func (ncm *NATSConnectorMock) QueueSubscribeMultipleCalls() int {
	ncm.mtx.RLock()
	defer ncm.mtx.RUnlock()
	return ncm.queueSubscribeMultipleCalls
}

func (ncm *NATSConnectorMock) PublishCalls() int {
	ncm.mtx.RLock()
	defer ncm.mtx.RUnlock()
	return ncm.publishCalls
}

func (ncm *NATSConnectorMock) DisconnectCalls() int {
	ncm.mtx.RLock()
	defer ncm.mtx.RUnlock()
	return ncm.disconnectCalls
}

func (ncm *NATSConnectorMock) UnsubscribeAllCalls() int {
	ncm.mtx.RLock()
	defer ncm.mtx.RUnlock()
	return ncm.unsubscribeAllCalls
}

func (ncm *NATSConnectorMock) Subscribe(subject string, fn nats2.ProcessEventFn) error {
	ncm.mtx.Lock()
	defer ncm.mtx.Unlock()
	if ncm.SubscribeFn != nil {
		return ncm.SubscribeFn(subject, fn)
	}
	panic("implement me")
}

func (ncm *NATSConnectorMock) QueueSubscribe(subject string, queueGroup string, fn nats2.ProcessEventFn) error {
	ncm.mtx.Lock()
	defer ncm.mtx.Unlock()
	if ncm.QueueSubscribeFn != nil {
		return ncm.QueueSubscribeFn(queueGroup, subject, fn)
	}
	panic("implement me")
}

func (ncm *NATSConnectorMock) SubscribeMultiple(subjects []string, fn nats2.ProcessEventFn) error {
	ncm.mtx.Lock()
	defer ncm.mtx.Unlock()
	if ncm.SubscribeMultipleFn != nil {
		return ncm.SubscribeMultipleFn(subjects, fn)
	}
	panic("implement me")
}

func (ncm *NATSConnectorMock) QueueSubscribeMultiple(subjects []string, queueGroup string, fn nats2.ProcessEventFn) error {
	ncm.mtx.Lock()
	defer ncm.mtx.Unlock()
	ncm.ProcessEventFn = fn
	ncm.queueSubscribeMultipleCalls++
	if ncm.QueueSubscribeMultipleFn != nil {
		return ncm.QueueSubscribeMultipleFn(subjects, queueGroup, fn)
	}
	panic("implement me")
}

func (ncm *NATSConnectorMock) Publish(event models.KeptnContextExtendedCE) error {
	ncm.mtx.Lock()
	defer ncm.mtx.Unlock()
	ncm.publishCalls++
	if ncm.PublishFn != nil {
		return ncm.PublishFn(event)
	}
	panic("implement me")
}

func (ncm *NATSConnectorMock) Disconnect() error {
	ncm.mtx.Lock()
	defer ncm.mtx.Unlock()
	ncm.disconnectCalls++
	if ncm.DisconnectFn != nil {
		return ncm.DisconnectFn()
	}
	panic("implement me")
}

func (ncm *NATSConnectorMock) UnsubscribeAll() error {
	ncm.mtx.Lock()
	defer ncm.mtx.Unlock()
	ncm.unsubscribeAllCalls++
	if ncm.UnsubscribeAllFn != nil {
		return ncm.UnsubscribeAllFn()
	}
	panic("implement me")
}

func TestEventSourceForwardsEventToChannel(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{
		QueueSubscribeMultipleFn: func(subjects []string, queueGroup string, fn nats2.ProcessEventFn) error { return nil },
		UnsubscribeAllFn:         func() error { return nil },
	}
	eventChannel := make(chan types.EventUpdate)
	eventSource := New(natsConnectorMock)
	wg := &sync.WaitGroup{}
	wg.Add(1)
	eventSource.Start(context.TODO(), types.RegistrationData{}, eventChannel, make(chan error), wg)
	eventSource.OnSubscriptionUpdate([]models.EventSubscription{{Event: "a"}})
	event := models.KeptnContextExtendedCE{ID: "id"}
	jsonEvent, _ := event.ToJSON()
	e := &nats.Msg{Data: jsonEvent, Sub: &nats.Subscription{Subject: "subscription"}} //models.KeptnContextExtendedCE{ID: "id"}
	go natsConnectorMock.ProcessEventFn(e)
	eventFromChan := <-eventChannel
	require.Equal(t, eventFromChan.KeptnEvent, event)
}

func TestEventSourceCancelDisconnectsFromBroker(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{
		QueueSubscribeMultipleFn: func(subjects []string, queueGroup string, fn nats2.ProcessEventFn) error { return nil },
		UnsubscribeAllFn:         func() error { return nil },
	}
	ctx, cancel := context.WithCancel(context.TODO())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	New(natsConnectorMock).Start(ctx, types.RegistrationData{}, make(chan types.EventUpdate), make(chan error), wg)
	cancel()
	require.Eventually(t, func() bool { return natsConnectorMock.UnsubscribeAllCalls() == 1 }, 2*time.Second, 100*time.Millisecond)
}

func TestEventSourceCallsWaitGroupDuringCancellation(t *testing.T) {
	t.Run("WaitGroup called", func(t *testing.T) {
		natsConnectorMock := &NATSConnectorMock{
			QueueSubscribeMultipleFn: func(subjects []string, queueGroup string, fn nats2.ProcessEventFn) error { return nil },
			UnsubscribeAllFn:         func() error { return nil },
		}
		ctx, cancel := context.WithCancel(context.TODO())
		wg := &sync.WaitGroup{}
		wg.Add(1)
		New(natsConnectorMock).Start(ctx, types.RegistrationData{}, make(chan types.EventUpdate), make(chan error), wg)
		cancel()
		wg.Wait()
	})
	t.Run("WaitGroup called - error in shutdown logic", func(t *testing.T) {
		natsConnectorMock := &NATSConnectorMock{
			QueueSubscribeMultipleFn: func(subjects []string, queueGroup string, fn nats2.ProcessEventFn) error { return nil },
			UnsubscribeAllFn:         func() error { return fmt.Errorf("ohoh") },
		}
		ctx, cancel := context.WithCancel(context.TODO())
		wg := &sync.WaitGroup{}
		wg.Add(1)
		New(natsConnectorMock).Start(ctx, types.RegistrationData{}, make(chan types.EventUpdate), make(chan error), wg)
		cancel()
		wg.Wait()
	})
}

func TestEventSourceCancelDisconnectFromBrokerFails(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{
		QueueSubscribeMultipleFn: func(subjects []string, queueGroup string, fn nats2.ProcessEventFn) error { return nil },
		UnsubscribeAllFn:         func() error { return fmt.Errorf("error occured") },
	}
	ctx, cancel := context.WithCancel(context.TODO())
	wg := &sync.WaitGroup{}
	wg.Add(1)
	New(natsConnectorMock).Start(ctx, types.RegistrationData{}, make(chan types.EventUpdate), make(chan error), wg)
	cancel()
	require.Eventually(t, func() bool { return natsConnectorMock.UnsubscribeAllCalls() == 1 }, 2*time.Second, 100*time.Millisecond)
}

func TestEventSourceQueueSubscribeFails(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{QueueSubscribeMultipleFn: func(strings []string, s string, fn nats2.ProcessEventFn) error { return fmt.Errorf("error occured") }}
	eventSource := New(natsConnectorMock)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	err := eventSource.Start(context.TODO(), types.RegistrationData{}, make(chan types.EventUpdate), make(chan error), wg)
	require.Error(t, err)
}

func TestEventSourceOnSubscriptionUpdate(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{
		QueueSubscribeMultipleFn: func(subjects []string, queueGroup string, fn nats2.ProcessEventFn) error { return nil },
		UnsubscribeAllFn:         func() error { return nil },
	}
	eventSource := New(natsConnectorMock)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	err := eventSource.Start(context.TODO(), types.RegistrationData{}, make(chan types.EventUpdate), make(chan error), wg)
	require.NoError(t, err)
	require.Equal(t, 1, natsConnectorMock.queueSubscribeMultipleCalls)
	eventSource.OnSubscriptionUpdate([]models.EventSubscription{{Event: "a"}})
	require.Equal(t, 1, natsConnectorMock.unsubscribeAllCalls)
	require.Equal(t, 2, natsConnectorMock.queueSubscribeMultipleCalls)
}

func TestEventSourceOnSubscriptionupdateWithDuplicatedSubjects(t *testing.T) {
	var receivedSubjects []string
	natsConnectorMock := &NATSConnectorMock{
		QueueSubscribeMultipleFn: func(subjects []string, queueGroup string, fn nats2.ProcessEventFn) error {
			receivedSubjects = subjects
			return nil
		},
		UnsubscribeAllFn: func() error { return nil },
	}
	eventSource := New(natsConnectorMock)
	err := eventSource.Start(context.TODO(), types.RegistrationData{}, make(chan types.EventUpdate), make(chan error), &sync.WaitGroup{})
	require.NoError(t, err)
	require.Equal(t, 1, natsConnectorMock.queueSubscribeMultipleCalls)
	eventSource.OnSubscriptionUpdate([]models.EventSubscription{{Event: "a"}, {Event: "a"}})
	require.Equal(t, 1, natsConnectorMock.unsubscribeAllCalls)
	require.Equal(t, 2, natsConnectorMock.queueSubscribeMultipleCalls)
	require.Equal(t, 1, len(receivedSubjects))
}

func TestEventSourceOnSubscriptiOnUpdateUnsubscribeAllFails(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{
		QueueSubscribeMultipleFn: func(subjects []string, queueGroup string, fn nats2.ProcessEventFn) error { return nil },
		UnsubscribeAllFn:         func() error { return fmt.Errorf("error occured") },
	}
	eventSource := New(natsConnectorMock)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	err := eventSource.Start(context.TODO(), types.RegistrationData{}, make(chan types.EventUpdate), make(chan error), wg)
	require.NoError(t, err)
	require.Equal(t, 1, natsConnectorMock.queueSubscribeMultipleCalls)
	eventSource.OnSubscriptionUpdate([]models.EventSubscription{{Event: "a"}})
	require.Equal(t, 1, natsConnectorMock.unsubscribeAllCalls)
	require.Equal(t, 1, natsConnectorMock.queueSubscribeMultipleCalls)
}

func TestEventSourceOnSubscriptionUpdateQueueSubscribeMultipleFails(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{
		QueueSubscribeMultipleFn: func(subjects []string, queueGroup string, fn nats2.ProcessEventFn) error { return nil },
		UnsubscribeAllFn:         func() error { return nil },
	}
	eventSource := New(natsConnectorMock)
	wg := &sync.WaitGroup{}
	wg.Add(1)

	err := eventSource.Start(context.TODO(), types.RegistrationData{}, make(chan types.EventUpdate), make(chan error), wg)
	require.NoError(t, err)
	require.Equal(t, 1, natsConnectorMock.queueSubscribeMultipleCalls)
	natsConnectorMock.QueueSubscribeMultipleFn = func(subjects []string, queueGroup string, fn nats2.ProcessEventFn) error {
		return fmt.Errorf("error occured")
	}
	eventSource.OnSubscriptionUpdate([]models.EventSubscription{{Event: "a"}})
	require.Equal(t, 1, natsConnectorMock.unsubscribeAllCalls)
	require.Equal(t, 2, natsConnectorMock.queueSubscribeMultipleCalls)
}

func TestEventSourceGetSender(t *testing.T) {
	event := models.KeptnContextExtendedCE{ID: "id", Type: strutils.Stringp("something")}
	natsConnectorMock := &NATSConnectorMock{
		PublishFn: func(ce models.KeptnContextExtendedCE) error {
			require.Equal(t, event, ce)
			return nil
		},
	}
	sendFn := New(natsConnectorMock).Sender()
	require.NotNil(t, sendFn)
	err := sendFn(event)
	require.NoError(t, err)
	require.Equal(t, 1, natsConnectorMock.publishCalls)
}

func TestEventSourceSenderFails(t *testing.T) {
	event := models.KeptnContextExtendedCE{ID: "id", Type: strutils.Stringp("something")}
	natsConnectorMock := &NATSConnectorMock{
		PublishFn: func(ce models.KeptnContextExtendedCE) error {
			require.Equal(t, event, ce)
			return fmt.Errorf("error occured")
		},
	}
	sendFn := New(natsConnectorMock).Sender()
	require.NotNil(t, sendFn)
	err := sendFn(event)
	require.Error(t, err)
	require.Equal(t, 1, natsConnectorMock.publishCalls)
}

func TestEventSourceStopUnsubscribesFromEventBroker(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{
		QueueSubscribeMultipleFn: func(strings []string, s string, fn nats2.ProcessEventFn) error {
			return nil
		},
		UnsubscribeAllFn: func() error {
			return nil
		},
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)

	es := New(natsConnectorMock)
	err := es.Start(context.TODO(), types.RegistrationData{}, make(chan types.EventUpdate), make(chan error), wg)
	require.NoError(t, err)
	err = es.Stop()
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		return natsConnectorMock.unsubscribeAllCalls == 1
	}, 100*time.Millisecond, 10*time.Millisecond)
}

func TestEventSourceCleanup(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{
		DisconnectFn: func() error {
			return nil
		},
	}

	es := New(natsConnectorMock)
	err := es.Cleanup()
	require.NoError(t, err)
	require.Eventually(t, func() bool {
		return natsConnectorMock.DisconnectCalls() == 1
	}, 100*time.Millisecond, 10*time.Millisecond)
}

func TestEventSourceCleanupFails(t *testing.T) {
	natsConnectorMock := &NATSConnectorMock{
		DisconnectFn: func() error {
			return errors.New("oops")
		},
	}

	es := New(natsConnectorMock)
	err := es.Cleanup()
	require.Error(t, err)
	require.Eventually(t, func() bool {
		return natsConnectorMock.DisconnectCalls() == 1
	}, 100*time.Millisecond, 10*time.Millisecond)
}
