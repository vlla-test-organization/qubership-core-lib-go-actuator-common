package health

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
	"time"
)

func TestHealthUpStatus(t *testing.T) {
	assert_ := assert.New(t)
	handler, err := NewHealthService(HealthConfiguration{time.Hour})
	assert.Nil(t, err)
	testCheck := func() (status Status) { return Status{Name: StatusUp} }
	handler.
		AddCheck("test-check", testCheck).
		RunChecksOnStartup(true).
		Start()

	healthResult := handler.GetHealth()
	assert_.Equal(http.StatusOK, healthResult.StatusCode)

	status_value, status_ok := healthResult.HealthMap["status"]
	assert_.Equal(StatusUp, status_value)
	assert_.Equal(true, status_ok)

	_, test_check_ok := healthResult.HealthMap["test-check"]
	assert_.Equal(true, test_check_ok)
}

func TestHealthWithoutDuration(t *testing.T) {
	assert_ := assert.New(t)

	handler, err := NewHealthService()
	assert.Nil(t, err)
	testCheck := func() (status Status) { return Status{Name: StatusUp} }
	handler.
		AddCheck("test-check", testCheck).
		Start()

	healthResult := handler.GetHealth()
	assert_.Equal(http.StatusOK, healthResult.StatusCode)

	status_value, status_ok := healthResult.HealthMap["status"]
	assert_.Equal(StatusUp, status_value)
	assert_.Equal(true, status_ok)

	_, test_check_ok := healthResult.HealthMap["test-check"]
	assert_.Equal(true, test_check_ok)
}

func TestHealthWithWarningAndUp(t *testing.T) {
	assert_ := assert.New(t)

	handler, err := NewHealthService(HealthConfiguration{time.Hour})
	assert.Nil(t, err)
	testCheckUp := func() (status Status) { return Status{Name: StatusUp} }
	testCheckWarning := func() (status Status) { return Status{Name: StatusWarning} }

	handler.
		AddCheck("test-check-up", testCheckUp).
		AddCheck("test-check-warning", testCheckWarning).
		RunChecksOnStartup(true).
		Start()

	healthResult := handler.GetHealth()
	assert_.Equal(http.StatusOK, healthResult.StatusCode)

	status_value, status_ok := healthResult.HealthMap["status"]
	assert_.Equal(StatusWarning, status_value)
	assert_.Equal(true, status_ok)

	_, test_check_ok := healthResult.HealthMap["test-check-up"]
	assert_.Equal(true, test_check_ok)
}

func TestHealthWithWarningAndUpAndProblem(t *testing.T) {
	assert_ := assert.New(t)

	handler, err := NewHealthService(HealthConfiguration{time.Hour})
	assert.Nil(t, err)
	testCheckUp := func() (status Status) { return Status{Name: StatusUp} }
	testCheckWarning := func() (status Status) { return Status{Name: StatusWarning} }
	testCheckProblem := func() (status Status) { return Status{Name: StatusProblem} }

	handler.
		AddCheck("test-check-up", testCheckUp).
		AddCheck("test-check-warning", testCheckWarning).
		AddCheck("test-check-problem", testCheckProblem).
		RunChecksOnStartup(true).
		Start()

	healthResult := handler.GetHealth()
	assert_.Equal(http.StatusServiceUnavailable, healthResult.StatusCode)

	status_value, status_ok := healthResult.HealthMap["status"]
	assert_.Equal(StatusProblem, status_value)
	assert_.Equal(true, status_ok)

	_, test_check_ok := healthResult.HealthMap["test-check-problem"]
	assert_.Equal(true, test_check_ok)
}

func TestErrorHealthCheckCreation(t *testing.T) {
	assert_ := assert.New(t)

	_, err := NewHealthService(HealthConfiguration{time.Hour}, HealthConfiguration{time.Millisecond})
	assert_.NotNil(t, err)
}

func TestHealthWithProblemAndError(t *testing.T) {
	assert_ := assert.New(t)

	handler, err := NewHealthService(HealthConfiguration{time.Hour})
	assert.Nil(t, err)

	test_error := errors.New("test error")
	availabilityFunc := func() Status {
		details := make(map[string]interface{})
		details["error"] = test_error.Error()
		return Status{Name: StatusProblem, Details: details}
	}
	handler.AddCheck("test-check-problem", availabilityFunc).
		RunChecksOnStartup(true).
		Start()

	healthResult := handler.GetHealth()
	assert_.Equal(http.StatusServiceUnavailable, healthResult.StatusCode)

	status_value, status_ok := healthResult.HealthMap["status"]
	assert_.Equal(StatusProblem, status_value)
	assert_.Equal(true, status_ok)

	problem, test_check_ok := healthResult.HealthMap["test-check-problem"]
	assert_.Equal(true, test_check_ok)
	problemAsMap, ok := problem.(map[string]interface{})
	assert_.Equal(true, ok)
	assert_.Equal(test_error.Error(), problemAsMap["error"])
}
