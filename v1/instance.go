package v1

import (
	"context"
	"time"

	"github.com/alecthomas/units"
)

type CloudCreateTerminateInstance interface {
	// CreateInstance expects an instance object to exist if successful, and no instance to exist if there is ANY error
	//      CloudClient Implementers: ensure that the instance is terminated if there is an error
	// Public ip is not always returned from create, but will exist when instance is in running state
	CreateInstance(ctx context.Context, attrs CreateInstanceAttrs) (*Instance, error)
	GetInstance(ctx context.Context, id string) (*Instance, error)                 // may or may not be locationally scoped
	TerminateInstance(ctx context.Context, instanceID string) error                // may or may not be locationally scoped
	ListInstances(ctx context.Context, args ListInstancesArgs) ([]Instance, error) // return all known instances from cloud api perspective
	GetInstanceTypes(ctx context.Context, args GetInstanceTypeArgs) ([]InstanceType, error)
}

type CloudRebootInstance interface {
	RebootInstance(ctx context.Context, instanceID CloudProviderInstanceID) error
}

type CloudStopStartInstance interface {
	StopInstance(ctx context.Context, instanceID CloudProviderInstanceID) error
	StartInstance(ctx context.Context, instanceID CloudProviderInstanceID) error
}

type CloudChangeInstanceType interface {
	ChangeInstanceType(ctx context.Context, instanceID CloudProviderInstanceID, instanceType string) error
}

type CloudInstanceTags interface {
	UpdateInstanceTags(ctx context.Context, args UpdateInstanceTagsArgs) error
}

// this is used by the control plane to efficiently update instances
type UpdateHandler interface {
	MergeInstanceForUpdate(currInst Instance, newInst Instance) Instance
	MergeInstanceTypeForUpdate(currIt InstanceType, newIt InstanceType) InstanceType
}

type Instance struct {
	Name                            string
	RefID                           string
	CloudCredRefID                  string // cloudCred used to create the Instance
	CreatedAt                       time.Time
	CloudID                         string
	IPAllocationID                  *string
	PublicIP                        string // Public ip is not always returned from create, but will exist when instance is in running state
	PublicDNS                       string
	PrivateIP                       string
	Hostname                        string
	ImageID                         string
	InstanceType                    string
	DiskSize                        units.Base2Bytes
	VolumeType                      string
	PubKeyFingerprint               string
	SSHUser                         string
	SSHPort                         int
	Status                          Status
	MetaEndpointEnabled             bool
	MetaTagsEnabled                 bool
	VPCID                           string
	SubnetID                        string
	Spot                            bool
	FirewallRules                   FirewallRules
	RetiredAt                       *time.Time
	RetireTimeout                   *time.Duration
	LastStopTransitionTime          *time.Time
	Location                        string
	SubLocation                     string
	Tags                            Tags
	Stoppable                       bool
	Rebootable                      bool
	IsContainer                     bool
	UserPrivilegeEscalationDisabled bool
	NotPrivileged                   bool
	InstanceTypeID                  InstanceTypeID
	AdditionalDisks                 []Disk

	// As of 08/26/2024 only used for Launchpad cloud.
	// Because there is port forwarding from a GPU node to Bastion node,
	// there is port mappings from the GPU node itself to the Bastion node.
	// i.e. Verb SSH port 2222 is mapped to 2022 on the Bastion node
	InternalPortMappings []PortMapping
}

type Status struct {
	LifecycleStatus LifecycleStatus
	Messages        []string
}

type LifecycleStatus string

type CloudProviderInstanceID string

type ListInstancesArgs struct {
	InstanceIDs []CloudProviderInstanceID
	TagFilters  map[string][]string
	Locations   LocationsFilter
}

type CreateInstanceAttrs struct {
	Location             string
	SubLocation          string
	Name                 string
	RefID                string // required also can be used for idempotency
	VPCID                string
	SubnetID             string
	PublicKey            string // must be in openssh format
	KeyPairName          *string
	ImageID              string
	InstanceType         string
	UserDataBase64       string
	DiskSize             units.Base2Bytes
	Tags                 Tags
	FirewallRules        FirewallRules
	UseSpot              bool
	UsePersistentIP      bool
	UseMultiAttachVolume bool
	RetireTimeout        *time.Duration
	// Additional Environment Variables.
	// Note: As of May 2024, the only cloud provider we have this implemented for
	// is the Akash provider.
	AdditionalEnvVars map[string]string
	AdditionalDisks   []Disk
}

type UpdateInstanceTagsArgs struct {
	InstanceID CloudProviderInstanceID
	Tags       Tags
}
