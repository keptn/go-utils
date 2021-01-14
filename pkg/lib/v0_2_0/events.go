package v0_2_0

import (
	"context"
	"errors"
	"fmt"
	"github.com/cloudevents/sdk-go/v2/protocol"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/mitchellh/mapstructure"
	"log"
	"time"

	"github.com/keptn/go-utils/pkg/lib/keptn"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	httpprotocol "github.com/cloudevents/sdk-go/v2/protocol/http"
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
	Project string            `json:"project,omitempty"`
	Stage   string            `json:"stage,omitempty"`
	Service string            `json:"service,omitempty"`
	Labels  map[string]string `json:"labels,omitempty"`

	Status  StatusType `json:"status,omitempty" jsonschema:"enum=succeeded,enum=errored,enum=unknown"`
	Result  ResultType `json:"result,omitempty" jsonschema:"enum=pass,enum=warning,enum=fail"`
	Message string     `json:"message,omitempty"`
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
		log.Println(fmt.Printf("%v", string(event.Data())))
		return nil
	}

	ctx := cloudevents.ContextWithTarget(context.Background(), k.EventBrokerURL)
	ctx = cloudevents.WithEncodingStructured(ctx)

	p, err := cloudevents.NewHTTP()
	if err != nil {
		log.Fatalf("failed to create protocol: %s", err.Error())
	}

	c, err := cloudevents.NewClient(p, cloudevents.WithTimeNow(), cloudevents.WithUUIDs())
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}

	var result protocol.Result
	for i := 0; i <= MAX_SEND_RETRIES; i++ {
		result = c.Send(ctx, event)
		httpResult, ok := result.(*httpprotocol.Result)
		if ok {
			if httpResult.StatusCode >= 200 && httpResult.StatusCode < 300 {
				return nil
			} else {
				<-time.After(keptn.GetExpBackoffTime(i + 1))
			}
		} else if cloudevents.IsUndelivered(result) {
			<-time.After(keptn.GetExpBackoffTime(i + 1))
		} else {
			return nil
		}
	}
	return errors.New("Failed to send cloudevent: " + result.Error())
}

// Decode decodes the given raw interface to the target pointer specified
// by the out paratmeter
func Decode(in, out interface{}) error {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		Squash: true,
		Result: out,
	})
	if err != nil {
		return err
	}
	return decoder.Decode(in)
}

// EventDataAs decodes the event data of the given keptn cloud event to the
// target pointer specified by the out parameter
func EventDataAs(in models.KeptnContextExtendedCE, out interface{}) error {
	return Decode(in.Data, out)
}
