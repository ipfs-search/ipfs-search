// Package instr contains components pertaining to observability and instrumentation.
package instr

import (
	"log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/api/global"
	"go.opentelemetry.io/otel/api/metric"
	"go.opentelemetry.io/otel/api/trace"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	"go.opentelemetry.io/otel/propagators"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

const (
	name = "github.com/ipfs-search"
)

// Instrumentation provides a canonical representation of instrumentation.
type Instrumentation struct {
	Tracer trace.Tracer
	Meter  metric.Meter
}

// Install configures and installs a Jaeger tracing pipeline. The first returned argument is a flusher, which should be called on program exit.
func Install(config *Config, serviceName string) (func(), error) {
	log.Printf("Creating Jaeger pipeline for service '%s' at ratio %f to endpoint %s", serviceName, config.SamplingRatio, config.JaegerEndpoint)

	// Configure context propagation
	global.SetTextMapPropagator(otel.NewCompositeTextMapPropagator(propagators.TraceContext{}, propagators.Baggage{}))

	// Configure sampler; default 1% of incoming requests (sniffed hashes)
	sampler := sdktrace.ParentBased(sdktrace.TraceIDRatioBased(config.SamplingRatio))

	return jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint(config.JaegerEndpoint),
		jaeger.WithProcess(jaeger.Process{ServiceName: serviceName}),
		jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sampler}),
	)
}

// New generates a representation of instrumentation containing the globally registered tracer and meter.
func New() *Instrumentation {
	return &Instrumentation{
		Tracer: global.Tracer(name),
		Meter:  global.Meter(name),
	}
}
