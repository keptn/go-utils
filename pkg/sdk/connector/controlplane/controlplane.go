package controlplane

import (
	"context"
	"errors"
	"fmt"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/sdk/connector/eventmatcher"
	"github.com/keptn/go-utils/pkg/sdk/connector/eventsource"
	"github.com/keptn/go-utils/pkg/sdk/connector/logforwarder"
	"github.com/keptn/go-utils/pkg/sdk/connector/logger"
	"github.com/keptn/go-utils/pkg/sdk/connector/subscriptionsource"
	"github.com/keptn/go-utils/pkg/sdk/connector/types"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type EventSender = types.EventSender
type EventSenderKeyType = types.EventSenderKeyType
type RegistrationData = types.RegistrationData

const tmpDataDistributorKey = "distributor"

var ErrEventHandleFatal = errors.New("fatal event handling error")

// Integration represents a Keptn Service that wants to receive events from the Keptn Control plane
type Integration interface {
	// OnEvent is called when a new event was received
	OnEvent(context.Context, models.KeptnContextExtendedCE) error

	// RegistrationData is called to get the initial registration data
	RegistrationData() types.RegistrationData
}

// ControlPlane can be used to connect to the Keptn Control Plane
type ControlPlane struct {
	subscriptionSource    subscriptionsource.SubscriptionSource
	eventSource           eventsource.EventSource
	currentSubscriptions  []models.EventSubscription
	logger                logger.Logger
	registered            bool
	integrationID         string
	logForwarder          logforwarder.LogForwarder
	mtx                   *sync.RWMutex
	eventHandlerWaitGroup *sync.WaitGroup
}

// WithLogger sets the logger to use
func WithLogger(logger logger.Logger) func(plane *ControlPlane) {
	return func(ns *ControlPlane) {
		ns.logger = logger
	}
}

// RunWithGracefulShutdown starts the controlplane component which takes care of registering
// the integration and handling events and subscriptions. Further, it supports graceful shutdown handling
// when receiving a SIGHUB, SIGINT, SIGQUIT, SIGARBT or SIGTERM signal.
//
// This call is blocking.
//
//If you want to start the controlPlane component with an own context you need to call the Register(ctx,integration)
// method on your own
func RunWithGracefulShutdown(controlPlane *ControlPlane, integration Integration, shutdownTimeout time.Duration) error {
	ctxShutdown, cancel := context.WithCancel(context.Background())
	defer cancel()

	ctxShutdown, _ = signal.NotifyContext(ctxShutdown, syscall.SIGHUP, syscall.SIGINT, syscall.SIGQUIT, syscall.SIGABRT, syscall.SIGTERM)
	go func() {
		<-ctxShutdown.Done()
		time.Sleep(shutdownTimeout) // shutdown timeout
		log.Printf("failed to gracefully shutdown")
		os.Exit(1)
	}()

	return controlPlane.Register(ctxShutdown, integration)

}

// New creates a new ControlPlane
// It is using a SubscriptionSource source to get information about current uniform subscriptions
// as well as an EventSource to actually receive events from Keptn
// and a LogForwarder to forward error logs
func New(subscriptionSource subscriptionsource.SubscriptionSource, eventSource eventsource.EventSource, logForwarder logforwarder.LogForwarder, opts ...func(plane *ControlPlane)) *ControlPlane {
	cp := &ControlPlane{
		subscriptionSource:    subscriptionSource,
		eventSource:           eventSource,
		currentSubscriptions:  []models.EventSubscription{},
		logger:                logger.NewDefaultLogger(),
		logForwarder:          logForwarder,
		registered:            false,
		mtx:                   &sync.RWMutex{},
		eventHandlerWaitGroup: &sync.WaitGroup{},
	}
	for _, o := range opts {
		o(cp)
	}
	return cp
}

// Register is initially used to register the Keptn integration to the Control Plane
func (cp *ControlPlane) Register(ctx context.Context, integration Integration) error {
	eventUpdates := make(chan types.EventUpdate)
	subscriptionUpdates := make(chan []models.EventSubscription)
	errC := make(chan error, 1)

	var err error
	registrationData := integration.RegistrationData()
	cp.logger.Debugf("Registering integration %s", integration.RegistrationData().Name)
	cp.integrationID, err = cp.subscriptionSource.Register(models.Integration(registrationData))
	if err != nil {
		return fmt.Errorf("could not register integration: %w", err)
	}
	cp.logger.Debugf("Registered with integration ID %s", cp.integrationID)
	registrationData.ID = cp.integrationID

	// WaitGroup used for synchronized shutdown of eventsource and subscription source
	// during cancellation of the context
	wg := &sync.WaitGroup{}
	wg.Add(2)

	cp.logger.Debugf("Starting event source for integration ID %s", cp.integrationID)
	if err := cp.eventSource.Start(ctx, registrationData, eventUpdates, errC, wg); err != nil {
		return err
	}
	cp.logger.Debugf("Event source started with data: %+v", registrationData)
	cp.logger.Debugf("Starting subscription source for integration ID %s", cp.integrationID)
	if err := cp.subscriptionSource.Start(ctx, registrationData, subscriptionUpdates, errC, wg); err != nil {
		return err
	}
	cp.logger.Debug("Subscription source started")
	cp.setRegistrationStatus(true)
	for {
		select {
		// event updates
		case event := <-eventUpdates:
			cp.logger.Debug("Got new event update")
			err := cp.handle(ctx, event, integration)
			if errors.Is(err, ErrEventHandleFatal) {
				return err
			}

		// subscription updates
		case subscriptions := <-subscriptionUpdates:
			cp.logger.Debugf("ControlPlane: Got a subscription update with %d subscriptions", len(subscriptions))
			cp.currentSubscriptions = subscriptions
			cp.eventSource.OnSubscriptionUpdate(subscriptions)

		// control plane cancelled via context
		case <-ctx.Done():
			cp.logger.Info("ControlPlane cancelled via context. Unregistering...")
			cp.stopComponents()
			wg.Wait()
			cp.waitForEventHandlers()
			cp.cleanup()
			cp.setRegistrationStatus(false)
			return nil

		// control plane cancelled via error in either one of the sub components
		case e := <-errC:
			cp.logger.Errorf("Stopping control plane due to error: %v", e)
			cp.logger.Info("Waiting for components to shutdown")
			cp.stopComponents()
			wg.Wait()
			cp.waitForEventHandlers()
			cp.cleanup()
			cp.setRegistrationStatus(false)
			return nil
		}
	}
}

func (cp *ControlPlane) waitForEventHandlers() {
	cp.logger.Info("Wait for all event handlers to finish")
	cp.eventHandlerWaitGroup.Wait()
	cp.logger.Info("All event handlers done - ready to shut down")
}

// IsRegistered can be called to detect whether the controlPlane is registered and ready to receive events
func (cp *ControlPlane) IsRegistered() bool {
	cp.mtx.RLock()
	defer cp.mtx.RUnlock()
	return cp.registered
}

func (cp *ControlPlane) stopComponents() {
	cp.logger.Info("Stopping subscription source...")
	if err := cp.subscriptionSource.Stop(); err != nil {
		log.Fatalf("Unable to stop subscription source: %v", err)
	}
	cp.logger.Info("Stopping event source...")
	if err := cp.eventSource.Stop(); err != nil {
		log.Fatalf("Unable to stop event source: %v", err)
	}
}

func (cp *ControlPlane) handle(ctx context.Context, eventUpdate types.EventUpdate, integration Integration) error {
	cp.logger.Debugf("Received an event of type: %s", *eventUpdate.KeptnEvent.Type)

	// if we already know the subscription ID we can just forward the event to be handled
	if eventUpdate.SubscriptionID != "" {
		if err := cp.forwardMatchedEvent(ctx, eventUpdate, integration, eventUpdate.SubscriptionID); err != nil {
			return err
		}
	} else {
		for _, subscription := range cp.currentSubscriptions {
			if subscription.Event == eventUpdate.MetaData.Subject {
				cp.logger.Debugf("Check if event matches subscription %s", subscription.ID)
				matcher := eventmatcher.New(subscription)
				if matcher.Matches(eventUpdate.KeptnEvent) {
					cp.logger.Info("Forwarding matched event update: ", eventUpdate.KeptnEvent.ID)
					if err := cp.forwardMatchedEvent(ctx, eventUpdate, integration, subscription.ID); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}

func (cp *ControlPlane) getSender(sender types.EventSender) types.EventSender {
	if cp.logForwarder != nil {
		return func(ce models.KeptnContextExtendedCE) error {
			err := cp.logForwarder.Forward(ce, cp.integrationID)
			if err != nil {
				cp.logger.Warnf("could not forward event")
			}
			return sender(ce)
		}
	} else {
		return sender
	}
}

func (cp *ControlPlane) forwardMatchedEvent(ctx context.Context, eventUpdate types.EventUpdate, integration Integration, subscriptionID string) error {
	// increase the eventHandler WaitGroup
	cp.eventHandlerWaitGroup.Add(1)
	// when the event handler is done, decrease the WaitGroup again
	defer cp.eventHandlerWaitGroup.Done()

	err := eventUpdate.KeptnEvent.AddTemporaryData(
		tmpDataDistributorKey,
		types.AdditionalSubscriptionData{
			SubscriptionID: subscriptionID,
		},
		models.AddTemporaryDataOptions{
			OverwriteIfExisting: true,
		},
	)
	if err != nil {
		cp.logger.Warnf("Could not append subscription data to event: %v", err)
	}
	if err := integration.OnEvent(context.WithValue(ctx, types.EventSenderKey, cp.getSender(cp.eventSource.Sender())), eventUpdate.KeptnEvent); err != nil {
		if errors.Is(err, ErrEventHandleFatal) {
			cp.logger.Errorf("Fatal error during handling of event: %v", err)
			return err
		}
		cp.logger.Warnf("Error during handling of event: %v", err)
	}
	return nil
}

func (cp *ControlPlane) setRegistrationStatus(registered bool) {
	cp.mtx.Lock()
	defer cp.mtx.Unlock()
	cp.registered = registered
}

func (cp *ControlPlane) cleanup() {
	cp.logger.Info("Cleaning up event source...")
	if err := cp.eventSource.Cleanup(); err != nil {
		log.Fatalf("Unable to clean up event source: %v", err)
	}
}
