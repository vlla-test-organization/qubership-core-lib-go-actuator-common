package apiversion

import (
	"context"
	"encoding/json"
	"errors"
	"os"

	"github.com/vlla-test-organization/qubership-core-lib-go/v3/configloader"
	"github.com/vlla-test-organization/qubership-core-lib-go/v3/logging"
)

var (
	logger logging.Logger
)

func init() {
	logger = logging.GetLogger("apiversion")
}

type apiVersionServiceImpl struct {
	apiVersionConfig        ApiVersionConfig
	apiVersionResponseCache *ApiVersionResponse
}

type ApiVersionConfig struct {
	PathToApiVersionInfoFile string
}

func (apiVersionServiceImpl *apiVersionServiceImpl) GetApiVersion(ctx context.Context) (*ApiVersionResponse, error) {
	if apiVersionServiceImpl.apiVersionResponseCache != nil {
		return apiVersionServiceImpl.apiVersionResponseCache, nil
	}
	resp, err := apiVersionServiceImpl.getSpecsFromFile()
	if err != nil {
		logger.Error("Error during parsing api version file: %v", err)
		return nil, err
	}
	return resp, nil
}

func NewApiVersionService(config ApiVersionConfig) (ApiVersionService, error) {
	if config.PathToApiVersionInfoFile == "" {
		config.PathToApiVersionInfoFile = configloader.GetOrDefaultString("api-version.path", "./api-version-info.json")
	}
	return &apiVersionServiceImpl{apiVersionConfig: config}, nil
}

func (apiVersionServiceImpl *apiVersionServiceImpl) getSpecsFromFile() (*ApiVersionResponse, error) {
	data, err := os.ReadFile(apiVersionServiceImpl.apiVersionConfig.PathToApiVersionInfoFile)
	if err != nil {
		return nil, err
	}

	var specs ApiVersionResponse
	err = json.Unmarshal(data, &specs)
	if err != nil {
		return nil, err
	}

	for _, spec := range specs.Specs {
		if spec.SpecRootUrl == "" {
			return nil, errors.New("spec.SpecRootUrl field can not be empty")
		}
		if spec.Minor == nil {
			return nil, errors.New("spec.Minor field can not be empty")
		}
		if spec.Major == nil {
			return nil, errors.New("spec.Major field can not be empty")
		}
		if len(spec.SupportedMajors) == 0 {
			return nil, errors.New("spec.SupportedMajors field can not be empty")
		}
	}
	apiVersionServiceImpl.apiVersionResponseCache = &specs
	return &specs, nil
}
