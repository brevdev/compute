package v1

type Configuration struct {
	AllowedInstanceTypes map[string][]string
}

func (c *Configuration) isAllowed(cloud string, shadeInstanceType string) bool {
	allowedInstanceTypes, found := c.AllowedInstanceTypes[cloud]
	if !found {
		return false
	}
	for _, allowedInstanceType := range allowedInstanceTypes {
		if shadeInstanceType == allowedInstanceType {
			return true
		}
	}
	return false
}
