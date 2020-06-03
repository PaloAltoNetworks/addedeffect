package appcreds

import "time"

type config struct {
	subnets     []string
	maxValidity time.Duration
}

func newConfig() config {
	return config{}
}

// An Optin can be used to configure a new appcred.
type Option func(*config)

// OptionSubnets configures the appcred to only
// work when used from one of the provided subnet.
func OptionSubnets(subnets []string) Option {
	return func(c *config) {
		c.subnets = subnets
	}
}

// OptionMaxValidity configures the appcred to only capable
// of delivering token with the provided max validity.
func OptionMaxValidity(max time.Duration) Option {
	return func(c *config) {
		c.maxValidity = max
	}
}
