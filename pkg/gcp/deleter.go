package gcp

import (
	"context"
	"fmt"

	"google.golang.org/api/compute/v1"
)

type GCPDeleter struct {
	service   *compute.Service
	projectID string
	zone      string
}

func NewGCPDeleter(client compute.Service, projectID string, zone string) (*GCPDeleter, error) {

	return &GCPDeleter{
		service:   &client,
		projectID: projectID,
		zone:      zone,
	}, nil
}

func (d *GCPDeleter) DeleteResource(resourceType string, resourceID string) error {
	ctx := context.Background()
	switch resourceType {
	case "disk":
		_, err := d.service.Disks.Delete(d.projectID, d.zone, resourceID).Context(ctx).Do()
		return err
	default:
		return fmt.Errorf("unknown resource type: %s", resourceType)
	}
}
