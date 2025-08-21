package tracing

import (
	"fmt"
	"github.com/vlla-test-organization/qubership-core-lib-go-actuator-common/v5/tracing/utils"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
	"math"
)

// RateLimitingSampler samples at most maxTracesPerSecond. The distribution of sampled traces follows
// burstiness of the service, i.e. a service with uniformly distributed requests will have those
// requests sampled uniformly as well, but if requests are bursty, especially sub-second, then a
// number of sequential requests can be sampled each second.
type RateLimitingSampler struct {
	//legacySamplerV1Base
	maxTracesPerSecond float64
	rateLimiter        *utils.ReconfigurableRateLimiter
	description        string
}

// NewRateLimitingSampler creates new RateLimitingSampler.
func NewRateLimitingSampler(maxTracesPerSecond float64) sdktrace.Sampler {
	s := new(RateLimitingSampler)
	s.description = fmt.Sprintf("RateLimitingSampler{%f}", maxTracesPerSecond)
	return s.init(maxTracesPerSecond)
}

func (s *RateLimitingSampler) init(maxTracesPerSecond float64) *RateLimitingSampler {
	if s.rateLimiter == nil {
		s.rateLimiter = utils.NewRateLimiter(maxTracesPerSecond, math.Max(maxTracesPerSecond, 1.0))
	} else {
		s.rateLimiter.Update(maxTracesPerSecond, math.Max(maxTracesPerSecond, 1.0))
	}
	s.maxTracesPerSecond = maxTracesPerSecond
	return s
}

// String is used to log sampler details.
func (s *RateLimitingSampler) String() string {
	return fmt.Sprintf("RateLimitingSampler(maxTracesPerSecond=%v)", s.maxTracesPerSecond)
}

func (ts RateLimitingSampler) ShouldSample(p sdktrace.SamplingParameters) sdktrace.SamplingResult {
	psc := trace.SpanContextFromContext(p.ParentContext)
	if ts.rateLimiter.CheckCredit(1.0) {
		return sdktrace.SamplingResult{
			Decision:   sdktrace.RecordAndSample,
			Tracestate: psc.TraceState(),
		}
	}
	return sdktrace.SamplingResult{
		Decision:   sdktrace.Drop,
		Tracestate: psc.TraceState(),
	}
}

func (ts RateLimitingSampler) Description() string {
	return ts.description
}
