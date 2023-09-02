package main

import (
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
)

func otelResource() *resource.Resource {
	newResource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName("go-simple-app"),
		semconv.ServiceVersion("v0.1.0"),
		attribute.String("environment", "local"),
	)
	r := newResource

	return r
}

func newOTlPExporter() sdktrace.SpanExporter {
	// defaults to localhost:4317

	client := otlptracegrpc.NewClient(otlptracegrpc.WithInsecure(), otlptracegrpc.WithDialOption(grpc.WithBlock()))
	exporter, err := otlptrace.New(context.Background(), client)
	if err != nil {
		panic(err)
	}
	return exporter
}

func newTracerProvider() *sdktrace.TracerProvider {
	resource := otelResource()
	exporter := newOTlPExporter()
	return sdktrace.NewTracerProvider(sdktrace.WithBatcher(exporter), sdktrace.WithResource(resource))
}
