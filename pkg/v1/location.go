package v1

import "context"

type CloudLocation interface {
	GetLocations(ctx context.Context, args GetLocationsArgs) ([]Location, error)
}

type GetLocationsArgs struct {
	IncludeUnavailable bool
}

type Location struct {
	Name        string // basically the id
	Description string
	Available   bool
	Endpoint    string
	Priority    int
	Country     string // ISO 3166-1 alpha-3 https://en.wikipedia.org/wiki/ISO_3166-1_alpha-3
}

type LocationsFilter []string
