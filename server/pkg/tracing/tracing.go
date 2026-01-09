package tracing

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
	"go.opentelemetry.io/otel/trace"
	nooptrace "go.opentelemetry.io/otel/trace/noop"
)

type contextKey string

const TraceIDKey contextKey = "trace_id"

var tracer trace.Tracer

// Config holds tracing configuration
type Config struct {
	ServiceName    string
	ServiceVersion string
	Environment    string
	JaegerEndpoint string
	SamplingRate   float64 // 0.0 to 1.0 (1.0 = 100%)
}

// getOTLPEndpoint reads the OTLP endpoint from config or environment variables
// Priority: cfg.JaegerEndpoint > JAEGER_ENDPOINT env > OTLP_ENDPOINT env > default
func getOTLPEndpoint(cfg Config) string {
	if cfg.JaegerEndpoint != "" {
		return cfg.JaegerEndpoint
	}

	if endpoint := os.Getenv("JAEGER_ENDPOINT"); endpoint != "" {
		return endpoint
	}

	if endpoint := os.Getenv("OTLP_ENDPOINT"); endpoint != "" {
		return endpoint
	}

	return "http://jaeger:4318"
}

// Init initializes OpenTelemetry tracing with OTLP HTTP exporter (for Jaeger)
func Init(cfg Config) (func(), error) {
	otlpEndpoint := getOTLPEndpoint(cfg)
	otlpEndpoint = normalizeEndpoint(otlpEndpoint)

	if cfg.SamplingRate == 0 {
		env := strings.ToLower(cfg.Environment)
		if env == "production" || env == "prod" {
			cfg.SamplingRate = 0.1
		} else {
			cfg.SamplingRate = 1.0
		}
	}

	endpointHost := extractHostPort(otlpEndpoint)

	ctx := context.Background()
	exp, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(endpointHost),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	res, err := resource.New(
		context.Background(),
		resource.WithAttributes(
			semconv.ServiceNameKey.String(cfg.ServiceName),
			semconv.ServiceVersionKey.String(cfg.ServiceVersion),
			attribute.String("environment", cfg.Environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(res),
		tracesdk.WithSampler(tracesdk.TraceIDRatioBased(cfg.SamplingRate)),
	)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	tracer = otel.Tracer(cfg.ServiceName)

	return func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			fmt.Printf("Error shutting down tracer provider: %v\n", err)
		}
	}, nil
}

// Tracer returns the global tracer
func Tracer() trace.Tracer {
	if tracer == nil {
		return nooptrace.NewTracerProvider().Tracer("noop")
	}
	return tracer
}

// StartSpan starts a new span with the given name and options
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	return Tracer().Start(ctx, name, opts...)
}

// SpanFromContext extracts span from context
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// TraceIDFromContext extracts trace ID from context and returns it as a string
func TraceIDFromContext(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		return span.SpanContext().TraceID().String()
	}
	if traceID, ok := ctx.Value(TraceIDKey).(string); ok {
		return traceID
	}
	return ""
}

// ContextWithTraceID creates a context with trace_id and OpenTelemetry span
func ContextWithTraceID(ctx context.Context, traceID string) context.Context {
	if traceID == "" {
		return ctx
	}
	return context.WithValue(ctx, TraceIDKey, traceID)
}

// SetSpanAttributes sets attributes on the span from context
func SetSpanAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span != nil {
		span.SetAttributes(attrs...)
	}
}

// RecordError records an error on the span
func RecordError(ctx context.Context, err error) {
	span := trace.SpanFromContext(ctx)
	if span != nil && err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
	}
}

// normalizeEndpoint normalizes the OTLP endpoint URL to prevent double-prefixing
// and ensures it's in the correct format for OpenTelemetry
func normalizeEndpoint(endpoint string) string {
	if endpoint == "" {
		return "http://jaeger:4318"
	}

	endpoint = strings.TrimSpace(endpoint)

	if strings.Contains(endpoint, "%") {
		if decoded, err := url.QueryUnescape(endpoint); err == nil && decoded != endpoint {
			endpoint = decoded
		}
	}

	maxIterations := 10
	for i := 0; i < maxIterations; i++ {
		original := endpoint
		endpoint = strings.Replace(endpoint, "http://http://", "http://", 1)
		endpoint = strings.Replace(endpoint, "https://https://", "https://", 1)
		endpoint = strings.Replace(endpoint, "http://https://", "https://", 1)
		endpoint = strings.Replace(endpoint, "https://http://", "http://", 1)
		if endpoint == original {
			break
		}
	}

	parsedURL, err := url.Parse(endpoint)
	if err != nil {
		cleanEndpoint := endpoint
		for strings.HasPrefix(cleanEndpoint, "http://") || strings.HasPrefix(cleanEndpoint, "https://") {
			cleanEndpoint = strings.TrimPrefix(cleanEndpoint, "http://")
			cleanEndpoint = strings.TrimPrefix(cleanEndpoint, "https://")
		}

		if idx := strings.IndexAny(cleanEndpoint, "/?#"); idx > 0 {
			cleanEndpoint = cleanEndpoint[:idx]
		}

		if cleanEndpoint != "" {
			if strings.Contains(cleanEndpoint, ":") {
				return "http://" + cleanEndpoint
			}
			return "http://" + cleanEndpoint + ":4318"
		}

		return "http://jaeger:4318"
	}

	if parsedURL.Host == "" {
		return "http://jaeger:4318"
	}

	normalized := fmt.Sprintf("%s://%s", parsedURL.Scheme, parsedURL.Host)
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		normalized = fmt.Sprintf("http://%s", parsedURL.Host)
	}

	return normalized
}

// extractHostPort extracts just the host:port from a URL string
func extractHostPort(endpoint string) string {
	if endpoint == "" {
		return "jaeger:4318"
	}

	endpoint = strings.TrimSpace(endpoint)
	endpoint = strings.TrimPrefix(endpoint, "http://")
	endpoint = strings.TrimPrefix(endpoint, "https://")

	if decoded, err := url.QueryUnescape(endpoint); err == nil && decoded != endpoint {
		endpoint = decoded
		endpoint = strings.TrimPrefix(endpoint, "http://")
		endpoint = strings.TrimPrefix(endpoint, "https://")
	}

	if idx := strings.Index(endpoint, "/"); idx > 0 {
		endpoint = endpoint[:idx]
	}
	if idx := strings.Index(endpoint, "?"); idx > 0 {
		endpoint = endpoint[:idx]
	}

	if strings.Contains(endpoint, ":") {
		return endpoint
	}
	if endpoint != "" {
		return endpoint + ":4318"
	}

	return "jaeger:4318"
}

