package v1

import (
	"os"
	"testing"

	"github.com/brevdev/cloud/internal/validation"
)

func TestValidationFunctions(t *testing.T) {
	apiKey := os.Getenv("LAMBDALABS_API_KEY")
	if apiKey == "" {
		t.Skip("LAMBDALABS_API_KEY not set, skipping LambdaLabs validation tests")
	}

	config := validation.ProviderConfig{
		Credential: NewLambdaLabsCredential("validation-test", apiKey),
	}

	validation.RunValidationSuite(t, config)
}

func TestInstanceLifecycleValidation(t *testing.T) {
	apiKey := os.Getenv("LAMBDALABS_API_KEY")
	if apiKey == "" {
		t.Skip("LAMBDALABS_API_KEY not set, skipping LambdaLabs validation tests")
	}

	config := validation.ProviderConfig{
		Credential: NewLambdaLabsCredential("validation-test", apiKey),
	}

	validation.RunInstanceLifecycleValidation(t, config)
}
