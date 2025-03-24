package apiversion

import (
	"context"
	"os"
	"testing"

	"github.com/netcracker/qubership-core-lib-go/v3/configloader"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TestSuite struct {
	suite.Suite
}

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(TestSuite))
}

func (suite *TestSuite) TestGetApiVersion() {
	config := ApiVersionConfig{
		PathToApiVersionInfoFile: "./testdata/api-version-info.json",
	}
	apiVersionService, err := NewApiVersionService(config)
	assert.Nil(suite.T(), err)
	resp, err := apiVersionService.GetApiVersion(context.Background())
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), len(resp.Specs), 2)

	respFromCache, err := apiVersionService.GetApiVersion(context.Background())
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), respFromCache, resp)
}

func (suite *TestSuite) TestGetApiVersionPathError() {
	config := ApiVersionConfig{
		PathToApiVersionInfoFile: "../api-version-info.json",
	}
	apiVersionService, err := NewApiVersionService(config)
	assert.Nil(suite.T(), err)
	_, err = apiVersionService.GetApiVersion(context.Background())
	assert.NotNil(suite.T(), err)
}

func (suite *TestSuite) TestGetApiVersionJsonError() {
	config := ApiVersionConfig{
		PathToApiVersionInfoFile: "../apiversion/testdata/api-version-info-wrong.json",
	}
	apiVersionService, err := NewApiVersionService(config)
	assert.Nil(suite.T(), err)
	_, err = apiVersionService.GetApiVersion(context.Background())
	assert.NotNil(suite.T(), err)
}

func (suite *TestSuite) TestGetApiVersionWithProperty() {
	os.Setenv("API-VERSION_PATH", "./testdata/api-version-info.json")
	configloader.InitWithSourcesArray([]*configloader.PropertySource{configloader.EnvPropertySource()})
	config := ApiVersionConfig{}
	apiVersionService, err := NewApiVersionService(config)
	assert.Nil(suite.T(), err)
	resp, err := apiVersionService.GetApiVersion(context.Background())
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), len(resp.Specs), 2)

	respFromCache, err := apiVersionService.GetApiVersion(context.Background())
	assert.Nil(suite.T(), err)
	assert.Equal(suite.T(), respFromCache, resp)
}
