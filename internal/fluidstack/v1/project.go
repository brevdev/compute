package v1

import (
	"context"
)

type Project struct {
	ID    string            `json:"id"`
	Name  string            `json:"name"`
	State string            `json:"state"`
	Tags  map[string]string `json:"tags"`
}

type CreateProjectRequest struct {
	Name string            `json:"name"`
	Tags map[string]string `json:"tags,omitempty"`
}

func (c *FluidStackClient) CreateProject(_ context.Context, _ string, _ map[string]string) (*Project, error) {
	return nil, nil
}

func (c *FluidStackClient) DeleteProject(_ context.Context, _ string) error {
	return nil
}

func (c *FluidStackClient) ListProjects(_ context.Context) ([]*Project, error) {
	return nil, nil
}

func (c *FluidStackClient) GetProject(_ context.Context, _ string) (*Project, error) {
	return nil, nil
}
