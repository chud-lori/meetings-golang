package infrastructure

import (
	"context"
	"errors"
	"fmt"
	logku "log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"

	//"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"

	//"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	//semconv "go.opentelemetry.io/otel/semconv/v1.12.0"
)

// setupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func SetupOTelSDK(ctx context.Context) (trace *trace.TracerProvider, shutdown func(context.Context) error, err error) {
	var shutdownFuncs []func(context.Context) error

	// shutdown calls cleanup functions registered via shutdownFuncs.
	// The errors from the calls are joined.
	// Each registered cleanup will be invoked once.
	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	// handleErr calls shutdown for cleanup and makes sure that all errors are returned.
	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up trace provider.
	tracerProvider, err := newTraceProvider(ctx)
	if err != nil {
        fmt.Println("EROROROROR: ", err)
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Set up meter provider.
	meterProvider, err := newMeterProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	// Set up logger provider.
	loggerProvider, err := newLoggerProvider()
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, loggerProvider.Shutdown)
	global.SetLoggerProvider(loggerProvider)

	return tracerProvider, shutdown, nil
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTraceProvider(ctx context.Context) (*trace.TracerProvider, error) {
    traceExporter, err := otlptracegrpc.New(
        ctx,
        otlptracegrpc.WithEndpoint("meetings-otel-collector:4317"),
        otlptracegrpc.WithInsecure(),
    )

    if err != nil {
        logku.Printf("Error creating trace exporter: %v", err)
        return nil, err
    }

    resource, err := resource.New(
        context.Background(),  // Context is needed here
        resource.WithAttributes(
            attribute.String("service.name", "meetings-app"),  // Correct service name
        ),
    )
    if err != nil {
        logku.Printf("Error creating trace exporter: %v", err)
    }

    traceProvider := trace.NewTracerProvider(
        trace.WithBatcher(traceExporter,
            trace.WithBatchTimeout(time.Second)),
        trace.WithResource(resource),
    )
	return traceProvider, nil
}

//func newMeterProvider() (*metric.MeterProvider, error) {
//    exporter, err := prometheus.New()
//
//	//metricExporter, err := stdoutmetric.New()
//	if err != nil {
//		return nil, err
//	}
//
//	meterProvider := metric.NewMeterProvider(
//		metric.WithReader(metric.NewPeriodicReader(exporter,
//			// Default is 1m. Set to 3s for demonstrative purposes.
//			metric.WithInterval(3*time.Second))),
//	)
//    // Expose the Prometheus metrics via HTTP server
//	go func() {
//		http.Handle("/metrics", exporter) // Expose /metrics for Prometheus
//		logku.Fatal(http.ListenAndServe(":8080", nil)) // Start server for metrics
//	}()
//	return meterProvider, nil
//}
func newMeterProvider() (*metric.MeterProvider, error) {
    // Create a new Prometheus exporter
    exporter, err := prometheus.New()
    if err != nil {
        return nil, err
    }
    // Create a new MeterProvider with the Prometheus exporter
    meterProvider := sdkmetric.NewMeterProvider(
        sdkmetric.WithReader(exporter),
    )
    // Expose the Prometheus metrics via HTTP server
    go func() {
        http.Handle("/metrics", promhttp.Handler())
        if err := http.ListenAndServe(":8080", nil); err != nil {
            panic(err)
        }
    }()

    return meterProvider, nil
}
func newLoggerProvider() (*log.LoggerProvider, error) {
	logExporter, err := stdoutlog.New()
	if err != nil {
		return nil, err
	}

	loggerProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
	)
	return loggerProvider, nil
}

