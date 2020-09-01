package v0_2_0

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/keptn/go-utils/pkg/lib/keptn"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

const MAX_SEND_RETRIES = 3

// GetTriggeredEventType returns for the given task the name of the triggered event type
func GetTriggeredEventType(task string) string {
	return "sh.keptn.event." + task + ".triggered"
}

// GetStartedEventType returns for the given task the name of the started event type
func GetStartedEventType(task string) string {
	return "sh.keptn.event." + task + ".started"
}

// GetStatusChangedEventType returns for the given task the name of the status.changed event type
func GetStatusChangedEventType(task string) string {
	return "sh.keptn.event." + task + ".status.changed"
}

// GetFinishedEventType returns for the given task the name of the finished event type
func GetFinishedEventType(task string) string {
	return "sh.keptn.event." + task + ".finished"
}

// EventData contains mandatory fields of all Keptn CloudEvents
type EventData struct {
	Project string            `json:"project"`
	Stage   string            `json:"stage"`
	Service string            `json:"service"`
	Labels  map[string]string `json:"labels"`

	Status  StatusType `json:"status"`
	Result  ResultType `json:"result"`
	Message string     `json:"message"`
}

func (e EventData) GetProject() string {
	return e.Project
}

func (e EventData) GetStage() string {
	return e.Stage
}

func (e EventData) GetService() string {
	return e.Service
}

func (e EventData) GetLabels() map[string]string {
	return e.Labels
}

// SendCloudEvent sends a cloudevent to the event broker
func (k *Keptn) SendCloudEvent(event cloudevents.Event) error {
	if k.UseLocalFileSystem {
		log.Println(fmt.Printf("%v", event.Data()))
		return nil
	}

	c, err := cloudevents.NewDefaultClient()
	if err != nil {
		return fmt.Errorf("failed to create client, %v", err)
	}

	// Set a target.
	ctx := cloudevents.ContextWithTarget(context.Background(), k.EventBrokerURL)

	for i := 0; i <= MAX_SEND_RETRIES; i++ {
		result := c.Send(ctx, event)
		if cloudevents.IsACK(result) {
			return nil
		}
		<-time.After(keptn.GetExpBackoffTime(i + 1))
	}
	return errors.New("Failed to send cloudevent:, " + err.Error())
}
