package kubernetes

import (
	"errors"
	"strings"

	v1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	_ "k8s.io/apimachinery/pkg/api/resource"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/klog"
)

// K8sChecker es una estructura que implementa la interfaz Checker para Kubernetes.
type K8sChecker struct {
	//clientset kubernetes.Interface // Cliente de Kubernetes para interactuar con el cluster.
	clientset kubernetesService
}

// NewK8sChecker crea una nueva instancia de K8sChecker.
func NewK8sChecker(clientset kubernetesService) (*K8sChecker, error) {
	return &K8sChecker{
		clientset: clientset,
	}, nil
}

// Verifica si un disco en Kubernetes está en uso o no.
// Si no existe el PV, el recurso no está en uso y se puede borrar
// Si existe el PV, pero la vincuacion PV contra PVC no es bidireccional
// O sea, que si el pvc asociado al PV , apunta a su vez a otro pv, el disco no está en uso y
// se puede borrar
// Si existe el PV y el PVC al que está asociado tambien existe, aunque el estado sea Released
// no borraremos el disco, como medida de precaucion.
func (c *K8sChecker) IsResourceUnused(resourceID string) (bool, error) {
	// En este caso, asumiremos que resourceID es el nombre del PersistentVolume en Kubernetes.
	pv, err := c.getPVFromDisk(resourceID)
	if err != nil {
		return false, err
	}

	if pv == nil {
		klog.Infof("PV was not found. Disk => %s", resourceID)
		return true, nil
	}

	// Verifica si el PersistentVolume está en uso.
	// Un PV se considera en uso si tiene una PersistentVolumeClaim asociada.
	/* 	isUnused := pv.Status.Phase == "Released"
	   	if isUnused {
	   		fmt.Println("PV was found, but pvc does not. Disk => ", resourceID)
	   	} */

	pvc, err := c.getPVCFromPV(pv)
	if err != nil {
		if k8serrors.IsNotFound(err) {
			// El PVC no se encontró, por lo que el recurso no está siendo utilizado.
			// Pero por seguridad no vamos a eliminar el disco. Quiza será utilizado en el futuro.
			klog.Infof("The pvc bind to the disk: %s does not exists", resourceID)
			return false, nil
		}
		return false, err
	}

	isBound, err := c.isPVCBoundToPV(pv.Name, pvc.Namespace, pvc.Name)
	if err != nil {
		return false, err
	}

	isUnused := !isBound // Invierte la lógica: si no está vinculado, entonces no está en uso.
	if isUnused {
		klog.Infof("PV was found, PVC was found, but pvc and pv are not binded. Disk => %s", resourceID)
	}

	return isUnused, nil // Retorna el resultado de la verificación y nil para el error.
}

func (c *K8sChecker) ListResources() ([]interface{}, error) {
	// Retornar nil o una lista vacía, o tal vez un error que indique que esta operación no está soportada.
	return nil, nil
}

func (c *K8sChecker) getPVFromDisk(diskName string) (*v1.PersistentVolume, error) {
	if c == nil {
		return nil, errors.New("clientset is nil")
	}

	pvList, err := c.clientset.PersistentVolumesList()
	if err != nil {
		return nil, err
	}

	for _, pv := range pvList.Items {
		if pv.Spec.GCEPersistentDisk != nil && pv.Spec.GCEPersistentDisk.PDName == diskName {
			klog.Infof("Entra en gce persistent disk. PV => %s", pv.Name)
			return &pv, nil // return the PV found for the GCE traditional disk
		}

		// Para el caso de un disco CSI, el identificador del disco podría ser parte de volumeHandle.
		// La comparación es más relajada aquí, solo comprobamos si el nombre del disco es parte de volumeHandle.
		if pv.Spec.CSI != nil && strings.Contains(pv.Spec.CSI.VolumeHandle, diskName) {
			klog.Infof("Entra en CSI PV => %s", pv.Name)
			return &pv, nil // return the PV found for the CSI disk
		}
	}

	return nil, nil
}

// GetPVCFromPV obtiene el PVC asociado a un PV.
func (c *K8sChecker) getPVCFromPV(pv *v1.PersistentVolume) (*v1.PersistentVolumeClaim, error) {
	if pv.Spec.ClaimRef == nil {
		return nil, nil
	}

	namespace := pv.Spec.ClaimRef.Namespace
	name := pv.Spec.ClaimRef.Name

	pvc, err := c.clientset.PersistentVolumeClaimsGet(name, namespace)
	if err != nil {
		// Error al obtener el PVC.
		return nil, err
	}

	return pvc, nil
}

func (c *K8sChecker) isPVCBoundToPV(pvName string, pvcNamespace string, pvcName string) (bool, error) {
	// Obtén el PVC del clúster
	pvc, err := c.clientset.PersistentVolumeClaimsGet(pvcName, pvcNamespace)
	if err != nil {
		return false, err
	}

	// Compara el nombre del PV con el valor de spec.volumeName en el PVC
	return pvc.Spec.VolumeName == pvName, nil
}
