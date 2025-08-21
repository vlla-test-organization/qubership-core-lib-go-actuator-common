package monitoring

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/vlla-test-organization/qubership-core-lib-go/v5/logging"
)

var logger logging.Logger

type Config struct {
	HttpRequestTimeBuckets []float64
}

func init() {
	logger = logging.GetLogger("common-monitoring")
}

type PlatformPrometheusMetrics struct {
	RequestStatusCounter    *prometheus.CounterVec
	RequestLatencyHistogram *prometheus.HistogramVec
}

func (this *PlatformPrometheusMetrics) IncRequestStatusCounter(code, method, path string) {
	this.RequestStatusCounter.WithLabelValues(code, method, path).Inc()
}

func (this *PlatformPrometheusMetrics) ObserveRequestLatencyHistogram(code, method, path string, begin time.Time) {
	this.RequestLatencyHistogram.WithLabelValues(code, method, path).Observe(
		float64(time.Since(begin)) / float64(time.Second),
	)
}

func RegisterPlatformPrometheusMetrics(config *Config) (*PlatformPrometheusMetrics, error) {
	logger.Debugf("Start register core prometheus metrics")
	prometheusMiddleware := new(PlatformPrometheusMetrics)
	prometheusMiddleware.RequestStatusCounter = getRequestCounter()
	if err := prometheus.Register(prometheusMiddleware.RequestStatusCounter); err != nil {
		return nil, err
	}
	prometheusMiddleware.RequestLatencyHistogram = getRequestLatency(config.HttpRequestTimeBuckets)
	if err := prometheus.Register(prometheusMiddleware.RequestLatencyHistogram); err != nil {
		return nil, err
	}
	return prometheusMiddleware, nil
}

func getRequestLatency(buckets []float64) *prometheus.HistogramVec {
	if buckets == nil {
		buckets = []float64{0.2, 0.8, 1.0}
	}
	return prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request_time",
		Help:    "How long it took to process a request, partitioned by status code, method and HTTP path.",
		Buckets: buckets,
	},
		[]string{"status", "method", "path"},
	)
}

func getRequestCounter() *prometheus.CounterVec {
	return prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_request_status",
			Help: "How many HTTP requests processed, partitioned by status code, method and HTTP path.",
		},
		[]string{"status", "method", "path"},
	)
}
