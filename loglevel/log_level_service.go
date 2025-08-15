package loglevel

import (
	"github.com/vlla-test-organization/qubership-core-lib-go/v3/logging"
)

var logger logging.Logger

func init() {
	logger = logging.GetLogger("loglevel")
}

type LogLevelService interface {
	GetLogLevels() (*logging.LogLevels, error)
}
