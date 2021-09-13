package observability

import (
	"context"
	"log"
	"net/url"

	"github.com/keptn/go-utils/pkg/common/osutils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

// EnvVarOTelCollectorEndpoint the env var containing the OpenTelemetry Collector endpoint
const EnvVarOTelCollectorEndpoint = "OTEL_COLLECTOR_ENDPOINT"

// InitOTelTraceProvider configures the OpenTelemetry SDK to export the spans to collector via OTLP/GRPC
//
// The SDK uses the collector endpoint defined in the environment variable: OTEL_COLLECTOR_ENDPOINT
// The environment variable can be set by adding it to the values.yaml file and doing a `helm upgrade`.
func InitOTelTraceProvider(serviceName string) func() {
	ctx := context.Background()

	collectorGrpcEndpoint := osutils.GetOSEnv(EnvVarOTelCollectorEndpoint)

	_, err := url.ParseRequestURI(collectorGrpcEndpoint)
	if err != nil {
		otel.SetTracerProvider(trace.NewNoopTracerProvider())
		return func() {}
	}

	// TODO: Depending how the collector is deployed, we might need
	// more things to be able to talk to it (authorization for ex).
	// So most likely we will need to more settings for it. For now, we are assuming
	// that there's a local collector running inside the cluster.
	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithEndpoint(collectorGrpcEndpoint),
	)
	if err != nil {
		return func() { log.Printf("Failed to create the trace exporter: %v", err) }
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	return func() {
		if tp == nil {
			return
		}
		if err := tp.Shutdown(ctx); err != nil {
			log.Printf("Error shutting down the tracer provider: %v", err)
		}
	}
}
