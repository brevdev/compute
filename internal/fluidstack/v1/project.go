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

func (c *FluidStackClient) CreateProject(ctx context.Context, name string, tags map[string]string) (*Project, error) {
	return nil, nil
}

func (c *FluidStackClient) DeleteProject(ctx context.Context, projectID string) error {
	return nil
}

func (c *FluidStackClient) ListProjects(ctx context.Context) ([]*Project, error) {
	return nil, nil
}

func (c *FluidStackClient) GetProject(ctx context.Context, projectID string) (*Project, error) {
	return nil, nil
}
