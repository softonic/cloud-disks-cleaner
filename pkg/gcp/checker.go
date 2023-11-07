package gcp

import (
	"context"
	"fmt"

	"google.golang.org/api/compute/v1"
)

// GCPChecker es una estructura que implementa la interfaz Checker para GCP.
type GCPChecker struct {
	service   *compute.Service // servicio de GCP para interactuar con los recursos de GCP.
	projectID string           // ID del proyecto de GCP.
	zone      string           // Zona de GCP en la que se deben verificar los recursos.
}

// NewGCPChecker crea una nueva instancia de GCPChecker.
func NewGCPChecker(client compute.Service, projectID string, zone string) (*GCPChecker, error) {

	return &GCPChecker{
		service:   &client,
		projectID: projectID,
		zone:      zone,
	}, nil
}

func (c *GCPChecker) IsResourceUnused(resourceID string) (bool, error) {
	ctx := context.Background()

	disk, err := c.service.Disks.Get(c.projectID, c.zone, resourceID).Context(ctx).Do()
	if err != nil {
		return false, err // Si hay un error, retorna false y el error.
	}

	// Verifica si el disco está en uso.
	// En este ejemplo, supondremos que un disco se considera en uso si tiene usuarios asociados.
	isUnused := len(disk.Users) == 0

	return isUnused, nil // Retorna el resultado de la verificación y nil para el error.
}

func (c *GCPChecker) ListResources() ([]interface{}, error) {
	ctx := context.Background()

	req := c.service.Disks.List(c.projectID, c.zone)
	var disks []compute.Disk
	var disksNames []string
	var resources []interface{}
	if err := req.Pages(ctx, func(page *compute.DiskList) error {
		for _, disk := range page.Items {
			disks = append(disks, *disk)
			resources = append(resources, disk)
			disksNames = append(disksNames, disk.Name)
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("Failed to list disks: %w", err)
	}

	return resources, nil
}
