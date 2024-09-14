package config

import (
	"errors"
	"strings"

	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/v2"
)

// Global variable but private
var config *koanf.Koanf = nil

// Load config from environment.
// Every env var is coverted to lowercase and plitted by underscore "_".
//
// Example: `BACKEND_URL` becomes `backend.url`
func LoadConfig() error {
	k := koanf.New(".")

	if err := k.Load(env.Provider("", ".", func(s string) string {
		return strings.Replace(strings.ToLower(s), "_", ".", -1)
	}), nil); err != nil {
		return err
	}

	config = k
	return nil
}

// Return the instance or error if the config is not laoded yet
func GetConfig() (*koanf.Koanf, error) {
	if config == nil {
		return nil, errors.New("You must call `LoadConfig()` first.")
	}
	return config, nil
}
