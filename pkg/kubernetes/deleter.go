package kubernetes

type K8sDeleter struct {
	clientset kubernetesService
}

func NewK8sDeleter(clientset kubernetesService) (*K8sDeleter, error) {
	return &K8sDeleter{
		clientset: clientset,
	}, nil
}

func (d *K8sDeleter) DeleteResource(resourceID string) error {
	return d.clientset.PersistentVolumeDelete(resourceID)
}
