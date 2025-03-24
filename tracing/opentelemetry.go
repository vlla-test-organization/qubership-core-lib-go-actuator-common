package tracing

type OpenTelemetryExporter interface {
	RegisterTracerProvider() (bool, error)
	ServerName() string
}
