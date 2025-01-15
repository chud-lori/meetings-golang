package infrastructure

import (
	"context"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

type statusCodeRecorder struct {
    http.ResponseWriter
    statusCode int
}

func (r *statusCodeRecorder) WriteHeader(statusCode int) {
    r.statusCode = statusCode
    r.ResponseWriter.WriteHeader(statusCode)
}

func TraceMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

        tracer := otel.Tracer("meetings-app")

        // Create a span
        ctx, span := tracer.Start(r.Context(), "http-request")
        defer span.End()

        recorder := &statusCodeRecorder{
            ResponseWriter: w,
            statusCode:     http.StatusOK, // Default status code is OK
        }

        r = r.WithContext(ctx)
        start := time.Now()
		// call the next handler
        next.ServeHTTP(w, r)

        execTime := time.Since(start)
        span.SetAttributes(
            attribute.String("http.method", r.Method),
            attribute.String("http.url", r.URL.String()),
            attribute.Float64("execution.time_ms", float64(execTime.Milliseconds())),
            attribute.Int("http.status_code", recorder.statusCode), // Add the status code to the span

        )
	})
}

func TraceFunction(ctx context.Context, spanName string, attributes ...attribute.KeyValue) (context.Context, func()) {
    tracer := otel.Tracer("meetings-app")

    // Start a child span
    ctx, span := tracer.Start(ctx, spanName)
    span.SetAttributes(attributes...)

    start := time.Now()

    return ctx, func() {
        execTime := time.Since(start)
        span.SetAttributes(attribute.Float64("execution.time_ms", float64(execTime.Milliseconds())))
        span.End()
    }
}

