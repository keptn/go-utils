package v0_2_0

import (
	"context"
	"errors"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/keptn/go-utils/pkg/lib/keptn"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0/fake"
	"github.com/stretchr/testify/assert"
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
	inputEvent.SetExtension(keptnSpecVersionCEExtension, "0.2.0")
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
		eventType        string
		keptnContext     string
		triggeredID      string
		keptnSpecVersion string
		eventData        keptn.EventProperties
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
					EventSender: &fake.EventSender{},
				},
			},
			args: args{
				source: "my-source",
			},
			sendEventType: keptnStartedEventSuffix,
			wantErr:       false,
			wantEvents: []wantEventProperties{
				{
					eventType:        GetStartedEventType(EvaluationTaskName),
					keptnContext:     "my-context",
					triggeredID:      "my-triggered-id",
					keptnSpecVersion: "0.2.0",
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
					EventSender: &fake.EventSender{},
				},
			},
			args: args{
				source: "my-source",
			},
			sendEventType: keptnStatusChangedEventSuffix,
			wantErr:       false,
			wantEvents: []wantEventProperties{
				{
					eventType:        GetStatusChangedEventType(EvaluationTaskName),
					keptnContext:     "my-context",
					triggeredID:      "my-triggered-id",
					keptnSpecVersion: "0.2.0",
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
					EventSender: &fake.EventSender{},
				},
			},
			args: args{
				source: "my-source",
			},
			sendEventType: keptnFinishedEventSuffix,
			wantErr:       false,
			wantEvents: []wantEventProperties{
				{
					eventType:        GetFinishedEventType(EvaluationTaskName),
					keptnContext:     "my-context",
					triggeredID:      "my-triggered-id",
					keptnSpecVersion: "0.2.0",
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
					EventSender: &fake.EventSender{},
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
					eventType:        GetFinishedEventType(EvaluationTaskName),
					keptnContext:     "my-context",
					triggeredID:      "my-triggered-id",
					keptnSpecVersion: "0.2.0",
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
		{
			name: "send finished event - error when sending event",
			fields: fields{
				KeptnBase: keptn.KeptnBase{
					EventSender: &fake.EventSender{
						Reactors: map[string]func(event cloudevents.Event) error{
							"*": func(event cloudevents.Event) error {
								return errors.New("")
							},
						},
					},
				},
			},
			args: args{
				data: &EventData{
					Result: ResultPass,
				},
				source: "my-source",
			},
			sendEventType: keptnFinishedEventSuffix,
			wantErr:       true,
			wantEvents:    []wantEventProperties{},
		},
		{
			name: "send event without source - return error",
			fields: fields{
				KeptnBase: keptn.KeptnBase{
					EventSender: &fake.EventSender{
						Reactors: map[string]func(event cloudevents.Event) error{
							"*": func(event cloudevents.Event) error {
								return errors.New("")
							},
						},
					},
				},
			},
			args: args{
				data: &EventData{
					Result: ResultPass,
				},
				source: "",
			},
			sendEventType: keptnFinishedEventSuffix,
			wantErr:       true,
			wantEvents:    []wantEventProperties{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k, err := NewKeptn(&inputEvent, keptn.KeptnOpts{
				EventSender: tt.fields.KeptnBase.EventSender,
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
			if !tt.wantErr && got == "" {
				t.Errorf("did not return event ID")
			}

			if len(k.EventSender.(*fake.EventSender).SentEvents) != len(tt.wantEvents) {
				t.Errorf("did not receive expected number of events. Expected %d, got %d", len(tt.wantEvents), len(k.EventSender.(*fake.EventSender).SentEvents))
			}

			for index, event := range k.EventSender.(*fake.EventSender).SentEvents {
				assert.Equal(t, event.Type(), tt.wantEvents[index].eventType)
				triggeredID, err := event.Context.GetExtension(triggeredIDCEExtension)
				assert.Nil(t, err)
				assert.Equal(t, triggeredID.(string), tt.wantEvents[index].triggeredID)
				keptnContext, err := event.Context.GetExtension(keptnContextCEExtension)
				assert.Nil(t, err)
				assert.Equal(t, keptnContext.(string), tt.wantEvents[index].keptnContext)
				keptnSpecVersion, err := event.Context.GetExtension(keptnSpecVersionCEExtension)
				assert.Nil(t, err)
				assert.Equal(t, keptnSpecVersion.(string), tt.wantEvents[index].keptnSpecVersion)
				data := &EventData{}
				err = event.DataAs(data)
				assert.Nil(t, err)
				assert.EqualValues(t, data, tt.wantEvents[index].eventData)
			}
		})
	}
}

func TestKeptn_EnsureGoContextIsSet(t *testing.T) {
	event := cloudevents.NewEvent()
	event.SetType("test-type")
	event.SetSource("test-source")

	testCases := []struct {
		name    string
		ctx     context.Context
		wantCtx context.Context
	}{
		{
			name:    "no provided context",
			ctx:     nil,
			wantCtx: context.Background(),
		},
		{
			name:    "with provided enriched context",
			ctx:     cloudevents.WithEncodingStructured(context.Background()),
			wantCtx: cloudevents.WithEncodingStructured(context.Background()),
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			k, err := NewKeptn(&event, keptn.KeptnOpts{
				Context: tc.ctx,
			})

			if err != nil {
				t.Error(err.Error())
			}
			assert.Equal(t, tc.wantCtx, k.Context)
		})
	}
}
