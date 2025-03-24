package loglevel

import (
	"github.com/netcracker/qubership-core-lib-go/v3/logging"
)

type loglevelServiceImpl struct {
}

func (loglevelServiceImpl *loglevelServiceImpl) GetLogLevels() (*logging.LogLevels, error) {
	logger.Debug("Get log levels for all loggers")
	levels := logging.GetLogLevels()
	return &levels, nil
}

func NewLogLevelService() (LogLevelService, error) {
	return &loglevelServiceImpl{}, nil
}
