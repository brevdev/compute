package config

import (
	"reflect"
	"testing"
)

func TestRegisteredProvidersHaveRefIDField(t *testing.T) {
	for providerID, factory := range providerRegistry {
		t.Run(providerID, func(t *testing.T) {
			cred := factory()
			credType := reflect.TypeOf(cred)

			if credType.Kind() == reflect.Ptr {
				credType = credType.Elem()
			}

			hasRefID := false
			for i := 0; i < credType.NumField(); i++ {
				field := credType.Field(i)
				if field.Name == "RefID" {
					hasRefID = true
					break
				}
				if field.Anonymous && field.Type.Kind() == reflect.Ptr {
					embeddedType := field.Type.Elem()
					for j := 0; j < embeddedType.NumField(); j++ {
						embeddedField := embeddedType.Field(j)
						if embeddedField.Name == "RefID" {
							hasRefID = true
							break
						}
					}
				}
			}

			if !hasRefID {
				t.Errorf("Provider %s does not have a RefID field in its credential struct", providerID)
			}
		})
	}
}
