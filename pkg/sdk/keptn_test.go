package sdk

import (
	"context"
	"fmt"
	"github.com/keptn/go-utils/pkg/sdk/internal/config"
	"math"
	"testing"

	"github.com/google/uuid"
	"github.com/keptn/go-utils/pkg/api/models"
	"github.com/keptn/go-utils/pkg/common/strutils"
	"github.com/keptn/go-utils/pkg/lib/v0_2_0"
	"github.com/stretchr/testify/require"
)

func Test_NewKeptn(t *testing.T) {
	t.Run("Create SDK instance", func(t *testing.T) {
		keptnSDK := NewKeptn("my-service")
		require.NotNil(t, keptnSDK)
	})
	t.Run("Create SDK instance with option(s)", func(t *testing.T) {
		myLoggerInstance := newDefaultLogger()
		keptnSDK := NewKeptn("my-service",
			WithAutomaticResponse(true),
			WithLogger(myLoggerInstance),
			WithGracefulShutdown(true),
			WithTaskHandler("event.type", &TaskHandlerMock{}),
			WithTaskHandler("event2.type", &TaskHandlerMock{}))
		require.NotNil(t, keptnSDK)
		require.True(t, keptnSDK.automaticEventResponse)
		require.Equal(t, myLoggerInstance, keptnSDK.Logger())
		require.True(t, keptnSDK.gracefulShutdown)
		require.IsType(t, keptnSDK.taskRegistry.Get("event.type").taskHandler, &TaskHandlerMock{})
		require.IsType(t, keptnSDK.taskRegistry.Get("event2.type").taskHandler, &TaskHandlerMock{})
		require.NotNil(t, keptnSDK.GetResourceHandler())
		require.NotNil(t, keptnSDK.APIV1())
		require.NotNil(t, keptnSDK.APIV2())

	})
}

func Test_ReceivingInvalidEvent(t *testing.T) {
	taskHandler := &TaskHandlerMock{}
	taskHandler.ExecuteFunc = func(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error) { return FakeTaskData{}, nil }
	fakeKeptn := NewFakeKeptn("fake")
	fakeKeptn.AddTaskHandler("sh.keptn.event.faketask.triggered", taskHandler)
	fakeKeptn.NewEvent(models.KeptnContextExtendedCE{
		Data:           math.Inf(1),
		ID:             "id",
		Shkeptncontext: "context",
		Source:         strutils.Stringp("source"),
		Type:           strutils.Stringp("sh.keptn.event.faketask.triggered"),
	})
	fakeKeptn.AssertNumberOfEventSent(t, 1)
}

func Test_SendEvents(t *testing.T) {
	t.Run("Send Started Event", func(t *testing.T) {
		var sentEvent models.KeptnContextExtendedCE
		keptnSDK := NewKeptn("my-service")
		keptnSDK.eventSender = func(ce models.KeptnContextExtendedCE) error {
			sentEvent = ce
			return nil
		}
		err := keptnSDK.SendStartedEvent(KeptnEvent{
			Contenttype:    "application/json",
			Data:           v0_2_0.EventData{Project: "prj", Stage: "stg", Service: "svc"},
			ID:             "id",
			Shkeptncontext: "context",
			Source:         strutils.Stringp("source"),
			Type:           strutils.Stringp("sh.keptn.event.faketask.triggered"),
		})
		require.NoError(t, err)
		require.NotNil(t, sentEvent)
		require.NotEmpty(t, sentEvent.ID)
		require.Equal(t, "sh.keptn.event.faketask.started", *sentEvent.Type)
		require.Equal(t, v0_2_0.EventData{Project: "prj", Stage: "stg", Service: "svc"}, sentEvent.Data)
	})

	t.Run("Send Finished Event", func(t *testing.T) {
		var sentEvent models.KeptnContextExtendedCE
		keptnSDK := NewKeptn("my-service")
		keptnSDK.eventSender = func(ce models.KeptnContextExtendedCE) error {
			sentEvent = ce
			return nil
		}
		err := keptnSDK.SendFinishedEvent(KeptnEvent{
			Contenttype:    "application/json",
			ID:             "id",
			Shkeptncontext: "context",
			Source:         strutils.Stringp("source"),
			Type:           strutils.Stringp("sh.keptn.event.faketask.triggered"),
		}, v0_2_0.EventData{Project: "prj", Stage: "stg", Service: "svc"})
		require.NoError(t, err)
		require.NotNil(t, sentEvent)
		require.NotEmpty(t, sentEvent.ID)
		require.Equal(t, "sh.keptn.event.faketask.finished", *sentEvent.Type)
		require.Equal(t, map[string]interface{}{"project": "prj", "result": "pass", "service": "svc", "stage": "stg", "status": "succeeded"}, sentEvent.Data)
	})

}

func Test_ReceivingEventWithMissingType(t *testing.T) {
	taskHandler := &TaskHandlerMock{}
	taskHandler.ExecuteFunc = func(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error) { return FakeTaskData{}, nil }
	fakeKeptn := NewFakeKeptn("fake")
	fakeKeptn.AddTaskHandler("sh.keptn.event.faketask.triggered", taskHandler)
	fakeKeptn.NewEvent(models.KeptnContextExtendedCE{
		Data:           v0_2_0.EventData{Project: "prj", Stage: "stg", Service: "svc"},
		ID:             "id",
		Shkeptncontext: "context",
		Source:         strutils.Stringp("source"),
	})

	fakeKeptn.AssertNumberOfEventSent(t, 0)
}

func Test_CannotGetEventSenderFromContext(t *testing.T) {
	taskHandler := &TaskHandlerMock{}
	taskHandler.ExecuteFunc = func(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error) { return FakeTaskData{}, nil }
	fakeKeptn := NewFakeKeptn("fake")
	fakeKeptn.AddTaskHandler("sh.keptn.event.faketask.triggered", taskHandler)
	fakeKeptn.Keptn.OnEvent(context.TODO(), models.KeptnContextExtendedCE{
		Data:           v0_2_0.EventData{Project: "prj", Stage: "stg", Service: "svc"},
		ID:             "id",
		Shkeptncontext: "context",
		Source:         strutils.Stringp("source"),
		Type:           strutils.Stringp("sh.keptn.event.faketask.triggered"),
	})
	fakeKeptn.AssertNumberOfEventSent(t, 0)
}

func Test_WhenReceivingAnEvent_StartedEventAndFinishedEventsAreSent(t *testing.T) {
	taskHandler := &TaskHandlerMock{}
	taskHandler.ExecuteFunc = func(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error) { return FakeTaskData{}, nil }
	fakeKeptn := NewFakeKeptn("fake")
	fakeKeptn.AddTaskHandler("sh.keptn.event.faketask.triggered", taskHandler)
	fakeKeptn.NewEvent(models.KeptnContextExtendedCE{
		Data:           v0_2_0.EventData{Project: "prj", Stage: "stg", Service: "svc"},
		ID:             "id",
		Shkeptncontext: "context",
		Source:         strutils.Stringp("source"),
		Type:           strutils.Stringp("sh.keptn.event.faketask.triggered"),
	})

	fakeKeptn.AssertNumberOfEventSent(t, 2)
	fakeKeptn.AssertSentEventType(t, 0, "sh.keptn.event.faketask.started")
	fakeKeptn.AssertSentEventType(t, 1, "sh.keptn.event.faketask.finished")
}

func Test_WhenReceivingAnEvent_TaskHandlerFails(t *testing.T) {
	taskHandler := &TaskHandlerMock{}
	taskHandler.ExecuteFunc = func(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error) {
		return nil, &Error{
			StatusType: v0_2_0.StatusErrored,
			ResultType: v0_2_0.ResultFailed,
			Message:    "something went wrong",
			Err:        fmt.Errorf("something went wrong"),
		}
	}
	fakeKeptn := NewFakeKeptn("fake")
	fakeKeptn.AddTaskHandler("sh.keptn.event.faketask.triggered", taskHandler)
	fakeKeptn.NewEvent(models.KeptnContextExtendedCE{
		Data:           v0_2_0.EventData{Project: "prj", Stage: "stg", Service: "svc"},
		ID:             "id",
		Shkeptncontext: "context",
		Source:         strutils.Stringp("source"),
		Type:           strutils.Stringp("sh.keptn.event.faketask.triggered"),
	})

	fakeKeptn.AssertNumberOfEventSent(t, 2)
	fakeKeptn.AssertSentEventType(t, 0, "sh.keptn.event.faketask.started")
	fakeKeptn.AssertSentEventType(t, 1, "sh.keptn.event.faketask.finished")
	fakeKeptn.AssertSentEventStatus(t, 1, v0_2_0.StatusErrored)
	fakeKeptn.AssertSentEventResult(t, 1, v0_2_0.ResultFailed)
}

func Test_WhenReceivingBadEvent_NoEventIsSent(t *testing.T) {
	taskHandler := &TaskHandlerMock{}
	taskHandler.ExecuteFunc = func(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error) { return FakeTaskData{}, nil }
	fakeKeptn := NewFakeKeptn("fake")
	fakeKeptn.AddTaskHandler("sh.keptn.event.faketask.triggered", taskHandler)
	fakeKeptn.NewEvent(newTestTaskBadTriggeredEvent())
	fakeKeptn.AssertNumberOfEventSent(t, 0)
}

func Test_WhenReceivingAnEvent_AndNoFilterMatches_NoEventIsSent(t *testing.T) {
	taskHandler := &TaskHandlerMock{}
	taskHandler.ExecuteFunc = func(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error) { return FakeTaskData{}, nil }
	fakeKeptn := NewFakeKeptn("fake")
	fakeKeptn.AddTaskHandler("sh.keptn.event.faketask.triggered", taskHandler, func(keptnHandle IKeptn, event KeptnEvent) bool { return false })
	fakeKeptn.NewEvent(models.KeptnContextExtendedCE{
		Data:           v0_2_0.EventData{Project: "prj", Stage: "stg", Service: "svc"},
		ID:             "id",
		Shkeptncontext: "context",
		Source:         strutils.Stringp("source"),
		Type:           strutils.Stringp("sh.keptn.event.faketask.triggered"),
	})

	fakeKeptn.AssertNumberOfEventSent(t, 0)
}

func Test_NoFinishedEventDataProvided(t *testing.T) {
	taskHandler := &TaskHandlerMock{}
	taskHandler.ExecuteFunc = func(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error) {
		return nil, nil
	}
	fakeKeptn := NewFakeKeptn("fake")
	fakeKeptn.AddTaskHandler("sh.keptn.event.faketask.triggered", taskHandler)
	fakeKeptn.NewEvent(models.KeptnContextExtendedCE{
		Data:           v0_2_0.EventData{Project: "prj", Stage: "stg", Service: "svc"},
		ID:             "id",
		Shkeptncontext: "context",
		Source:         strutils.Stringp("source"),
		Type:           strutils.Stringp("sh.keptn.event.faketask.triggered"),
	})

	fakeKeptn.AssertNumberOfEventSent(t, 1)
	fakeKeptn.AssertSentEventType(t, 0, "sh.keptn.event.faketask.started")
}

func Test_InitialRegistrationData(t *testing.T) {
	keptn := Keptn{env: config.EnvConfig{
		PubSubTopic:          "sh.keptn.event.task1.triggered,sh.keptn.event.task2.triggered",
		Location:             "localhost",
		K8sDeploymentVersion: "v1",
		K8sDeploymentName:    "k8s-deployment",
		K8sNamespace:         "k8s-namespace",
		K8sPodName:           "k8s-podname",
		K8sNodeName:          "k8s-nodename",
	}}

	regData := keptn.RegistrationData()
	require.Equal(t, "v1", regData.MetaData.IntegrationVersion)
	require.Equal(t, "localhost", regData.MetaData.Location)
	require.Equal(t, "k8s-deployment", regData.MetaData.KubernetesMetaData.DeploymentName)
	require.Equal(t, "k8s-namespace", regData.MetaData.KubernetesMetaData.Namespace)
	require.Equal(t, "k8s-podname", regData.MetaData.KubernetesMetaData.PodName)
	require.Equal(t, "k8s-nodename", regData.MetaData.Hostname)
	require.Equal(t, []models.EventSubscription{{Event: "sh.keptn.event.task1.triggered"}, {Event: "sh.keptn.event.task2.triggered"}}, regData.Subscriptions)
}

func Test_InitialRegistrationData_EmptyPubSubTopics(t *testing.T) {
	keptn := Keptn{env: config.EnvConfig{PubSubTopic: ""}}
	regData := keptn.RegistrationData()
	require.Equal(t, 0, len(regData.Subscriptions))
}

func newTestTaskTriggeredEvent() models.KeptnContextExtendedCE {
	return models.KeptnContextExtendedCE{
		Contenttype:    "application/json",
		Data:           FakeTaskData{},
		ID:             uuid.New().String(),
		Shkeptncontext: "keptncontext",
		Triggeredid:    "ID",
		GitCommitID:    "mycommitid",
		Source:         strutils.Stringp("unittest"),
		Type:           strutils.Stringp("sh.keptn.event.faketask.triggered"),
	}
}

func newTestTaskBadTriggeredEvent() models.KeptnContextExtendedCE {
	return models.KeptnContextExtendedCE{
		Contenttype:    "application/json",
		Data:           FakeTaskData{},
		ID:             uuid.New().String(),
		Shkeptncontext: "keptncontext",
		Triggeredid:    "ID",
		GitCommitID:    "mycommitid",
		Source:         strutils.Stringp("unittest"),
		Type:           strutils.Stringp("sh.keptn.event.faketask.finished.triggered"),
	}
}

type FakeTaskData struct {
}
type TaskHandlerMock struct {
	// ExecuteFunc mocks the Execute method.
	ExecuteFunc func(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error)
}

func (mock *TaskHandlerMock) Execute(keptnHandle IKeptn, event KeptnEvent) (interface{}, *Error) {
	if mock.ExecuteFunc == nil {
		panic("TaskHandlerMock.ExecuteFunc: method is nil but taskHandler.Execute was just called")
	}
	return mock.ExecuteFunc(keptnHandle, event)
}
