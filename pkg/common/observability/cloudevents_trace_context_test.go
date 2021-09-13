package observability

import (
	"context"
	"net/http"
	"testing"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/extensions"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/oteltest"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

type testcase struct {
	name                   string
	event                  cloudevents.Event
	header                 http.Header
	contextHasTraceContext bool
	want                   extensions.DistributedTracingExtension
}

var (
	traceparent = http.CanonicalHeaderKey("traceparent")
	tracestate  = http.CanonicalHeaderKey("tracestate")

	prop           = propagation.TraceContext{}
	eventTraceID   = "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	eventSpanID    = "bbbbbbbbbbbbbbbb"
	distributedExt = extensions.DistributedTracingExtension{
		TraceParent: "00-bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb-bbbbbbbbbbbbbbbb-00",
		TraceState:  "key1=value1,key2=value2",
	}
)

func TestExtractTraceContextFromEvent(t *testing.T) {
	sr := new(oteltest.SpanRecorder)
	provider := oteltest.NewTracerProvider(oteltest.WithSpanRecorder(sr))
	otel.SetTracerProvider(provider)
	tracer := provider.Tracer("test-tracer")

	tests := []testcase{
		{
			name:  "context with recording span",
			event: createCloudEvent(distributedExt),
			header: http.Header{
				traceparent: []string{"00-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-aaaaaaaaaaaaaaaa-00"},
			},
			contextHasTraceContext: true,
		},
		{
			name:                   "context without tracecontext",
			event:                  createCloudEvent(distributedExt),
			contextHasTraceContext: false,
		},
		{
			name:  "context with tracecontext and event with invalid tracecontext",
			event: createCloudEventWithInvalidTraceParent(),
			header: http.Header{
				traceparent: []string{"00-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-aaaaaaaaaaaaaaaa-00"},
			},
			contextHasTraceContext: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			expectedCtx := context.Background()

			if tc.contextHasTraceContext {
				// Simulates a case of a auto-instrumented client where the context
				// has the incoming parent span + the new span started by the auto-instrumented library (e.g Http)
				expectedCtx = prop.Extract(expectedCtx, propagation.HeaderCarrier(tc.header))
				expectedCtx, span := tracer.Start(expectedCtx, "http-autoinstrumentation")

				// act
				actualCtx := ExtractDistributedTracingExtension(expectedCtx, tc.event)

				// Because the ctx already had a traceContext, the new should be the same
				assert.Equal(t, trace.SpanContextFromContext(expectedCtx), trace.SpanContextFromContext(actualCtx))
				span.End()
			} else {

				// act
				actualCtx := ExtractDistributedTracingExtension(expectedCtx, tc.event)

				// the new context was enriched with the traceparent from the event
				assert.NotEqual(t, trace.SpanContextFromContext(expectedCtx), trace.SpanContextFromContext(actualCtx))

				sc := trace.SpanContextFromContext(actualCtx)

				// tracecontext should be the same as in the event
				assert.Equal(t, eventTraceID, sc.TraceID().String())
				assert.Equal(t, eventSpanID, sc.SpanID().String())
				assert.Equal(t, eventSpanID, sc.SpanID().String())
				assert.Equal(t, distributedExt.TraceState, sc.TraceState().String())
			}
		})
	}
}

func TestInjectDistributedTracingExtension(t *testing.T) {

	sr := new(oteltest.SpanRecorder)
	provider := oteltest.NewTracerProvider(oteltest.WithSpanRecorder(sr))
	otel.SetTracerProvider(provider)

	tests := []testcase{
		{
			name:  "inject tracecontext into event",
			event: createCloudEvent(extensions.DistributedTracingExtension{}),
			header: http.Header{
				traceparent: []string{"00-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-aaaaaaaaaaaaaaaa-00"},
				tracestate:  []string{"key1=value1,key2=value2"},
			},
			want: extensions.DistributedTracingExtension{
				TraceParent: "00-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-aaaaaaaaaaaaaaaa-00",
				TraceState:  "key1=value1,key2=value2",
			},
		},
		{
			name:  "ovewrite tracecontext in the event",
			event: createCloudEvent(distributedExt),
			header: http.Header{
				traceparent: []string{"00-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-aaaaaaaaaaaaaaaa-00"},
				tracestate:  []string{"key1=value1,key2=value2,key3=value3"},
			},
			want: extensions.DistributedTracingExtension{
				TraceParent: "00-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-aaaaaaaaaaaaaaaa-00",
				TraceState:  "key1=value1,key2=value2,key3=value3",
			},
		},
		{
			name:  "context without tracecontext",
			event: createCloudEvent(distributedExt),
			want:  distributedExt,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			ctx := context.Background()
			ctx = prop.Extract(ctx, propagation.HeaderCarrier(tc.header))

			// act
			InjectDistributedTracingExtension(ctx, tc.event)

			actual, ok := extensions.GetDistributedTracingExtension(tc.event)
			assert.True(t, ok)

			assert.Equal(t, tc.want, actual)
		})
	}

}

func createCloudEvent(distributedExt extensions.DistributedTracingExtension) cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSource("example/uri")
	event.SetType("example.type")
	event.SetData(cloudevents.ApplicationJSON, map[string]string{"hello": "world"})

	if distributedExt.TraceParent != "" {
		distributedExt.AddTracingAttributes(&event)
	}

	return event
}

func createCloudEventWithInvalidTraceParent() cloudevents.Event {
	event := cloudevents.NewEvent()
	event.SetSource("example/uri")
	event.SetType("example.type")
	event.SetData(cloudevents.ApplicationJSON, map[string]string{"hello": "world"})

	// set directly to force an invalid value
	event.SetExtension(extensions.TraceParentExtension, 123)

	return event
}
