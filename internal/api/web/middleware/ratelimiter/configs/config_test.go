package configs_test

import (
	"ratelim/internal/api/web/middleware/ratelimiter/configs"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadConfig_MustNot_WhenDefaultIsMissing(t *testing.T) {
	// Arrange
	cfg, err := configs.LoadConfig("services_first_test.yaml")

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "error validating config: default service is missing")
	assert.Nil(t, cfg)
}

func TestLoadConfig_MustNotPass_WhenServiceNameIsMissing(t *testing.T) {
	// Arrange
	cfg, err := configs.LoadConfig("services_third_test.yaml")

	// Assert
	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "error validating config: services cannot be empty")
	assert.Nil(t, cfg)
}

func TestLoadConfig_MustNotAllowDuplicateServiceNames(t *testing.T) {
	// Arrange
	cfg, err := configs.LoadConfig("services_fourth_test.yaml")

	// Assert
	assert.Nil(t, err)
	assert.Equal(t, len(cfg.Services), 2)
}

func TestLoadConfig_PassForProperConfig(t *testing.T) {
	// Arrange
	cfg, err := configs.LoadConfig("services_second_test.yaml")

	// Assert
	assert.Nil(t, err)
	assert.NotNil(t, cfg)
	assert.Equal(t, len(cfg.Services), 5)
	//Default Service
	assert.Equal(t, cfg.Services[0].Name, "default")
	assert.Equal(t, cfg.Services[0].Type, "ip")
	assert.Equal(t, cfg.Services[0].Address, "any")
	assert.Equal(t, cfg.Services[0].Key, "default")
	assert.Equal(t, cfg.Services[0].Valid, true)
	assert.Equal(t, cfg.Services[0].AllowedRPS, 10)
	assert.Equal(t, cfg.Services[0].WaitTimeIfLimitExceeded, "5m")
	//Service A
	assert.Equal(t, cfg.Services[1].Name, "service-a")
	assert.Equal(t, cfg.Services[1].Type, "token")
	assert.Equal(t, cfg.Services[1].Key, "abcd1234")
	assert.Equal(t, cfg.Services[1].Address, "")
	assert.Equal(t, cfg.Services[1].Valid, true)
	assert.Equal(t, cfg.Services[1].AllowedRPS, 20)
	assert.Equal(t, cfg.Services[1].WaitTimeIfLimitExceeded, "10s")
	//Service B
	assert.Equal(t, cfg.Services[2].Name, "service-b")
	assert.Equal(t, cfg.Services[2].Type, "token")
	assert.Equal(t, cfg.Services[2].Key, "efgh5678")
	assert.Equal(t, cfg.Services[2].Valid, true)
	assert.Equal(t, cfg.Services[2].AllowedRPS, 30)
	assert.Equal(t, cfg.Services[2].WaitTimeIfLimitExceeded, "5s")
	//Service C
	assert.Equal(t, cfg.Services[3].Name, "service-c")
	assert.Equal(t, cfg.Services[3].Type, "token")
	assert.Equal(t, cfg.Services[3].Key, "ijkl91011")
	assert.Equal(t, cfg.Services[3].Address, "")
	assert.Equal(t, cfg.Services[3].Valid, false)
	assert.Equal(t, cfg.Services[3].AllowedRPS, 0)
	assert.Equal(t, cfg.Services[3].WaitTimeIfLimitExceeded, "")
	// //Service D
	assert.Equal(t, cfg.Services[4].Name, "service-d")
	assert.Equal(t, cfg.Services[4].Type, "token")
	assert.Equal(t, cfg.Services[4].Key, "mnop121314")
	assert.Equal(t, cfg.Services[4].Address, "")
	assert.Equal(t, cfg.Services[4].Valid, true)
	assert.Equal(t, cfg.Services[4].AllowedRPS, 10)
	assert.Equal(t, cfg.Services[4].WaitTimeIfLimitExceeded, "5m")
}
