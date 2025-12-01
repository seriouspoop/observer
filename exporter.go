// Package main provides the exporter types and pre-built configs
package main

type ExporterType int

const (
	ConsoleExporter ExporterType = iota
	OTLPHttpExporter
	OTLPGrpcExporter
)

type Exporter struct {
	Type         ExporterType
	HTTPEndpoint string
	GRPCEndpoint string
}

func NewExporter(serviceName string, isProd bool) *Exporter {
	if isProd {
		return NewProductionExporter(serviceName)
	}
	return NewDevelopmentExporter()
}

// NewDevelopmentExporter provides console logging is handled by zap, this will not be used if it is development environment
func NewDevelopmentExporter() *Exporter {
	return &Exporter{
		Type:         ConsoleExporter,
		HTTPEndpoint: "localhost:4318",
		GRPCEndpoint: "localhost:4317",
	}
}

// NewProductionExporter returns a grpc otlp exporter.
// serviceName is the opentelemetry collector service name.
func NewProductionExporter(serviceName string) *Exporter {
	return &Exporter{
		Type:         OTLPGrpcExporter,
		HTTPEndpoint: serviceName + ":4318",
		GRPCEndpoint: serviceName + ":4317",
	}
}
