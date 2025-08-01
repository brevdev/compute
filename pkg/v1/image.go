package v1

import (
	"context"
	"time"
)

type CloudMachineImage interface {
	GetImages(ctx context.Context, args GetImageArgs) ([]Image, error)
}

type GetImageArgs struct {
	Owners        []string // self, amazon, aws-marketplace, <account id>, project id for GCP
	Architectures []string // i386, x86_64, arm64
	NameFilters   []string // name of the image (wildcard permitted)
	ImageIDs      []string
}

type Image struct {
	ID           string
	Architecture string
	Description  string
	Name         string
	CreatedAt    time.Time
}
