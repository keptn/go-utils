package nats

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/keptn/go-utils/pkg/sdk/connector/types"
	"reflect"
	"sort"
	"sync"

	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/sdk/connector/logger"
	natseventsource "github.com/keptn/go-utils/pkg/sdk/connector/nats"
	"github.com/nats-io/nats.go"
)

// NATSEventSource is an implementation of EventSource
// that is using the NATS event broker internally
type NATSEventSource struct {
	currentSubjects []string
	connector       natseventsource.NATS
	eventProcessFn  natseventsource.ProcessEventFn
	queueGroup      string
	logger          logger.Logger
	quitC           chan struct{}
}

// New creates a new NATSEventSource
func New(natsConnector natseventsource.NATS, opts ...func(source *NATSEventSource)) *NATSEventSource {
	e := &NATSEventSource{
		currentSubjects: []string{},
		connector:       natsConnector,
		eventProcessFn:  func(event *nats.Msg) error { return nil },
		quitC:           make(chan struct{}, 1),
		logger:          logger.NewDefaultLogger(),
	}
	for _, o := range opts {
		o(e)
	}
	return e
}

// WithLogger sets the logger to use
func WithLogger(logger logger.Logger) func(*NATSEventSource) {
	return func(ns *NATSEventSource) {
		ns.logger = logger
	}
}

func (n *NATSEventSource) Start(ctx context.Context, registrationData types.RegistrationData, eventChannel chan types.EventUpdate, errChan chan error, wg *sync.WaitGroup) error {
	n.queueGroup = registrationData.Name
	n.eventProcessFn = func(event *nats.Msg) error {
		keptnEvent := models.KeptnContextExtendedCE{}
		if err := json.Unmarshal(event.Data, &keptnEvent); err != nil {
			return fmt.Errorf("could not unmarshal message: %w", err)
		}
		eventChannel <- types.EventUpdate{
			KeptnEvent: keptnEvent,
			MetaData:   types.EventUpdateMetaData{event.Sub.Subject},
		}
		return nil
	}
	if err := n.connector.QueueSubscribeMultiple(n.currentSubjects, n.queueGroup, n.eventProcessFn); err != nil {
		return fmt.Errorf("could not start NATS event source: %w", err)
	}
	go func() {
		select {
		case <-ctx.Done():
			n.unsubscribe()
			wg.Done()
			return
		case <-n.quitC:
			n.unsubscribe()
			wg.Done()
			return
		}
	}()
	return nil
}

func (n *NATSEventSource) unsubscribe() {
	if err := n.connector.UnsubscribeAll(); err != nil {
		n.logger.Errorf("Unable to unsubscribe from NATS: %v", err)
	} else {
		n.logger.Debug("Unsubscribed from NATS")
	}
}

func (n *NATSEventSource) OnSubscriptionUpdate(subj []models.EventSubscription) {
	s := dedup(subjects(subj))
	n.logger.Debugf("Updating subscriptions")
	if !isEqual(n.currentSubjects, s) {
		n.logger.Debugf("Cleaning up %d old subscriptions", len(n.currentSubjects))
		err := n.connector.UnsubscribeAll()
		n.logger.Debug("Unsubscribed from previous subscriptions")
		if err != nil {
			n.logger.Errorf("Could not handle subscription update: %v", err)
			return
		}
		n.logger.Debugf("Subscribing to %d topics", len(s))
		if err := n.connector.QueueSubscribeMultiple(s, n.queueGroup, n.eventProcessFn); err != nil {
			n.logger.Errorf("Could not handle subscription update: %v", err)
			return
		}
		n.currentSubjects = s
		n.logger.Debugf("Subscription to %d topics successful", len(s))
	}
}

func (n *NATSEventSource) Sender() types.EventSender {
	return n.connector.Publish
}

func (n *NATSEventSource) Stop() error {
	n.quitC <- struct{}{}
	return nil
}

func (n *NATSEventSource) Cleanup() error {
	return n.connector.Disconnect()
}

func isEqual(a1 []string, a2 []string) bool {
	sort.Strings(a2)
	sort.Strings(a1)
	return reflect.DeepEqual(a1, a2)
}

func dedup(elements []string) []string {
	result := make([]string, 0, len(elements))
	temp := map[string]struct{}{}
	for _, el := range elements {
		if _, ok := temp[el]; !ok {
			temp[el] = struct{}{}
			result = append(result, el)
		}
	}
	return result
}

func subjects(subscriptions []models.EventSubscription) []string {
	var ret []string
	for _, s := range subscriptions {
		ret = append(ret, s.Event)
	}
	return ret
}
