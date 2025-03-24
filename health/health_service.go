package health

import (
	"github.com/netcracker/qubership-core-lib-go/v3/logging"
)

var logger logging.Logger

const (
	StatusUp      HealthStatus = "UP"
	StatusWarning HealthStatus = "WARNING"
	StatusProblem HealthStatus = "PROBLEM"
)

func init() {
	logger = logging.GetLogger("health")
}

type HealthService interface {
	// AddCheck adds health indicator which is responsible for doing some health checking.
	// checkName is a name of indicator, e.g.: PostgresqlHealthCheck or KubernatesHealthCheck.
	// Check is a function which does health checking.
	AddCheck(checkName string, check Check) HealthService
	// RemoveCheck removes a health indicator
	RemoveCheck(checkName string) HealthService
	// Map which determines outgoing request http code by resulting health status, e.g.:
	// "UP" : 200
	// "WARNING" : 200
	// "PROBLEM" : 503
	SetStatusMapping(statusMapping StatusMapping) HealthService
	// Set priority of HealthStatus by which the resulting HealthStatus will be selected
	SetStatusOrder(statusOrder []HealthStatus) HealthService
	// Allows to set function which calculates HealthResult
	SetStatusResolver(resolver func(Checks, StatusMapping, StatusOrder) *HealthResult) HealthService
	// Starting the periodic health calculation process
	Start() HealthService
	// Get health result which contains all health indicators, overall status and http code
	GetHealth() *HealthResult
	// Should health be calculated immediately after calling the method Start
	RunChecksOnStartup(run bool) HealthService
}

type (
	HealthStatus string
	// Check functional checks a state of some program unit
	Check func() (status Status)
	// Map where key is a name of checking system. For example: MongoHealthCheck or KubernetesHealthCheck
	Checks map[string]Check
	// Map which contains the priority of HealthStatus and by which the resulting status is determined
	StatusOrder map[HealthStatus]int
	// Map which determines outgoing request http code by result health status, e.g.:
	// "UP" : 200
	// "WARNING" : 200
	// "PROBLEM" : 503
	StatusMapping map[HealthStatus]int
	// Contains information about health probe of check
	Status struct {
		Name    HealthStatus
		Details map[string]interface{}
	}
	// Health response
	HealthResult struct {
		HealthMap  map[string]interface{}
		StatusCode int
	}
)

func (status Status) GetStatus() HealthStatus {
	return status.Name
}

func (status Status) GetDetails() map[string]interface{} {
	return status.Details
}

func (healthResult HealthResult) GetHealthMap() map[string]interface{} {
	return healthResult.HealthMap
}

func (healthResult HealthResult) GetStatusCode() int {
	return healthResult.StatusCode
}
