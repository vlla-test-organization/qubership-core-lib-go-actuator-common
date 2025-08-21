package clpropertyutils

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vlla-test-organization/qubership-core-lib-go/v5/configloader"
)

func TestGetHttpBufferHeaderMaxSizeBytes(t *testing.T) {
	envVarKey := "HTTP_BUFFER_HEADER_MAX_SIZE"
	envVarValue := 9920
	err := os.Setenv(envVarKey, fmt.Sprint(envVarValue))
	assert.Nil(t, err)
	defer os.Unsetenv(envVarKey)

	configloader.InitWithSourcesArray([]*configloader.PropertySource{configloader.EnvPropertySource()})
	gotValue := GetHttpBufferHeaderMaxSizeBytes()
	assert.Equal(t, envVarValue, gotValue)
}
