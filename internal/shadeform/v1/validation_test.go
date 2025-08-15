package v1

import (
	"github.com/brevdev/cloud/internal/validation"
	v1 "github.com/brevdev/cloud/pkg/v1"
	"os"
	"testing"
)

func TestValidationFunctions(t *testing.T) {
	checkSkip(t)
	apiKey := getAPIKey()

	config := validation.ProviderConfig{
		Credential: NewShadeformCredential("validation-test", apiKey),
		StableIDs:  []v1.InstanceTypeID{"datacrunch_B200_helsinki-finland-5", "massedcompute_L40_desmoines-usa-1"},
	}

	validation.RunValidationSuite(t, config)
}

func TestInstanceLifecycleValidation(t *testing.T) {
	checkSkip(t)
	apiKey := getAPIKey()

	config := validation.ProviderConfig{
		Credential: NewShadeformCredential("validation-test", apiKey),
	}

	validation.RunInstanceLifecycleValidation(t, config)
}

func checkSkip(t *testing.T) {
	apiKey := getAPIKey()
	isValidationTest := os.Getenv("VALIDATION_TEST")
	if apiKey == "" && isValidationTest != "" {
		t.Fatal("SHADEFORM_API_KEY not set, but VALIDATION_TEST is set")
	} else if apiKey == "" && isValidationTest == "" {
		t.Skip("SHADEFORM_API_KEY not set, skipping shadeform validation tests")
	}
}

func getAPIKey() string {
	return os.Getenv("SHADEFORM_API_KEY")
}
