package main

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

type Meter struct {
	name     string
	meter    metric.Meter
	provider *sdkmetric.MeterProvider
}

func NewMeter(ctx context.Context, name string, ex *Exporter) (*Meter, error) {
	exp, err := setupMetricExporter(ctx, ex)
	if err != nil {
		return nil, err
	}

	provider, err := newMeterProvider(exp, name)
	if err != nil {
		return nil, err
	}

	// set global provider for all libraries using otel - prometheus etc.
	otel.SetMeterProvider(provider)
	meter := provider.Meter(name)
	return &Meter{name: name, meter: meter, provider: provider}, nil
}

// Wrapper ---------------------------------------

// Name return name of the meter in use
func (m *Meter) Name() string {
	return m.name
}

func (m *Meter) Shutdown(ctx context.Context) error {
	return m.provider.Shutdown(ctx)
}

// SDK -------------------------------------------

func newMeterProvider(exp sdkmetric.Exporter, name string) (*sdkmetric.MeterProvider, error) {
	r, err := newResource(name)
	if err != nil {
		return nil, err
	}

	return sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(r),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exp))), nil
}

// Expoters --------------------------------------

func setupMetricExporter(ctx context.Context, ex *Exporter) (e sdkmetric.Exporter, err error) {
	switch ex.Type {
	case ConsoleExporter:
		e, err = newMetricConsoleExporter()
	case OTLPHttpExporter:
		e, err = newMetricOTLPHttpExporter(ctx, ex.HTTPEndpoint)
	case OTLPGrpcExporter:
		e, err = newMetricOTLPGrpcExporter(ctx, ex.GRPCEndpoint)
	default:
		e, err = newMetricConsoleExporter()
	}
	if err != nil {
		return nil, err
	}
	return
}

func newMetricOTLPHttpExporter(ctx context.Context, otlpEndpoint string) (sdkmetric.Exporter, error) {
	insecureOpts := otlpmetrichttp.WithInsecure()
	endpointOpts := otlpmetrichttp.WithEndpoint(otlpEndpoint)
	return otlpmetrichttp.New(ctx, insecureOpts, endpointOpts)
}

func newMetricOTLPGrpcExporter(ctx context.Context, otlpEndpoint string) (sdkmetric.Exporter, error) {
	insecureOpts := otlpmetricgrpc.WithInsecure()
	endpointOpts := otlpmetricgrpc.WithEndpoint(otlpEndpoint)
	return otlpmetricgrpc.New(ctx, insecureOpts, endpointOpts)
}

func newMetricConsoleExporter() (sdkmetric.Exporter, error) {
	return stdoutmetric.New(stdoutmetric.WithPrettyPrint())
}
