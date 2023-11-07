package gcp

import (
	"context"

	"google.golang.org/api/compute/v1"
)

// NewClient creates a new GCP client.
func NewClient() (*compute.Service, error) {
	// implicit uses Application Default Credentials to authenticate
	ctx := context.Background()

	computeService, err := compute.NewService(ctx)
	if err != nil {
		return nil, err

	}
	return computeService, nil
}
