package gcp

type Deleter interface {
	DeleteResource(resourceType string, resourceID string) error
}
