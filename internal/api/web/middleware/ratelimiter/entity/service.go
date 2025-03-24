package entity

import (
	"errors"
	"fmt"
)

type ServiceConfig struct {
	Name                    string `mapstructure:"name"`
	Type                    string `mapstructure:"type"`
	Address                 string `mapstructure:"address"`
	Key                     string `mapstructure:"key"`
	Valid                   bool   `mapstructure:"valid"`
	AllowedRPS              int    `mapstructure:"allowed_rps"`
	WaitTimeIfLimitExceeded string `mapstructure:"wait_time_if_limit_exceeded"`
}

type Config struct {
	Services []*ServiceConfig `mapstructure:"services"`
}

func (c *Config) Validate() []error {
	var Errors []error
	var ValidServices []*ServiceConfig
	seenNames := make(map[string]bool)

	// Check if services is empty
	if len(c.Services) == 0 {
		return append(Errors, errors.New("services cannot be empty"))
	}

	// Check if default service is missing
	var defaultWaitTime string
	var defaultAllowedRPS int
	var hasDefault bool

	for _, s := range c.Services {
		if s.Name == "default" {
			defaultWaitTime = s.WaitTimeIfLimitExceeded
			defaultAllowedRPS = s.AllowedRPS
			hasDefault = true
			s.Key = "default"
		}
	}
	if !hasDefault {
		c.Services = []*ServiceConfig{}
		return append(Errors, errors.New("default service is missing"))
	}

	// Check which services are valid
	for _, s := range c.Services {
		if s.Name == "" {
			Errors = append(Errors, errors.New("service name cannot be empty"))
			continue
		}

		// Skip duplicates: only the first valid service with a name is accepted
		if seenNames[s.Name] {
			continue
		}

		if s.Type != "ip" && s.Type != "token" {
			Errors = append(Errors, fmt.Errorf("invalid type for service '%s': must be 'ip' or 'token'", s.Name))
			continue
		}

		if s.AllowedRPS < 0 {
			Errors = append(Errors, fmt.Errorf("allowed_rps must be >= 0 for service '%s'", s.Name))
			continue
		}

		if s.Type == "token" && s.Key == "" {
			Errors = append(Errors, fmt.Errorf("key cannot be empty for service '%s' of type 'token'", s.Name))
			continue
		}

		if s.Type == "ip" && s.Address == "" {
			Errors = append(Errors, fmt.Errorf("address cannot be empty for service '%s' of type 'ip'", s.Name))
			continue
		}

		// Mark name as used and add to valid list
		seenNames[s.Name] = true
		ValidServices = append(ValidServices, s)
	}

	// Set defaults for missing config
	for _, vs := range ValidServices {
		if vs.WaitTimeIfLimitExceeded == "" && vs.Valid {
			vs.WaitTimeIfLimitExceeded = defaultWaitTime
		}
		if vs.AllowedRPS == 0 && vs.Valid {
			vs.AllowedRPS = defaultAllowedRPS
		}
	}

	c.Services = ValidServices
	return Errors
}
