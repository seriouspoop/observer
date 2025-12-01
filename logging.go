package main

import (
	"context"

	"go.opentelemetry.io/contrib/bridges/otelzap"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"

	sdklogs "go.opentelemetry.io/otel/sdk/log"
)

type Logger struct {
	name     string
	provider *sdklogs.LoggerProvider
}

func NewLogger(ctx context.Context, name string, exp *Exporter) (*Logger, error) {
	ex, err := setupLogExporter(ctx, exp)
	if err != nil {
		return nil, err
	}
	provider, err := newLogProvider(ex, name)
	if err != nil {
		return nil, err
	}
	return &Logger{name: name, provider: provider}, nil
}

// Wrapper ---------------------------------------

func (l *Logger) NewLoggerCore() *otelzap.Core {
	core := otelzap.NewCore(l.name, otelzap.WithLoggerProvider(l.provider))
	return core
}

func (l *Logger) Shutdown(ctx context.Context) error {
	return l.provider.Shutdown(ctx)
}

// SDK -------------------------------------------

func newLogProvider(exp sdklogs.Exporter, name string) (*sdklogs.LoggerProvider, error) {
	r, err := newResource(name)
	if err != nil {
		return nil, err
	}

	processor := sdklogs.NewBatchProcessor(exp)

	// TODO - add sampling option
	return sdklogs.NewLoggerProvider(
		sdklogs.WithResource(r),
		sdklogs.WithProcessor(processor),
	), nil
}

// Exporters -----------------------------------------

func setupLogExporter(ctx context.Context, ex *Exporter) (e sdklogs.Exporter, err error) {
	switch ex.Type {
	case ConsoleExporter:
		e, err = newLogConsoleExporter()
	case OTLPHttpExporter:
		e, err = newLogOTLPHttpExporter(ctx, ex.HTTPEndpoint)
	case OTLPGrpcExporter:
		e, err = newLogOTLPGrpcExporter(ctx, ex.GRPCEndpoint)
	default:
		e, err = newLogConsoleExporter()
	}
	if err != nil {
		return nil, err
	}
	return
}

func newLogOTLPHttpExporter(ctx context.Context, otlpEndpoint string) (sdklogs.Exporter, error) {
	insecureOpts := otlploghttp.WithInsecure()
	endpoint := otlploghttp.WithEndpoint(otlpEndpoint)
	return otlploghttp.New(ctx, endpoint, insecureOpts)
}

func newLogOTLPGrpcExporter(ctx context.Context, otlpEndpoint string) (sdklogs.Exporter, error) {
	insecureOpts := otlploggrpc.WithInsecure()
	endpoint := otlploggrpc.WithEndpoint(otlpEndpoint)
	return otlploggrpc.New(ctx, endpoint, insecureOpts)
}

func newLogConsoleExporter() (sdklogs.Exporter, error) {
	return stdoutlog.New(stdoutlog.WithPrettyPrint())
}
