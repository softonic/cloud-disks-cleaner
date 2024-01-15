package usage

type Checker interface {
	IsResourceUnused(resourceID string) (bool, string, error)
	ListResources() ([]interface{}, error)
}
