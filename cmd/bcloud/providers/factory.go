package providers

import (
	"fmt"

	lambdalabs "github.com/brevdev/cloud/internal/lambdalabs/v1"
	v1 "github.com/brevdev/cloud/pkg/v1"
)

func CreateCredential(yamlKey string, config map[string]interface{}) (v1.CloudCredential, error) {
	provider, ok := config["provider"].(string)
	if !ok {
		return nil, fmt.Errorf("provider field is required")
	}

	var refID string
	if explicitRefID, exists := config["ref_id"].(string); exists && explicitRefID != "" {
		refID = explicitRefID
	} else {
		refID = yamlKey
	}

	switch provider {
	case "lambdalabs":
		apiKey, ok := config["api_key"].(string)
		if !ok {
			return nil, fmt.Errorf("api_key required for lambdalabs provider")
		}
		return lambdalabs.NewLambdaLabsCredential(refID, apiKey), nil

	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
}
