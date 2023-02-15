package tracer

import (
	"context"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"runtime"
	"strings"
)

type Closer func(context.Context) error

func InitTracerWithJaegerExporter(collectorEndpoint, serviceName, serviceVersion string) (Closer, error) {
	// Create and install Jaeger export pipeline.
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(collectorEndpoint)))
	if err != nil {
		return nil, err
	}

	return InitProviderWithCustomExporter(exporter, serviceName, serviceVersion), nil
}

func InitProviderWithCustomExporter(exporter sdktrace.SpanExporter, serviceName, serviceVersion string) Closer {
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			attribute.String("service.name", serviceName),
			attribute.String("service.version", serviceVersion),
			attribute.String("runtime.version", runtime.Version()),
		)),
	)
	otel.SetTracerProvider(tp)

	return tp.ForceFlush
}

func NewSpan(ctx context.Context, serviceName string, spanName ...string) (context.Context, trace.Span) {
	var name string
	if len(spanName) > 0 {
		name = strings.Join(spanName, " ")
	} else {
		name = getCallerFunctionName()
	}
	tr := otel.Tracer(serviceName)
	return tr.Start(ctx, name)
}

func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// getCallerFunctionName return name of function caller to
// the caller of this function.
func getCallerFunctionName() string {
	rpc := make([]uintptr, 1)
	runtime.Callers(3, rpc[:])
	fn := runtime.FuncForPC(rpc[0])
	funcName := fn.Name()

	for i := len(funcName) - 1; i >= 0; i-- {
		if funcName[i] == '.' {
			return funcName[i+1:]
		}
	}

	return funcName
}
