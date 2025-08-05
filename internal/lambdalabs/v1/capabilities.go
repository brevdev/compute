package v1

import (
	"context"

	v1 "github.com/brevdev/compute/pkg/v1"
)

// GetCapabilities returns the capabilities of Lambda Labs
// Based on API documentation at https://cloud.lambda.ai/api/v1/openapi.json
func (c *LambdaLabsClient) GetCapabilities(_ context.Context) (v1.Capabilities, error) {
	capabilities := v1.Capabilities{
		// SUPPORTED FEATURES (with API evidence):

		// Instance Management
		v1.CapabilityCreateInstance,          // POST /api/v1/instance-operations/launch
		v1.CapabilityTerminateInstance,       // POST /api/v1/instance-operations/terminate
		v1.CapabilityCreateTerminateInstance, // Combined create/terminate capability
		v1.CapabilityRebootInstance,          // POST /api/v1/instance-operations/restart

		// Firewall Management
		v1.CapabilityModifyFirewall, // Firewall rulesets API available

		// UNSUPPORTED FEATURES (no API evidence found):
		// - v1.CapabilityStopStartInstance     // No stop/start endpoints
		// - v1.CapabilityResizeInstanceVolume  // No volume resizing endpoints
		// - v1.CapabilityMachineImage          // No image endpoints
		// - v1.CapabilityTags                  // No tagging endpoints
	}

	return capabilities, nil
}
