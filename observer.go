package main

import (
	"context"

	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

type Observer struct {
	name         string
	traceSDK     *Tracer
	logSDK       *Logger
	metricSDK    *Meter
	exporterType *Exporter
}

func New(ctx context.Context, name string, ex *Exporter) (*Observer, error) {
	tracer, err := NewTracer(ctx, name, ex)
	if err != nil {
		return nil, err
	}

	meter, err := NewMeter(ctx, name, ex)
	if err != nil {
		return nil, err
	}
	logger, err := NewLogger(ctx, name, ex)
	if err != nil {
		return nil, err
	}
	return &Observer{
		name:         name,
		traceSDK:     tracer,
		logSDK:       logger,
		metricSDK:    meter,
		exporterType: ex,
	}, nil
}

func (o *Observer) Shutdown(ctx context.Context) error {
	err := o.metricSDK.Shutdown(ctx)
	if err != nil {
		return err
	}
	err = o.logSDK.Shutdown(ctx)
	if err != nil {
		return err
	}
	return o.traceSDK.Shutdown(ctx)
}

// LogSDK returns logging SDK of the observer, use to extract log related instrumentation
func (o *Observer) LogSDK() *Logger {
	return o.logSDK
}

// TraceSDK returns tracing SDK of observer, use to extract trace related instrumentation
func (o *Observer) TraceSDK() *Tracer {
	return o.traceSDK
}

// MeterSDK returns metric SDK of observer, use to extract metric related instrumentation
func (o *Observer) MeterSDK() *Meter {
	return o.metricSDK
}

func newResource(name string) (*resource.Resource, error) {
	return resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(name),
		),
	)
}
