package v1

import (
	"testing"

	"github.com/brevdev/cloud/internal/validation"
	v1 "github.com/brevdev/compute/pkg/v1"
)

func TestValidationFunctions(t *testing.T) {
	config := validation.ProviderConfig{
		ProviderName: "LambdaLabs",
		EnvVarName:   "LAMBDALABS_API_KEY",
		ClientFactory: func(apiKey string) v1.CloudClient {
			return NewLambdaLabsClient("validation-test", apiKey)
		},
	}

	validation.RunValidationSuite(t, config)
}

func TestInstanceLifecycleValidation(t *testing.T) {
	config := validation.ProviderConfig{
		ProviderName: "LambdaLabs",
		EnvVarName:   "LAMBDALABS_API_KEY",
		ClientFactory: func(apiKey string) v1.CloudClient {
			return NewLambdaLabsClient("validation-test", apiKey)
		},
	}

	validation.RunInstanceLifecycleValidation(t, config)
}
