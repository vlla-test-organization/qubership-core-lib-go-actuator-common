package health

import (
	"errors"
	"net/http"
	"sync"
	"time"
)

type healthServiceImpl struct {
	checks             Checks
	statusMapping      StatusMapping
	statusOrder        StatusOrder
	healthAsyncResult  *HealthResult
	config             HealthConfiguration
	hasTimeDuration    bool
	statusResolver     func(Checks, StatusMapping, StatusOrder) *HealthResult
	mux                sync.RWMutex
	runChecksOnStartup bool
	started            bool
}

type HealthConfiguration struct {
	HealthCheckInterval time.Duration
}

func NewHealthService(healthCheckConfig ...HealthConfiguration) (HealthService, error) {
	if len(healthCheckConfig) > 1 {
		return nil, errors.New("more than one HealthConfiguration")
	}
	if len(healthCheckConfig) == 1 {
		return &healthServiceImpl{checks: make(map[string]Check), config: healthCheckConfig[0], hasTimeDuration: healthCheckConfig[0].HealthCheckInterval > 0, runChecksOnStartup: false}, nil
	}
	return &healthServiceImpl{checks: make(map[string]Check), hasTimeDuration: false, runChecksOnStartup: false}, nil
}

func (healthServiceImpl *healthServiceImpl) AddCheck(checkName string, check Check) HealthService {
	logger.Debugf("Add health indicator with name %s", checkName)
	healthServiceImpl.checks[checkName] = check
	return healthServiceImpl
}

func (healthServiceImpl *healthServiceImpl) RemoveCheck(checkName string) HealthService {
	logger.Debugf("Remove health indicator with name %s", checkName)
	delete(healthServiceImpl.checks, checkName)
	return healthServiceImpl
}

func (healthServiceImpl *healthServiceImpl) SetStatusMapping(statusMapping StatusMapping) HealthService {
	healthServiceImpl.statusMapping = statusMapping
	return healthServiceImpl
}

func (healthServiceImpl *healthServiceImpl) SetStatusOrder(statusOrder []HealthStatus) HealthService {
	healthServiceImpl.statusOrder = getStatusOrders(statusOrder)
	return healthServiceImpl
}

func (healthServiceImpl *healthServiceImpl) SetStatusResolver(resolver func(Checks, StatusMapping, StatusOrder) *HealthResult) HealthService {
	healthServiceImpl.statusResolver = resolver
	return healthServiceImpl
}

func (healthServiceImpl *healthServiceImpl) Start() HealthService {
	if healthServiceImpl.statusMapping == nil {
		healthServiceImpl.statusMapping = getDefaultStatusMapping()
		logger.Debug("No statusMapping provided, using default one: %s", healthServiceImpl.statusMapping)
	}
	if healthServiceImpl.statusOrder == nil {
		healthServiceImpl.statusOrder = getDefaultStatusOrder()
		logger.Debug("No StatusOrder provided, using default one: %s", healthServiceImpl.statusOrder)
	}

	if healthServiceImpl.statusResolver == nil {
		healthServiceImpl.statusResolver = getDefaultResolver()
		logger.Debug("No statusResolver provided, using default one")
	}

	if healthServiceImpl.runChecksOnStartup == true {
		logger.Debug("Start health calculation")
		healthServiceImpl.setResult(healthServiceImpl.statusResolver(healthServiceImpl.checks, healthServiceImpl.statusMapping, healthServiceImpl.statusOrder))
	} else {
		// set unavailable until first async result received from healthChecksRoutine
		healthServiceImpl.setResult(&HealthResult{
			StatusCode: http.StatusServiceUnavailable,
			HealthMap:  map[string]interface{}{"async-checks": "not-ready-yet"}})
	}
	if healthServiceImpl.hasTimeDuration {
		// start go routine for async health checks
		go healthServiceImpl.startHealthChecksInRoutine()
		logger.Info("Started health checks in go routine with interval: %s", (healthServiceImpl.config.HealthCheckInterval).String())
	} else {
		logger.Info("Health checks will be calculated by request")
	}
	return healthServiceImpl
}

func (healthServiceImpl *healthServiceImpl) GetHealth() *HealthResult {
	return healthServiceImpl.getResult()
}

func (healthServiceImpl *healthServiceImpl) RunChecksOnStartup(run bool) HealthService {
	healthServiceImpl.runChecksOnStartup = run
	return healthServiceImpl
}

func getDefaultStatusMapping() map[HealthStatus]int {
	statusToReturnCodes := make(map[HealthStatus]int)
	statusToReturnCodes[StatusUp] = http.StatusOK
	statusToReturnCodes[StatusWarning] = http.StatusOK
	statusToReturnCodes[StatusProblem] = http.StatusServiceUnavailable
	return statusToReturnCodes
}

func getDefaultStatusOrder() map[HealthStatus]int {
	return getStatusOrders([]HealthStatus{StatusProblem, StatusWarning, StatusUp})
}

func getStatusOrders(statusOrder []HealthStatus) map[HealthStatus]int {
	statusOrderMap := make(map[HealthStatus]int)
	for order, status := range statusOrder {
		statusOrderMap[status] = order - len(statusOrder)
	}
	return statusOrderMap
}
func getDefaultResolver() func(Checks, StatusMapping, StatusOrder) *HealthResult {
	return func(checks Checks, statusMapping StatusMapping, statusOrder StatusOrder) (result *HealthResult) {
		resultHealthMap := make(map[string]interface{})
		currentLowestStatus := StatusUp
		for checkName, check := range checks {
			checkStatus := check()
			statusMap := make(map[string]interface{})
			statusMap["status"] = checkStatus.GetStatus()
			for key, value := range checkStatus.GetDetails() {
				statusMap[key] = value
			}

			resultHealthMap[checkName] = statusMap
			if getStatusOrder(checkStatus.GetStatus(), statusOrder) < getStatusOrder(currentLowestStatus, statusOrder) {
				currentLowestStatus = checkStatus.GetStatus()
			}
		}
		// set summary status
		resultHealthMap["status"] = currentLowestStatus
		returnCode := getCodeForStatus(currentLowestStatus, statusMapping)
		result = &HealthResult{HealthMap: resultHealthMap, StatusCode: returnCode}
		return result
	}
}

func getStatusOrder(status HealthStatus, statusOrder StatusOrder) int {
	if statusOrder == nil {
		return -1
	}
	if order, ok := statusOrder[status]; ok {
		return order
	} else {
		return len(statusOrder) - 1
	}
}

func getCodeForStatus(status HealthStatus, statusMapping StatusMapping) int {
	if statusMapping == nil {
		return -1
	}
	if returnCode, ok := statusMapping[status]; ok {
		return returnCode
	} else {
		return http.StatusServiceUnavailable
	}
}

func (healthServiceImpl *healthServiceImpl) setResult(result *HealthResult) {
	healthServiceImpl.mux.Lock()
	defer healthServiceImpl.mux.Unlock()
	healthServiceImpl.healthAsyncResult = result
}

func (healthServiceImpl *healthServiceImpl) getResult() *HealthResult {
	if !healthServiceImpl.hasTimeDuration {
		healthServiceImpl.setResult(healthServiceImpl.statusResolver(healthServiceImpl.checks, healthServiceImpl.statusMapping, healthServiceImpl.statusOrder))
	}
	healthServiceImpl.mux.RLock()
	defer healthServiceImpl.mux.RUnlock()
	return healthServiceImpl.healthAsyncResult
}

func (healthServiceImpl *healthServiceImpl) startHealthChecksInRoutine() {
	tick := time.Tick(healthServiceImpl.config.HealthCheckInterval)
	for {
		select {
		case <-tick:
			healthServiceImpl.setResult(healthServiceImpl.statusResolver(healthServiceImpl.checks,
				healthServiceImpl.statusMapping,
				healthServiceImpl.statusOrder))
		}
	}
}
