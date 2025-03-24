package apiversion

import "context"

type ApiVersionService interface {
	// Get API version info
	GetApiVersion(ctx context.Context) (*ApiVersionResponse, error)
}
