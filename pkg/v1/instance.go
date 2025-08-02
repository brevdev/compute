package v1

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/alecthomas/units"
	"github.com/google/uuid"
)

type CloudCreateTerminateInstance interface {
	// CreateInstance expects an instance object to exist if successful, and no instance to exist if there is ANY error
	//      CloudClient Implementers: ensure that the instance is terminated if there is an error
	// Public ip is not always returned from create, but will exist when instance is in running state
	CreateInstance(ctx context.Context, attrs CreateInstanceAttrs) (*Instance, error)
	GetInstance(ctx context.Context, id CloudProviderInstanceID) (*Instance, error)  // may or may not be locationally scoped
	TerminateInstance(ctx context.Context, instanceID CloudProviderInstanceID) error // may or may not be locationally scoped
	ListInstances(ctx context.Context, args ListInstancesArgs) ([]Instance, error)   // return all known instances from cloud api perspective
	GetInstanceTypes(ctx context.Context, args GetInstanceTypeArgs) ([]InstanceType, error)
}

func ValidateCreateInstance(ctx context.Context, client CloudCreateTerminateInstance, attrs CreateInstanceAttrs) (*Instance, error) {
	t0 := time.Now()
	attrs.RefID = uuid.New().String()
	name, err := makeDebuggableName(attrs.Name)
	if err != nil {
		return nil, err
	}
	attrs.Name = name
	i, err := client.CreateInstance(ctx, attrs)
	if err != nil {
		return nil, err
	}
	var validationErr error
	t1 := time.Now()
	diff := t1.Sub(t0)
	if diff > 1*time.Minute {
		validationErr = errors.Join(validationErr, fmt.Errorf("create instance took too long: %s", diff))
	}
	if i.CreatedAt.Before(t0) {
		validationErr = errors.Join(validationErr, fmt.Errorf("createdAt is before t0: %s", i.CreatedAt))
	}
	if i.CreatedAt.After(t1) {
		validationErr = errors.Join(validationErr, fmt.Errorf("createdAt is after t1: %s", i.CreatedAt))
	}
	if i.Name != name {
		validationErr = errors.Join(validationErr, fmt.Errorf("name mismatch: %s != %s", i.Name, name))
	}
	if i.RefID != attrs.RefID {
		validationErr = errors.Join(validationErr, fmt.Errorf("refID mismatch: %s != %s", i.RefID, attrs.RefID))
	}
	if attrs.Location != "" && attrs.Location != i.Location {
		validationErr = errors.Join(validationErr, fmt.Errorf("location mismatch: %s != %s", attrs.Location, i.Location))
	}
	if attrs.SubLocation != "" && attrs.SubLocation != i.SubLocation {
		validationErr = errors.Join(validationErr, fmt.Errorf("subLocation mismatch: %s != %s", attrs.SubLocation, i.SubLocation))
	}
	if attrs.InstanceType != "" && attrs.InstanceType != i.InstanceType {
		validationErr = errors.Join(validationErr, fmt.Errorf("instanceType mismatch: %s != %s", attrs.InstanceType, i.InstanceType))
	}

	return i, validationErr
}

func ValidateListCreatedInstance(ctx context.Context, client CloudCreateTerminateInstance, i *Instance) error {
	ins, err := client.ListInstances(ctx, ListInstancesArgs{
		Locations: []string{i.Location},
	})
	if err != nil {
		return err
	}
	var validationErr error
	if len(ins) == 0 {
		validationErr = errors.Join(validationErr, fmt.Errorf("no instances found"))
	}
	if ins[0].Location != i.Location {
		validationErr = errors.Join(validationErr, fmt.Errorf("location mismatch: %s != %s", ins[0].Location, i.Location))
	}
	instanceIDsMap := map[CloudProviderInstanceID]Instance{}
	for _, inst := range ins {
		instanceIDsMap[inst.CloudID] = inst
	}
	inst, ok := instanceIDsMap[i.CloudID]
	if !ok {
		validationErr = errors.Join(validationErr, fmt.Errorf("instance not found: %s", i.CloudID))
		return validationErr
	}
	if inst.RefID != i.RefID {
		validationErr = errors.Join(validationErr, fmt.Errorf("refID mismatch: %s != %s", inst.RefID, i.RefID))
	}
	return validationErr
}

func ValidateTerminateInstance(ctx context.Context, client CloudCreateTerminateInstance, instance Instance) error {
	err := client.TerminateInstance(ctx, instance.CloudID)
	if err != nil {
		return err
	}
	// TODO wait for terminated
	return nil
}

type CloudStopStartInstance interface {
	StopInstance(ctx context.Context, instanceID CloudProviderInstanceID) error
	StartInstance(ctx context.Context, instanceID CloudProviderInstanceID) error
}

func ValidateStopStartInstance(ctx context.Context, client CloudStopStartInstance, instance Instance) error {
	err := client.StopInstance(ctx, instance.CloudID)
	if err != nil {
		return err
	}
	// TODO wait for stopped
	err = client.StartInstance(ctx, instance.CloudID)
	if err != nil {
		return err
	}
	// TODO wait for running
	return nil
}

type CloudRebootInstance interface {
	RebootInstance(ctx context.Context, instanceID CloudProviderInstanceID) error
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

func ValidateMergeInstanceForUpdate(client UpdateHandler, currInst Instance, newInst Instance) error {
	mergedInst := client.MergeInstanceForUpdate(currInst, newInst)

	var validationErr error
	if currInst.Name != mergedInst.Name {
		validationErr = errors.Join(validationErr, fmt.Errorf("name mismatch: %s != %s", currInst.Name, mergedInst.Name))
	}
	if currInst.RefID != mergedInst.RefID {
		validationErr = errors.Join(validationErr, fmt.Errorf("refID mismatch: %s != %s", currInst.RefID, mergedInst.RefID))
	}
	if currInst.Location != mergedInst.Location {
		validationErr = errors.Join(validationErr, fmt.Errorf("location mismatch: %s != %s", currInst.Location, newInst.Location))
	}
	if currInst.SubLocation != mergedInst.SubLocation {
		validationErr = errors.Join(validationErr, fmt.Errorf("subLocation mismatch: %s != %s", currInst.SubLocation, mergedInst.SubLocation))
	}
	if currInst.InstanceType != "" && currInst.InstanceType != mergedInst.InstanceType {
		validationErr = errors.Join(validationErr, fmt.Errorf("instanceType mismatch: %s != %s", currInst.InstanceType, mergedInst.InstanceType))
	}
	if currInst.InstanceTypeID != "" && currInst.InstanceTypeID != mergedInst.InstanceTypeID {
		validationErr = errors.Join(validationErr, fmt.Errorf("instanceTypeID mismatch: %s != %s", currInst.InstanceTypeID, mergedInst.InstanceTypeID))
	}
	if currInst.CloudCredRefID != mergedInst.CloudCredRefID {
		validationErr = errors.Join(validationErr, fmt.Errorf("cloudCredRefID mismatch: %s != %s", currInst.CloudCredRefID, mergedInst.CloudCredRefID))
	}
	if currInst.VolumeType != "" && currInst.VolumeType != mergedInst.VolumeType {
		validationErr = errors.Join(validationErr, fmt.Errorf("volumeType mismatch: %s != %s", currInst.VolumeType, mergedInst.VolumeType))
	}
	if currInst.Spot != mergedInst.Spot {
		validationErr = errors.Join(validationErr, fmt.Errorf("spot mismatch: %v != %v", currInst.Spot, mergedInst.Spot))
	}
	return validationErr
}

type Instance struct {
	Name                            string
	RefID                           string
	CloudCredRefID                  string // cloudCred used to create the Instance
	CreatedAt                       time.Time
	CloudID                         CloudProviderInstanceID
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

const (
	LifecycleStatusPending     LifecycleStatus = "pending"
	LifecycleStatusRunning     LifecycleStatus = "running"
	LifecycleStatusStopping    LifecycleStatus = "stopping"
	LifecycleStatusStopped     LifecycleStatus = "stopped"
	LifecycleStatusSuspending  LifecycleStatus = "suspending"
	LifecycleStatusSuspended   LifecycleStatus = "suspended"
	LifecycleStatusTerminating LifecycleStatus = "terminating"
	LifecycleStatusTerminated  LifecycleStatus = "terminated"
	LifecycleStatusFailed      LifecycleStatus = "failed"
)

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

func makeDebuggableName(name string) (string, error) {
	pt, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%s", name, time.Now().In(pt).Format("2006-01-02-15-04-05")), nil
}
