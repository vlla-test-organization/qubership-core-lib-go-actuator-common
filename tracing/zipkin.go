package tracing

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/vlla-test-organization/qubership-core-lib-go/v3/configloader"
	"github.com/vlla-test-organization/qubership-core-lib-go/v3/logging"
	"go.opentelemetry.io/otel"
	zipkintr "go.opentelemetry.io/otel/exporters/zipkin"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

var logger logging.Logger

var ZIPKIN_ENDPOINT = "api/v2/spans"

func init() {
	logger = logging.GetLogger("tracer")
}

type ZipkinOptions struct {
	// required parameter
	TracingEnabled bool

	// required parameter
	TracingHost string

	// required parameter
	TracingSamplerRateLimiting int

	// required parameter
	ServiceName string

	// required parameter
	Namespace string
}

type zipkinTracer struct {
	zipkinOptions *ZipkinOptions
}

func NewZipkinTracer() *zipkinTracer {
	tracingEnabled, err := strconv.ParseBool(configloader.GetOrDefaultString("tracing.enabled", "false"))
	tracingHost := configloader.GetOrDefaultString("tracing.host", "")
	tracingSamplerRate, err := strconv.Atoi(configloader.GetOrDefaultString("tracing.sampler.ratelimiting", "10"))
	microserviceName := configloader.GetOrDefaultString("microservice.name", "")
	namespace := configloader.GetKoanf().MustString("microservice.namespace")
	logger.Debugf("tracingHost %s, microserviceName %s, namespace %s", tracingHost, microserviceName, namespace)
	if err != nil {
		panic(err)
	}

	zipkinOptions := ZipkinOptions{
		TracingEnabled:             tracingEnabled,
		TracingHost:                tracingHost,
		TracingSamplerRateLimiting: tracingSamplerRate,
		ServiceName:                microserviceName,
		Namespace:                  namespace,
	}
	return &zipkinTracer{zipkinOptions: &zipkinOptions}
}

func NewZipkinTracerWithOpts(zipkinOptions ZipkinOptions) *zipkinTracer {
	return &zipkinTracer{zipkinOptions: &zipkinOptions}
}

// Create Zipkin Exporter and install it as a global tracer.
func (this *zipkinTracer) RegisterTracerProvider() (bool, error) {
	if !this.zipkinOptions.TracingEnabled {
		logger.Debugf("zipkin tracer is disabled")
		return false, nil
	}
	logger.Debugf("register zipkin provider as opentracing provider")
	if err := this.checkConfigs(this.zipkinOptions); err != nil {
		return false, err
	}
	this.getSampler()
	exporter, err := zipkintr.New(
		fmt.Sprintf("http://%s:9411/%s", this.zipkinOptions.TracingHost, ZIPKIN_ENDPOINT),
	)
	if err != nil {
		return false, err
	}
	batcher := sdktrace.NewSimpleSpanProcessor(exporter)

	/* Note: service.namespace and service.name are not intended to be concatenated automatically
	for the purpose of forming a single globally unique name for the service.
	see https://github.com/open-telemetry/opentelemetry-specification/blob/main/specification/resource/semantic_conventions/README.md
	*/
	var tp = sdktrace.NewTracerProvider(
		sdktrace.WithSpanProcessor(batcher),
		sdktrace.WithResource(resource.NewWithAttributes(semconv.SchemaURL,
			semconv.ServiceNameKey.String(this.zipkinOptions.ServiceName+"-"+this.zipkinOptions.Namespace),
			semconv.ServiceNamespaceKey.String(this.zipkinOptions.Namespace))),
		sdktrace.WithSampler(this.getSampler()))

	otel.SetTracerProvider(tp)
	logger.Debug("zipkin tracer was registered as global tracer provider")
	return true, nil
}

func (this *zipkinTracer) getSampler() trace.Sampler {
	if this.zipkinOptions.TracingEnabled {
		return NewRateLimitingSampler(float64(this.zipkinOptions.TracingSamplerRateLimiting))
	}
	return sdktrace.NeverSample()
}

func (*zipkinTracer) checkConfigs(options *ZipkinOptions) error {
	if options.ServiceName == "" {
		return errors.New("you must specify serviceName in zipkin options or set microservice.name configuration parameter")
	}
	if options.TracingEnabled && options.TracingHost == "" {
		return ZipkinUrlIsEmptyError{}
	}
	if options.TracingSamplerRateLimiting <= 0 {
		return errors.New("tracing sampler rate limiting parameter must be more than 0")
	}
	return nil
}

func (this *zipkinTracer) ServerName() string {
	return this.zipkinOptions.ServiceName
}
