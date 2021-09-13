package observability

import (
	"context"
	"fmt"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/cloudevents/sdk-go/v2/binding"
	"github.com/cloudevents/sdk-go/v2/observability"
	"github.com/cloudevents/sdk-go/v2/protocol"
	cehttp "github.com/cloudevents/sdk-go/v2/protocol/http"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	// TODO: What should we put here?
	instrumentationName         = "github.com/keptn/go-utils/observability/cloudevents"
	keptnContextCEExtension     = "shkeptncontext"
	keptnSpecVersionCEExtension = "shkeptnspecversion"
	triggeredIDCEExtension      = "triggeredid"
)

// OTelObservabilityService implements the ObservabilityService interface from cloudevents
type OTelObservabilityService struct {
	TracerProvider trace.TracerProvider
	Tracer         trace.Tracer
}

// InboundContextDecorators returns a decorator function that allows enriching the context with the incoming parent trace.
// This method gets invoked automatically by passing the option 'WithObservabilityService' when creating the cloudevents HTTP client.
func (os OTelObservabilityService) InboundContextDecorators() []func(context.Context, binding.Message) context.Context {
	return []func(context.Context, binding.Message) context.Context{tracePropagatorContextDecorator}
}

// RecordReceivedMalformedEvent records the error from a malformed event in the span.
func (os OTelObservabilityService) RecordReceivedMalformedEvent(ctx context.Context, err error) {
	spanName := fmt.Sprintf("%s receive", observability.ClientSpanName)
	_, span := os.Tracer.Start(
		ctx,
		spanName,
		trace.WithAttributes(attribute.String(string(semconv.CodeFunctionKey), "RecordReceivedMalformedEvent")))

	span.RecordError(err)
	span.End()
}

// RecordCallingInvoker starts a new span before calling the invoker upon a received event.
// In case the operation fails, the error is recorded and the span is marked as failed.
func (os OTelObservabilityService) RecordCallingInvoker(ctx context.Context, event *cloudevents.Event) (context.Context, func(errOrResult error)) {
	spanName := fmt.Sprintf("%s.%s receive", observability.ClientSpanName, event.Context.GetType())

	// span name will be: cloudevents.client.sh.keptn.event.hardening.evaluation.triggered send
	ctx, span := os.Tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindConsumer))

	if span.IsRecording() {
		span.SetAttributes(eventSpanAttributes(event, "RecordCallingInvoker")...)
	}

	return ctx, func(errOrResult error) {
		recordSpanError(span, errOrResult)
		span.End()
	}
}

// RecordSendingEvent starts a new span before sending the event.
// In case the operation fails, the error is recorded and the span is marked as failed.
func (os OTelObservabilityService) RecordSendingEvent(ctx context.Context, event cloudevents.Event) (context.Context, func(errOrResult error)) {
	spanName := fmt.Sprintf("%s.%s send", observability.ClientSpanName, event.Context.GetType())
	ctx, span := os.Tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindProducer))

	// TODO: Should we add more things here? What about sensitive information?
	if span.IsRecording() {
		span.SetAttributes(eventSpanAttributes(&event, "RecordSendingEvent")...)
	}

	return ctx, func(errOrResult error) {
		recordSpanError(span, errOrResult)
		span.End()
	}
}

// RecordRequestEvent starts a new span before transmitting the given.
// In case the operation fails, the error is recorded and the span is marked as failed.
func (os OTelObservabilityService) RecordRequestEvent(ctx context.Context, event cloudevents.Event) (context.Context, func(errOrResult error, event *cloudevents.Event)) {
	spanName := fmt.Sprintf("%s.%s process", observability.ClientSpanName, event.Context.GetType())
	ctx, span := os.Tracer.Start(ctx, spanName, trace.WithSpanKind(trace.SpanKindConsumer))

	if span.IsRecording() {
		span.SetAttributes(eventSpanAttributes(&event, "RecordRequestEvent")...)
	}

	return ctx, func(errOrResult error, event *cloudevents.Event) {
		recordSpanError(span, errOrResult)
		span.End()
	}
}

// NewOTelObservabilityService returns a OpenTelemetry enabled observability service
func NewOTelObservabilityService() *OTelObservabilityService {
	s := &OTelObservabilityService{
		TracerProvider: otel.GetTracerProvider(),
	}

	s.Tracer = s.TracerProvider.Tracer(
		instrumentationName,
		trace.WithInstrumentationVersion("1.0.0"), // TODO: Get the package version from somewhere?
	)

	return s
}

// Extracts the traceparent from the msg and enriches the context to enable propagation
func tracePropagatorContextDecorator(ctx context.Context, msg binding.Message) context.Context {
	var messageCtx context.Context
	if mctx, ok := msg.(binding.MessageContext); ok {
		messageCtx = mctx.Context()
	} else if mctx, ok := binding.UnwrapMessage(msg).(binding.MessageContext); ok {
		messageCtx = mctx.Context()
	}

	if messageCtx == nil {
		return ctx
	}
	span := trace.SpanFromContext(messageCtx)
	if span == nil {
		return ctx
	}
	return trace.ContextWithSpan(ctx, span)
}

func eventSpanAttributes(e *cloudevents.Event, method string) []attribute.KeyValue {
	attr := []attribute.KeyValue{
		attribute.String(string(semconv.CodeFunctionKey), method),
		attribute.String(observability.SpecversionAttr, e.SpecVersion()),
		attribute.String(observability.IdAttr, e.ID()),
		attribute.String(observability.TypeAttr, e.Type()),
		attribute.String(observability.SourceAttr, e.Source()),
	}
	if sub := e.Subject(); sub != "" {
		attr = append(attr, attribute.String(observability.SubjectAttr, sub))
	}
	if dct := e.DataContentType(); dct != "" {
		attr = append(attr, attribute.String(observability.DatacontenttypeAttr, dct))
	}
	if keptnContext, err := e.Context.GetExtension("shkeptncontext"); err == nil {
		attr = append(attr, attribute.String("shkeptncontext", keptnContext.(string)))
	}
	return attr
}

func recordSpanError(span trace.Span, errOrResult error) {
	if protocol.IsACK(errOrResult) || !span.IsRecording() {
		return
	}

	var httpResult *cehttp.Result
	if cloudevents.ResultAs(errOrResult, &httpResult) {
		span.RecordError(httpResult)
		if httpResult.StatusCode > 0 {
			span.SetStatus(semconv.SpanStatusFromHTTPStatusCode(httpResult.StatusCode))
		}
	} else {
		span.RecordError(errOrResult)
	}
}
