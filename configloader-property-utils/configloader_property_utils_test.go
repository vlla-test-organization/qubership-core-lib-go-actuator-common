package clpropertyutils

import (
	"fmt"
	"os"
	"testing"

	"github.com/netcracker/qubership-core-lib-go/v3/configloader"
	"github.com/stretchr/testify/assert"
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
