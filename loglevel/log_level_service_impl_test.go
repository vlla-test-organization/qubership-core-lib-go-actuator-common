package loglevel

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/vlla-test-organization/qubership-core-lib-go/v5/configloader"
	"github.com/vlla-test-organization/qubership-core-lib-go/v5/logging"
)

func TestEnvPropertySource(t *testing.T) {
	var loggerName = "logpackage"

	os.Setenv("LOGGING_LEVEL_LOGPACKAGE", logging.LvlCrit.String())
	defer func() {
		os.Clearenv()
	}()

	configloader.Init(configloader.EnvPropertySource())

	_ = logging.GetLogger(loggerName)

	logLevelService, err := NewLogLevelService()
	require.Nil(t, err)

	resp, err := logLevelService.GetLogLevels()
	require.Nil(t, err)
	logLevel := (*resp)[loggerName]
	require.Equal(t, strings.ToUpper(logging.LvlCrit.String()), logLevel)
}
