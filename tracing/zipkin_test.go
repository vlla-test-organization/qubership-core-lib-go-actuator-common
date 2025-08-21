package tracing

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vlla-test-organization/qubership-core-lib-go/v5/configloader"
	"go.opentelemetry.io/otel"
)

func init() {

}

func setUpEnvs(envs ...string) {
	for _, env := range envs {
		splitedEnv := strings.Split(env, ":")
		if err := os.Setenv(splitedEnv[0], splitedEnv[1]); err != nil {
			panic(err)
		}
	}
	configloader.InitWithSourcesArray([]*configloader.PropertySource{configloader.EnvPropertySource()})
}

func removeEnvs(envs ...string) {
	for _, env := range envs {
		splitedEnv := strings.Split(env, ":")
		if err := os.Unsetenv(splitedEnv[0]); err != nil {
			panic(err)
		}
	}

}

var (
	TracingEnabled          = "TRACING_ENABLED:true"
	ServiceName             = "MICROSERVICE_NAME:someService"
	TracingHost             = "TRACING_HOST:localhost"
	TracingSamplerRateLimit = "TRACING_SAMPLER_RATELIMITING:10"
	Namespace               = "MICROSERVICE_NAMESPACE:someNamespace"
)

func TestZipkinOpenTelemetry(t *testing.T) {
	options := ZipkinOptions{ServiceName: "someService", TracingHost: "localhost", TracingSamplerRateLimiting: 10, TracingEnabled: true, Namespace: "test-namespace"}
	zipkinTracer := NewZipkinTracerWithOpts(options)
	registered, err := zipkinTracer.RegisterTracerProvider()
	assert.Nil(t, err)
	assert.True(t, registered)
	tracerProvider := otel.GetTracerProvider()
	assert.NotNil(t, tracerProvider)
}

func TestZipkinTracerIsOpenTelemetryExporter(t *testing.T) {
	options := ZipkinOptions{ServiceName: "someService", TracingHost: "localhost", TracingSamplerRateLimiting: 10, TracingEnabled: true, Namespace: "test-namespace"}
	zipkinTracer := NewZipkinTracerWithOpts(options)
	assert.Implements(t, (*OpenTelemetryExporter)(nil), zipkinTracer)
}

func TestZipkinOpenTelemetryViaEnv(t *testing.T) {
	setUpEnvs(TracingEnabled, ServiceName, TracingHost, TracingSamplerRateLimit, Namespace)
	zipkinTracer := NewZipkinTracer()
	registered, err := zipkinTracer.RegisterTracerProvider()
	assert.Nil(t, err)
	assert.True(t, registered)

	assert.Equal(t, zipkinTracer.zipkinOptions.TracingEnabled, true)
	assert.Equal(t, zipkinTracer.zipkinOptions.TracingHost, "localhost")
	assert.Equal(t, zipkinTracer.zipkinOptions.ServiceName, "someService")
	assert.Equal(t, zipkinTracer.zipkinOptions.TracingSamplerRateLimiting, 10)
	assert.Equal(t, zipkinTracer.zipkinOptions.Namespace, "someNamespace")
	removeEnvs(TracingEnabled, ServiceName, TracingHost, TracingSamplerRateLimit, Namespace)
}

func TestZipkinOpenTelemetryViaEnvDefaultValues(t *testing.T) {
	setUpEnvs(Namespace)
	zipkinTracer := NewZipkinTracer()
	registered, err := zipkinTracer.RegisterTracerProvider()
	assert.Nil(t, err)
	assert.False(t, registered)

	assert.Equal(t, zipkinTracer.zipkinOptions.TracingEnabled, false)
	assert.Equal(t, zipkinTracer.zipkinOptions.TracingHost, "")
	assert.Equal(t, zipkinTracer.zipkinOptions.ServiceName, "")
	assert.Equal(t, zipkinTracer.zipkinOptions.TracingSamplerRateLimiting, 10)

}

func TestZipkinOpenTelemetry_WithEmptyURL(t *testing.T) {
	options := ZipkinOptions{ServiceName: "someService", TracingHost: "", TracingEnabled: true, TracingSamplerRateLimiting: 10, Namespace: "test-namespace"}
	zipkinTracer := NewZipkinTracerWithOpts(options)
	registered, err := zipkinTracer.RegisterTracerProvider()
	assert.IsType(t, ZipkinUrlIsEmptyError{}, err)
	assert.False(t, registered)
}

func TestZipkinOpenTelemetry_WithEmptyServiceName(t *testing.T) {
	options := ZipkinOptions{ServiceName: "", TracingEnabled: true, Namespace: "test-namespace"}
	zipkinTracer := NewZipkinTracerWithOpts(options)
	if _, err := zipkinTracer.RegisterTracerProvider(); assert.Error(t, err) {
		assert.Contains(t, err.Error(), "you must specify serviceName in zipkin options or set microservice.name configuration parameter")
	}
}

func TestZipkinOpenTelemetry_WithIncorrectTracingSamplerRateLimiting(t *testing.T) {
	options := ZipkinOptions{TracingHost: "nc-diagnostic-agent", ServiceName: "someService", TracingEnabled: true, TracingSamplerRateLimiting: -10, Namespace: "test-namespace"}
	zipkinTracer := NewZipkinTracerWithOpts(options)
	if _, err := zipkinTracer.RegisterTracerProvider(); assert.Error(t, err) {
		assert.Contains(t, err.Error(), "tracing sampler rate limiting parameter must be more than 0")
	}
}
