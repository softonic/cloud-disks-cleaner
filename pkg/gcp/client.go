package gcp

import (
	"context"

	"google.golang.org/api/compute/v1"
)

type ComputeService interface {
	// Asegúrate de incluir aquí todos los métodos que usas de compute.Service.
	//NewService(ctx context.Context, opts ...option.ClientOption) (*compute.Service, error)
	DisksGet(ctx context.Context, projectID string, zone string, diskID string) (*compute.Disk, error)
	DisksDelete(ctx context.Context, projectID string, zone string, diskID string) (*compute.Operation, error)
	DisksList(projectID string, zone string) *compute.DisksListCall
}

type computeServiceImpl struct {
	service *compute.Service
}

func (c *computeServiceImpl) DisksGet(ctx context.Context, projectID string, zone string, diskID string) (*compute.Disk, error) {
	return c.service.Disks.Get(projectID, zone, diskID).Context(ctx).Do()
}

func (c *computeServiceImpl) DisksDelete(ctx context.Context, projectID string, zone string, diskID string) (*compute.Operation, error) {
	return c.service.Disks.Delete(projectID, zone, diskID).Context(ctx).Do()
}

func (c *computeServiceImpl) DisksList(projectID string, zone string) *compute.DisksListCall {
	return c.service.Disks.List(projectID, zone)
}

func NewComputeService(ctx context.Context) (ComputeService, error) {
	client, err := compute.NewService(ctx)
	if err != nil {
		return nil, err
	}
	return &computeServiceImpl{service: client}, nil
}
