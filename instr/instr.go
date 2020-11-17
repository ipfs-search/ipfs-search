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

const name = "github.com/ipfs-search/instr"

type Instrumentation struct {
	Tracer trace.Tracer
	Meter  metric.Meter
}

func Install(serviceName string) (func(), error) {
	// First parameter is a flusher, should be called on context close!
	log.Printf("Creating Jaeger pipeline: %s", serviceName)

	// Configure context propagation
	global.SetTextMapPropagator(otel.NewCompositeTextMapPropagator(propagators.TraceContext{}, propagators.Baggage{}))

	// Configure sampler; default 1% of incoming requests (sniffed hashes)
	sampler := sdktrace.ParentBased(sdktrace.TraceIDRatioBased(0.01))
	// sampler := sdktrace.AlwaysSample()

	return jaeger.InstallNewPipeline(
		jaeger.WithAgentEndpoint("localhost:6831"),
		jaeger.WithProcess(jaeger.Process{ServiceName: serviceName}),
		jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sampler}),
		// jaeger.WithCollectorEndpoint("http://localhost:14268/api/traces"),
	)
}

func New() *Instrumentation {
	return &Instrumentation{
		Tracer: global.Tracer(name),
		Meter:  global.Meter(name),
	}
}
