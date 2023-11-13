package gcp

import (
	"context"
	"fmt"
)

type GCPDeleter struct {
	service   ComputeService
	projectID string
	zone      string
}

func NewGCPDeleter(client ComputeService, projectID string, zone string) (*GCPDeleter, error) {

	return &GCPDeleter{
		service:   client,
		projectID: projectID,
		zone:      zone,
	}, nil
}

func (d *GCPDeleter) DeleteResource(resourceType string, resourceID string) error {
	ctx := context.Background()
	switch resourceType {
	case "disk":
		_, err := d.service.DisksDelete(ctx, d.projectID, d.zone, resourceID)
		return err
	default:
		return fmt.Errorf("unknown resource type: %s", resourceType)
	}
}
