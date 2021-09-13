package observability

import (
	"context"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// CloudEventTraceContext a wrapper around the OpenTelemetry TraceContext
// https://github.com/open-telemetry/opentelemetry-go/blob/main/propagation/trace_context.go
type CloudEventTraceContext struct {
	traceContext propagation.TraceContext
}

// NewCloudEventTraceContext creates a new CloudEventTraceContext
func NewCloudEventTraceContext() CloudEventTraceContext {
	return CloudEventTraceContext{traceContext: propagation.TraceContext{}}
}

// Extract extracts the tracecontext from the cloud event into the context.
//
// If the context has a recording span, then the same context is returned. If not, then the extraction
// from the cloud event happens. The reason for this is to avoid breaking the span order in the trace.
// For instrumented clients, the context *should* have the incoming span from the auto-instrumented library
// thus using this one is more appropriate.
func (etc CloudEventTraceContext) Extract(ctx context.Context, carrier CloudEventCarrier) context.Context {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		// if the context already has an active span we just return that
		return ctx
	}

	// Extracts the traceparent from the cloud event into the context
	// This is useful when there's no context (reading from the queue in a long running process)
	// In this case we can use the traceparent from the event to continue the trace flow.
	return etc.traceContext.Extract(ctx, carrier)
}

// Inject injects the current tracecontext from the context into the cloud event
func (etc CloudEventTraceContext) Inject(ctx context.Context, carrier CloudEventCarrier) {
	etc.traceContext.Inject(ctx, carrier)
}
