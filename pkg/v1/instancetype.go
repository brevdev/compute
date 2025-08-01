package v1

import (
	"time"

	"github.com/alecthomas/units"
	"github.com/bojanz/currency"
)

type InstanceTypeID string

type InstanceType struct {
	ID                              InstanceTypeID // this id should be unique across all regions and stable
	Location                        string
	AvailableAzs                    []string
	SubLocation                     string
	Type                            string
	SupportedGPUs                   []GPU
	SupportedStorage                []Storage
	ElasticRootVolume               bool
	SupportedUsageClasses           []string
	Memory                          units.Base2Bytes
	MaximumNetworkInterfaces        int32
	NetworkPerformance              string
	SupportedNumCores               []int32
	DefaultCores                    int32
	VCPU                            int32
	SupportedArchitectures          []string
	ClockSpeedInGhz                 float64
	Quota                           InstanceTypeQuota
	Stoppable                       bool
	Rebootable                      bool
	VariablePrice                   bool
	Preemptible                     bool
	IsAvailable                     bool
	BasePrice                       *currency.Amount
	SubLocationTypeChangeable       bool
	IsContainer                     bool
	UserPrivilegeEscalationDisabled bool
	NotPrivileged                   bool
	EstimatedDeployTime             *time.Duration
	Provider                        string
	CanModifyFirewallRules          bool
}

type GPU struct {
	Count          int32
	Memory         units.Base2Bytes
	MemoryDetails  string
	NetworkDetails string
	Manufacturer   string
	Name           string
	Type           string
}

type InstanceTypeQuota struct {
	OnDemand Quota
	Spot     Quota
	Reserved Quota
}

type GetInstanceTypeArgs struct {
	Locations              LocationsFilter
	SupportedArchitectures []string
	InstanceTypes          []string
}
