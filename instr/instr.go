// Package instr contains common instrumentation tooling.
package instr

import (
	"context"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	name = "github.com/ipfs-search"
)

// Instrumentation provides a canonical representation of instrumentation.
type Instrumentation struct {
	Tracer trace.Tracer
}

// Install configures and installs a Jaeger tracing pipeline. The first returned argument is a flusher, which should be called on program exit.
func Install(config *Config, serviceName string) (func(context.Context), error) {
	log.Printf("Creating Jaeger pipeline for service '%s' at ratio %f to endpoint %s", serviceName, config.SamplingRatio, config.JaegerEndpoint)

	// Configure context propagation
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	// Configure sampler; default 1% of incoming requests (sniffed hashes)
	sampler := tracesdk.ParentBased(tracesdk.TraceIDRatioBased(config.SamplingRatio))

	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.JaegerEndpoint)))
	if err != nil {
		return nil, err
	}

	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		tracesdk.WithSampler(sampler),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
		)),
	)

	// Register our TracerProvider as the global so any imported
	// instrumentation in the future will default to using it.
	otel.SetTracerProvider(tp)

	return func(ctx context.Context) {
		// Do not make the application hang when it is shutdown.
		ctx, cancel := context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}, nil
	// return jaeger.InstallNewPipeline(
	// 	jaeger.WithCollectorEndpoint(config.JaegerEndpoint),
	// 	jaeger.WithProcess(jaeger.Process{ServiceName: serviceName}),
	// 	jaeger.WithSDK(&sdktrace.Config{DefaultSampler: sampler}),
	// )
}

// New generates a representation of instrumentation containing the globally registered tracer and meter.
func New() *Instrumentation {
	return &Instrumentation{
		Tracer: otel.Tracer(name),
	}
}
