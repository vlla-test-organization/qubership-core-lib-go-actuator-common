package tracing

import (
	"fmt"
	"context"
	"github.com/stretchr/testify/assert"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"testing"
)

func TestNewRateLimitingSampler(t *testing.T) {
	var sampler sdktrace.Sampler
	sampler = NewRateLimitingSampler(float64(10))
	assert.Equal(t, sampler.Description(), (fmt.Sprintf("RateLimitingSampler{%f}", float64(10))))
}

func TestShouldSample_RecordAndSample_Drop(t *testing.T) {
	var sampler sdktrace.Sampler
	sampler = NewRateLimitingSampler(float64(1))
	traceID, _ := trace.TraceIDFromHex("4bf92f3577b34da6a3ce929d0e0e4736")
	spanID, _ := trace.SpanIDFromHex("00f067aa0ba902b7")
	parentCtx := trace.ContextWithSpanContext(
		context.Background(),
		trace.NewSpanContext(trace.SpanContextConfig{
			TraceID: traceID,
			SpanID:  spanID,
		}),
	)
	assert.Equal(t, sdktrace.RecordAndSample, sampler.ShouldSample(sdktrace.SamplingParameters{ParentContext: parentCtx}).Decision)
	assert.Equal(t, sdktrace.Drop, sampler.ShouldSample(sdktrace.SamplingParameters{ParentContext: parentCtx}).Decision)
}