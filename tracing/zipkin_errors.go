package tracing

import "fmt"

type ZipkinUrlIsEmptyError struct {
}

func (z ZipkinUrlIsEmptyError) Error() string {
	return fmt.Sprintf("Tracer is enabled but tracingHost is empty. You must specify tracing.host parameter or disable tracing")
}
