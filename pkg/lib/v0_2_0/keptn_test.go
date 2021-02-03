package v0_2_0

import (
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0/fake"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_ensureContextAttributesAreSet(t *testing.T) {
	type args struct {
		srcEvent keptn.EventProperties
		newEvent keptn.EventProperties
	}
	tests := []struct {
		name string
		args args
		want keptn.EventProperties
	}{
		{
			name: "copy context attributes to empty event",
			args: args{
				srcEvent: &EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
					Labels: map[string]string{
						"foo": "bar",
					},
				},
				newEvent: &EventData{},
			},
			want: &EventData{
				Project: "my-project",
				Stage:   "my-stage",
				Service: "my-service",
				Labels: map[string]string{
					"foo": "bar",
				},
			},
		},
		{
			name: "add new labels to event",
			args: args{
				srcEvent: &EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
					Labels: map[string]string{
						"foo": "bar",
					},
				},
				newEvent: &EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
					Labels: map[string]string{
						"bar": "foo",
					},
				},
			},
			want: &EventData{
				Project: "my-project",
				Stage:   "my-stage",
				Service: "my-service",
				Labels: map[string]string{
					"foo": "bar",
					"bar": "foo",
				},
			},
		},
		{
			name: "merge labels - do not overwrite existing ones",
			args: args{
				srcEvent: &EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
					Labels: map[string]string{
						"foo": "bar",
					},
				},
				newEvent: &EventData{
					Project: "my-project",
					Stage:   "my-stage",
					Service: "my-service",
					Labels: map[string]string{
						"foo": "foo",
						"bar": "foo",
					},
				},
			},
			want: &EventData{
				Project: "my-project",
				Stage:   "my-stage",
				Service: "my-service",
				Labels: map[string]string{
					"foo": "bar",
					"bar": "foo",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ensureContextAttributesAreSet(tt.args.srcEvent, tt.args.newEvent)

			assert.EqualValues(t, tt.want, got)
		})
	}
}

func TestKeptn_SendEventConvenienceFunctions(t *testing.T) {

	inputEvent := cloudevents.NewEvent()
	inputEvent.SetType(GetTriggeredEventType(EvaluationTaskName))
	inputEvent.SetExtension(keptnContextCEExtension, "my-context")
	inputEvent.SetID("my-triggered-id")
	inputEvent.SetDataContentType(cloudevents.ApplicationJSON)
	inputEvent.SetData(cloudevents.ApplicationJSON, &EventData{
		Project: "my-project",
		Stage:   "my-stage",
		Service: "my-service",
		Labels: map[string]string{
			"foo": "bar",
		},
	})

	type fields struct {
		KeptnBase keptn.KeptnBase
	}
	type args struct {
		data   keptn.EventProperties
		source string
	}
	type wantEventProperties struct {
		eventType    string
		keptnContext string
		triggeredID  string
		eventData    keptn.EventProperties
	}
	tests := []struct {
		name          string
		fields        fields
		args          args
		sendEventType string
		wantErr       bool
		wantEvents    []wantEventProperties
	}{
		{
			name: "send started event - no error",
			fields: fields{
				KeptnBase: keptn.KeptnBase{
					Event:       nil,
					CloudEvent:  nil,
					EventSender: nil,
				},
			},
			args: args{
				source: "my-source",
			},
			sendEventType: keptnStartedEventSuffix,
			wantErr:       false,
			wantEvents: []wantEventProperties{
				{
					eventType:    GetStartedEventType(EvaluationTaskName),
					keptnContext: "my-context",
					triggeredID:  "my-triggered-id",
					eventData: &EventData{
						Project: "my-project",
						Stage:   "my-stage",
						Service: "my-service",
						Labels: map[string]string{
							"foo": "bar",
						},
					},
				},
			},
		},
		{
			name: "send status.changed event - no error",
			fields: fields{
				KeptnBase: keptn.KeptnBase{
					Event:       nil,
					CloudEvent:  nil,
					EventSender: nil,
				},
			},
			args: args{
				source: "my-source",
			},
			sendEventType: keptnStatusChangedEventSuffix,
			wantErr:       false,
			wantEvents: []wantEventProperties{
				{
					eventType:    GetStatusChangedEventType(EvaluationTaskName),
					keptnContext: "my-context",
					triggeredID:  "my-triggered-id",
					eventData: &EventData{
						Project: "my-project",
						Stage:   "my-stage",
						Service: "my-service",
						Labels: map[string]string{
							"foo": "bar",
						},
					},
				},
			},
		},
		{
			name: "send finished event - no error",
			fields: fields{
				KeptnBase: keptn.KeptnBase{
					Event:       nil,
					CloudEvent:  nil,
					EventSender: nil,
				},
			},
			args: args{
				source: "my-source",
			},
			sendEventType: keptnFinishedEventSuffix,
			wantErr:       false,
			wantEvents: []wantEventProperties{
				{
					eventType:    GetFinishedEventType(EvaluationTaskName),
					keptnContext: "my-context",
					triggeredID:  "my-triggered-id",
					eventData: &EventData{
						Project: "my-project",
						Stage:   "my-stage",
						Service: "my-service",
						Labels: map[string]string{
							"foo": "bar",
						},
					},
				},
			},
		},
		{
			name: "send finished event with additional attributes- no error",
			fields: fields{
				KeptnBase: keptn.KeptnBase{
					Event:       nil,
					CloudEvent:  nil,
					EventSender: nil,
				},
			},
			args: args{
				data: &EventData{
					Result: ResultPass,
				},
				source: "my-source",
			},
			sendEventType: keptnFinishedEventSuffix,
			wantErr:       false,
			wantEvents: []wantEventProperties{
				{
					eventType:    GetFinishedEventType(EvaluationTaskName),
					keptnContext: "my-context",
					triggeredID:  "my-triggered-id",
					eventData: &EventData{
						Project: "my-project",
						Stage:   "my-stage",
						Service: "my-service",
						Labels: map[string]string{
							"foo": "bar",
						},
						Result: ResultPass,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fakeSender := &fake.EventSender{}
			k, err := NewKeptn(&inputEvent, keptn.KeptnOpts{
				EventSender: fakeSender,
			})
			if err != nil {
				t.Error(err.Error())
			}
			var got string

			switch tt.sendEventType {
			case keptnStartedEventSuffix:
				got, err = k.SendTaskStartedEvent(tt.args.data, tt.args.source)
			case keptnStatusChangedEventSuffix:
				got, err = k.SendTaskStatusChangedEvent(tt.args.data, tt.args.source)
			case keptnFinishedEventSuffix:
				got, err = k.SendTaskFinishedEvent(tt.args.data, tt.args.source)
			default:
				return
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == "" {
				t.Errorf("did not return event ID")
			}

			if len(fakeSender.SentEvents) != len(tt.wantEvents) {
				t.Errorf("did not receive expected number of events. Expected %d, got %d", len(tt.wantEvents), len(fakeSender.SentEvents))
			}

			for index, event := range fakeSender.SentEvents {
				assert.Equal(t, event.Type(), tt.wantEvents[index].eventType)
				triggeredID, err := event.Context.GetExtension(triggeredIDCEExtenstion)
				assert.Nil(t, err)
				assert.Equal(t, triggeredID.(string), tt.wantEvents[index].triggeredID)
				keptnContext, err := event.Context.GetExtension(keptnContextCEExtension)
				assert.Nil(t, err)
				assert.Equal(t, keptnContext.(string), tt.wantEvents[index].keptnContext)
				data := &EventData{}
				err = event.DataAs(data)
				assert.Nil(t, err)
				assert.EqualValues(t, data, tt.wantEvents[index].eventData)
			}
		})
	}
}
