package monitoring

import (
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRegisterRequestStatusMetric(t *testing.T) {
	updateRegister()
	platformPrometheusMetrics, err := RegisterPlatformPrometheusMetrics(&Config{})
	assert.Nil(t, err)
	assert.NotNil(t, platformPrometheusMetrics.RequestLatencyHistogram)

	platformPrometheusMetrics.RequestStatusCounter.WithLabelValues("404", "POST", "").Inc()

	metric := findMetric("http_request_status")
	assert.NotNil(t, metric)
	assert.Equal(t, float64(1), metric.GetMetric()[0].GetCounter().GetValue())

}

func TestRegisterRequestLatencyHistogramMetric(t *testing.T) {
	updateRegister()
	platformPrometheusMetrics, err := RegisterPlatformPrometheusMetrics(&Config{})
	assert.Nil(t, err)
	assert.NotNil(t, platformPrometheusMetrics.RequestLatencyHistogram)

	platformPrometheusMetrics.RequestLatencyHistogram.WithLabelValues("404", "POST", "").Observe(0.5)

	metric := findMetric("http_request_time")
	//fmt.Print(metric)
	assert.NotNil(t, metric)
	histogram := metric.GetMetric()[0].GetHistogram()
	assert.Equal(t, 0.5, histogram.GetSampleSum())
	assert.Equal(t, uint64(1), histogram.GetSampleCount())
	assert.Equal(t, 0.2, histogram.GetBucket()[0].GetUpperBound())
	assert.Equal(t, 0.8, histogram.GetBucket()[1].GetUpperBound())
	assert.Equal(t, 1.0, histogram.GetBucket()[2].GetUpperBound())
}

func TestRegisterRequestLatencyHistogramMetricWithConfig(t *testing.T) {
	updateRegister()
	httpRequestTimeBuckets := []float64{0.005, 0.01, 0.05, 0.1}
	platformPrometheusMetrics, err := RegisterPlatformPrometheusMetrics(&Config{HttpRequestTimeBuckets: httpRequestTimeBuckets})
	assert.Nil(t, err)
	assert.NotNil(t, platformPrometheusMetrics.RequestLatencyHistogram)

	platformPrometheusMetrics.RequestLatencyHistogram.WithLabelValues("404", "POST", "").Observe(0.5)

	metric := findMetric("http_request_time")
	//fmt.Print(metric)
	assert.NotNil(t, metric)
	histogram := metric.GetMetric()[0].GetHistogram()
	assert.Equal(t, httpRequestTimeBuckets[0], histogram.GetBucket()[0].GetUpperBound())
	assert.Equal(t, httpRequestTimeBuckets[1], histogram.GetBucket()[1].GetUpperBound())
	assert.Equal(t, httpRequestTimeBuckets[2], histogram.GetBucket()[2].GetUpperBound())
	assert.Equal(t, httpRequestTimeBuckets[3], histogram.GetBucket()[3].GetUpperBound())
}

func updateRegister() {
	registry := prometheus.NewRegistry()
	prometheus.DefaultRegisterer = registry
	prometheus.DefaultGatherer = registry
}

func findMetric(metricName string) *dto.MetricFamily {
	if metrics, err := prometheus.DefaultGatherer.Gather(); err != nil {
		panic(err)
	} else {
		for _, metric := range metrics {
			if metric.GetName() == metricName {
				return metric
			}
		}
	}
	return nil
}
