package observability

import (
	"context"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/extensions"
)

// CloudEventCarrier adapts the distributed trace extension to satisfy the TextMapCarrier interface.
// https://github.com/open-telemetry/opentelemetry-go/blob/4bfa00347d52e655899b13054198c0aefbd67f85/propagation/propagation.go#L23
type CloudEventCarrier struct {
	Extension *extensions.DistributedTracingExtension
}

// NewCloudEventCarrier creates a new CloudEventCarrier with an empty distributed tracing extension
func NewCloudEventCarrier() CloudEventCarrier {
	return CloudEventCarrier{Extension: &extensions.DistributedTracingExtension{}}
}

// NewCloudEventCarrierWithEvent creates a new CloudEventCarrier with a distributed tracing extension
// populated with the trace data from the event
func NewCloudEventCarrierWithEvent(event cloudevents.Event) CloudEventCarrier {
	var te, ok = extensions.GetDistributedTracingExtension(event)
	if !ok {
		return CloudEventCarrier{Extension: &extensions.DistributedTracingExtension{}}
	}
	return CloudEventCarrier{Extension: &te}
}

// Get returns the value associated with the passed key.
func (cec CloudEventCarrier) Get(key string) string {
	switch key {
	case extensions.TraceParentExtension:
		return cec.Extension.TraceParent
	case extensions.TraceStateExtension:
		return cec.Extension.TraceState
	default:
		return ""
	}
}

// Set stores the key-value pair.
func (cec CloudEventCarrier) Set(key string, value string) {
	switch key {
	case extensions.TraceParentExtension:
		cec.Extension.TraceParent = value
	case extensions.TraceStateExtension:
		cec.Extension.TraceState = value
	}
}

// Keys lists the keys stored in this carrier.
func (cec CloudEventCarrier) Keys() []string {
	return []string{extensions.TraceParentExtension, extensions.TraceStateExtension}
}

// InjectDistributedTracingExtension injects the tracecontext into the event as a DistributedTracingExtension
func InjectDistributedTracingExtension(ctx context.Context, event cloudevents.Event) {

	// TODO: Should we validate if there's already a tracecontext in the event?
	// Calling it will override any existing value..
	tc := NewCloudEventTraceContext()
	carrier := NewCloudEventCarrier()
	tc.Inject(ctx, carrier)
	carrier.Extension.AddTracingAttributes(&event)
}

// ExtractDistributedTracingExtension reads tracecontext from the cloudevent DistributedTracingExtension into a returned Context.
//
// The returned Context will be a copy of ctx and contain the extracted
// tracecontext as the remote SpanContext. If the extracted tracecontext is
// invalid, the passed ctx will be returned directly instead.
func ExtractDistributedTracingExtension(ctx context.Context, event cloudevents.Event) context.Context {
	tc := NewCloudEventTraceContext()
	carrier := NewCloudEventCarrierWithEvent(event)

	ctx = tc.Extract(ctx, carrier)

	return ctx
}
