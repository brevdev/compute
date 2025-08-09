# How to Add a Cloud Provider

This guide walks you through implementing a new cloud provider for the Brev Cloud SDK, using the Lambda Labs implementation as a reference example.

## Overview

Adding a new cloud provider involves implementing the `CloudClient` interface and its constituent interfaces. The SDK uses a composition pattern where `CloudClient` aggregates multiple smaller interfaces, allowing providers to implement only the features they support using the `NotImplCloudClient` pattern.

## Prerequisites

Before implementing a provider, ensure you have:

1. **API Documentation**: Complete API documentation for your cloud provider
2. **Authentication Method**: Understanding of how to authenticate with the provider's API
3. **Feature Analysis**: Clear understanding of which CloudClient features the provider supports
4. **Security Requirements**: Familiarity with [Security Requirements](SECURITY.md)

## Implementation Steps

### 1. Directory Structure

Create your provider directory following the established pattern:

```
internal/{provider-name}/
├── v1/
│   ├── client.go          # Credential and client implementation
│   ├── instance.go        # Instance lifecycle operations
│   ├── capabilities.go    # Feature capability declarations
│   ├── instancetype.go    # Instance type management (if supported)
│   ├── location.go        # Location management (if supported)
│   ├── networking.go      # Firewall/networking (if supported)
│   ├── errors.go          # Provider-specific error handling
│   └── README.md          # Provider-specific documentation
├── CONTRIBUTE.md          # Development setup and contribution guide
└── SECURITY.md            # Security implementation details
```

### 2. Core Interface Implementation

#### 2.1 Credential Implementation

Implement the `CloudCredential` interface in `client.go`:

```go
// {Provider}Credential implements the CloudCredential interface
type {Provider}Credential struct {
    RefID  string
    APIKey string  // or other auth fields
}

var _ v1.CloudCredential = &{Provider}Credential{}

func (c *{Provider}Credential) GetReferenceID() string {
    return c.RefID
}

func (c *{Provider}Credential) GetAPIType() v1.APIType {
    return v1.APITypeGlobal  // or APITypeLocational
}

func (c *{Provider}Credential) GetCloudProviderID() v1.CloudProviderID {
    return "{provider-id}"
}

func (c *{Provider}Credential) GetTenantID() (string, error) {
    // Implementation specific to your provider
}

func (c *{Provider}Credential) MakeClient(ctx context.Context, location string) (v1.CloudClient, error) {
    return New{Provider}Client(c.RefID, c.APIKey), nil
}
```

#### 2.2 Client Implementation

Implement the main client struct using the `NotImplCloudClient` pattern:

```go
// {Provider}Client implements the CloudClient interface
// It embeds NotImplCloudClient to handle unsupported features
type {Provider}Client struct {
    v1.NotImplCloudClient
    refID    string
    apiKey   string
    client   *{provider-sdk}.APIClient
    location string
}

var _ v1.CloudClient = &{Provider}Client{}

func New{Provider}Client(refID, apiKey string) *{Provider}Client {
    // Initialize your provider's SDK client
    return &{Provider}Client{
        refID:  refID,
        apiKey: apiKey,
        client: // your provider SDK client
    }
}
```

### 3. Instance Lifecycle Operations

Implement core instance operations in `instance.go`:

```go
// CreateInstance creates a new instance
func (c *{Provider}Client) CreateInstance(ctx context.Context, attrs v1.CreateInstanceAttrs) (*v1.Instance, error) {
    // 1. Handle SSH key setup if needed
    // 2. Convert Brev attributes to provider-specific request
    // 3. Call provider API to create instance
    // 4. Convert provider response to v1.Instance
    // 5. Return instance or error
}

// GetInstance retrieves an instance by ID
func (c *{Provider}Client) GetInstance(ctx context.Context, instanceID v1.CloudProviderInstanceID) (*v1.Instance, error) {
    // Implementation
}

// TerminateInstance terminates an instance
func (c *{Provider}Client) TerminateInstance(ctx context.Context, instanceID v1.CloudProviderInstanceID) error {
    // Implementation
}

// ListInstances lists all instances
func (c *{Provider}Client) ListInstances(ctx context.Context, args v1.ListInstancesArgs) ([]v1.Instance, error) {
    // Implementation
}
```

#### 3.1 Instance Conversion

Create conversion functions to map between provider types and v1 types:

```go
func convert{Provider}InstanceToV1Instance(providerInstance {provider}.Instance) *v1.Instance {
    return &v1.Instance{
        RefID:          // extract from provider instance
        CloudCredRefID: // extract from instance name or metadata
        CreatedAt:      // convert provider timestamp
        CloudID:        v1.CloudProviderInstanceID(providerInstance.ID),
        Name:           providerInstance.Name,
        PublicIP:       providerInstance.PublicIP,
        PrivateIP:      providerInstance.PrivateIP,
        Status: v1.Status{
            LifecycleStatus: convert{Provider}StatusToV1Status(providerInstance.Status),
        },
        InstanceType:  providerInstance.InstanceType,
        Location:      providerInstance.Region,
        SSHUser:       "ubuntu", // or provider default
        SSHPort:       22,
        Stoppable:     // based on provider capabilities
        Rebootable:    // based on provider capabilities
        // ... other fields
    }
}

func convert{Provider}StatusToV1Status(status string) v1.LifecycleStatus {
    switch status {
    case "running":
        return v1.LifecycleStatusRunning
    case "pending":
        return v1.LifecycleStatusPending
    case "terminated":
        return v1.LifecycleStatusTerminated
    // ... other status mappings
    default:
        return v1.LifecycleStatusPending
    }
}
```

### 4. Capabilities Declaration

Define your provider's capabilities in `capabilities.go`:

```go
func get{Provider}Capabilities() v1.Capabilities {
    return v1.Capabilities{
        // SUPPORTED FEATURES (with API evidence):
        v1.CapabilityCreateInstance,          // POST /api/v1/instances
        v1.CapabilityTerminateInstance,       // DELETE /api/v1/instances/{id}
        v1.CapabilityCreateTerminateInstance, // Combined capability
        v1.CapabilityRebootInstance,          // POST /api/v1/instances/{id}/reboot
        
        // UNSUPPORTED FEATURES (no API evidence found):
        // - v1.CapabilityStopStartInstance     // No stop/start endpoints
        // - v1.CapabilityResizeInstanceVolume  // No volume resizing
        // - v1.CapabilityMachineImage          // No image management
        // - v1.CapabilityTags                  // No tagging support
    }
}

func (c *{Provider}Client) GetCapabilities(ctx context.Context) (v1.Capabilities, error) {
    return get{Provider}Capabilities(), nil
}

func (c *{Provider}Credential) GetCapabilities(ctx context.Context) (v1.Capabilities, error) {
    return get{Provider}Capabilities(), nil
}
```

### 5. Optional Interface Implementations

Implement additional interfaces based on your provider's capabilities:

#### 5.1 Instance Types (if supported)

```go
// GetInstanceTypes retrieves available instance types
func (c *{Provider}Client) GetInstanceTypes(ctx context.Context, args v1.GetInstanceTypeArgs) ([]v1.InstanceType, error) {
    // Implementation
}
```

#### 5.2 Locations (if supported)

```go
// GetLocations retrieves available locations/regions
func (c *{Provider}Client) GetLocations(ctx context.Context, args v1.GetLocationsArgs) ([]v1.Location, error) {
    // Implementation
}
```

#### 5.3 Additional Operations (if supported)

```go
// RebootInstance reboots an instance
func (c *{Provider}Client) RebootInstance(ctx context.Context, instanceID v1.CloudProviderInstanceID) error {
    // Implementation
}

// StopInstance stops an instance (if supported)
func (c *{Provider}Client) StopInstance(ctx context.Context, instanceID v1.CloudProviderInstanceID) error {
    // Implementation
}

// StartInstance starts a stopped instance (if supported)
func (c *{Provider}Client) StartInstance(ctx context.Context, instanceID v1.CloudProviderInstanceID) error {
    // Implementation
}
```

### 6. Error Handling

Create provider-specific error handling in `errors.go`:

```go
func handleErrToCloudErr(err error) error {
    // Convert provider-specific errors to appropriate v1 errors
    // Handle rate limiting, authentication, not found, etc.
    return err
}
```

### 7. Update Handlers

Implement the `UpdateHandler` interface:

```go
func (c *{Provider}Client) MergeInstanceForUpdate(currInst v1.Instance, newInst v1.Instance) v1.Instance {
    // Merge logic specific to your provider
    // Usually just return newInst unless special handling needed
    return newInst
}

func (c *{Provider}Client) MergeInstanceTypeForUpdate(currIt v1.InstanceType, newIt v1.InstanceType) v1.InstanceType {
    return newIt
}
```

## Documentation Requirements

### 1. Provider README (`internal/{provider}/v1/README.md`)

Document your provider implementation:

```markdown
# {Provider} Provider

## Overview
Brief description of the provider and its API integration.

## Supported Features
List all supported CloudClient interface methods with API endpoint references:
- ✅ **Create Instance**: `POST /api/v1/instances`
- ✅ **Get Instance**: `GET /api/v1/instances/{id}`

## Unsupported Features
List unsupported features with explanations:
- ❌ **Stop/Start Instance**: No stop/start endpoints available

## Implementation Approach
Explain use of NotImplCloudClient pattern and interface compliance.

## API Integration
Document base URL, authentication method, and API documentation links.
```

### 2. Security Documentation (`internal/{provider}/SECURITY.md`)

Document security implementation following [Security Requirements](SECURITY.md):

```markdown
# {Provider} Security Implementation

## SSH Access Requirements
Document SSH server support and key-based authentication.

## Network Security
Explain firewall implementation and network isolation.

## Data Protection
Document encryption at rest and in transit.

## Implementation Checklist
- [ ] Inbound Deny by Default
- [ ] Outbound Allow by Default
- [ ] Explicit Inbound Access
- [ ] Persistent Disk Encryption
- [ ] TLS 1.2+ enforcement
```

### 3. Contribution Guide (`internal/{provider}/CONTRIBUTE.md`)

Document development setup and testing:

```markdown
# Contributing to {Provider}

## Setup
Environment setup instructions, API key configuration.

## Running Tests
Test execution instructions.

## Implementation Notes
Provider-specific development notes and gotchas.
```

## Testing Strategy

### 1. Unit Tests

Create comprehensive unit tests for your implementation:

```go
func TestCreateInstance(t *testing.T) {
    // Test instance creation with various scenarios
}

func TestGetInstance(t *testing.T) {
    // Test instance retrieval
}

func TestCapabilities(t *testing.T) {
    // Test capability declarations
}
```

### 2. Integration Tests

Use the SDK's validation functions:

```go
func TestValidateCreateInstance(t *testing.T) {
    client := New{Provider}Client(refID, apiKey)
    _, err := v1.ValidateCreateInstance(ctx, client, attrs)
    assert.NoError(t, err)
}
```

### 3. GitHub Actions

Set up CI/CD following the Lambda Labs pattern in `.github/workflows/`.

## Key Implementation Patterns

### 1. NotImplCloudClient Pattern

- Embed `v1.NotImplCloudClient` in your client struct
- Only implement methods your provider actually supports
- Unsupported methods automatically return `ErrNotImplemented`
- Maintains full `CloudClient` interface compliance

### 2. Capability-Based Design

- Declare only capabilities your provider actually supports
- Use evidence-based decisions from API documentation
- Comment unsupported features with explanations

### 3. Error Handling

- Convert provider-specific errors to appropriate SDK errors
- Handle common scenarios: rate limiting, authentication, not found
- Provide meaningful error messages

### 4. Type Conversion

- Create conversion functions between provider and v1 types
- Handle status mapping carefully
- Ensure all required v1.Instance fields are populated

## Common Pitfalls

1. **Over-implementing**: Don't implement features your provider doesn't support
2. **Status Mapping**: Ensure accurate lifecycle status conversion
3. **Authentication**: Handle API authentication and token refresh properly
4. **Error Handling**: Don't ignore provider-specific error conditions
5. **Testing**: Validate with real API calls, not just unit tests

## Reference Implementation

See the [Lambda Labs implementation](../internal/lambdalabs/) for a complete, well-tested example that demonstrates all these patterns.

## Related Documentation

- [Architecture Overview](ARCHITECTURE.md) - How the Cloud SDK fits into Brev's architecture
- [Security Requirements](SECURITY.md) - Security specifications and implementation requirements
- [V1 Design Notes](../pkg/v1/V1_DESIGN_NOTES.md) - Design decisions and AWS-inspired patterns

## Getting Help

For questions about provider implementation:
- Review existing provider implementations in `internal/`
- Check the core interfaces in `pkg/v1/`
- Refer to validation functions for testing patterns
