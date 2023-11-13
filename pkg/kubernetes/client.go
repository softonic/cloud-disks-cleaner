package kubernetes

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func NewKubernetesService(inCluster bool, k8sConfig string) (kubernetesService, error) {
	var config *rest.Config
	var err error

	if inCluster {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", k8sConfig)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get k8s config: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create k8s clientset: %w", err)
	}

	return &clientsetimpl{clientset: clientset}, nil

}

type kubernetesService interface {
	PersistentVolumesList() (*v1.PersistentVolumeList, error)
	PersistentVolumeClaimsGet(pvcName string, pvcNamespace string) (*v1.PersistentVolumeClaim, error)
}

type clientsetimpl struct {
	clientset *kubernetes.Clientset
}

func (c *clientsetimpl) PersistentVolumesList() (*v1.PersistentVolumeList, error) {
	return c.clientset.CoreV1().PersistentVolumes().List(metav1.ListOptions{})
}

func (c *clientsetimpl) PersistentVolumeClaimsGet(pvcName string, pvcNamespace string) (*v1.PersistentVolumeClaim, error) {
	return c.clientset.CoreV1().PersistentVolumeClaims(pvcNamespace).Get(pvcName, metav1.GetOptions{})
}
