package usage

type Checker interface {
	IsResourceUnused(resourceID string) (bool, error)
	ListResources() ([]interface{}, error)
}
