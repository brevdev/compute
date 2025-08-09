package providers

import (
	"fmt"

	"github.com/brevdev/cloud/cmd/bcloud/config"
	v1 "github.com/brevdev/cloud/pkg/v1"
)

type DefaultLocationProvider interface {
	GetDefaultLocation() string
}

func CreateCredential(_ string, credEntry config.CredentialEntry) (v1.CloudCredential, error) {
	if credEntry.Value == nil {
		return nil, fmt.Errorf("credential entry has no value")
	}

	return credEntry.Value, nil
}

func GetDefaultLocation(cred v1.CloudCredential) string {
	if provider, ok := cred.(DefaultLocationProvider); ok {
		return provider.GetDefaultLocation()
	}
	return ""
}
