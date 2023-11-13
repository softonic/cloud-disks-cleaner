package app

import (
	"github.com/softonic/cloud-disks-cleaner/pkg/gcp"
	"github.com/softonic/cloud-disks-cleaner/pkg/usage"
	"google.golang.org/api/compute/v1"
	"k8s.io/klog"
)

func ProcessUnusedDisks(gcpChecker usage.Checker, k8sChecker usage.Checker) ([]string, error) {

	disksToBeRemoved := []string{}

	computeDisks, err := gcpChecker.ListResources()
	if err != nil {
		return nil, err
	}

	for _, diskInterface := range computeDisks {
		disk, ok := diskInterface.(*compute.Disk)
		if !ok {
			klog.Errorf("Error: expected item of type compute.Disk, got %T", diskInterface)
			continue // skip to the next disk if there's a type mismatch
		}
		isNotUsedByAnyNode, err := gcpChecker.IsResourceUnused(disk.Name)
		if err == nil && !isNotUsedByAnyNode {
			continue
		} else if err != nil {
			klog.Errorf("Error checking usage of disk %s: %v", disk.Name, err)
			return nil, err
		}

		isUnused, err := k8sChecker.IsResourceUnused(disk.Name)
		if err != nil {
			klog.Errorf("Error checking usage of disk %s: %v", disk.Name, err)
			continue // skip to the next disk if there's an error
		}

		if isUnused {
			disksToBeRemoved = append(disksToBeRemoved, disk.Name)
		}

	}
	return disksToBeRemoved, nil
}

func RemoveUnusedDisks(gcpDeleter *gcp.GCPDeleter, disks []string) {

	for _, disk := range disks {
		klog.Infof("Delete disk %s", disk)
		/* 		err := gcpDeleter.DeleteResource("disk", disk)
		   		if err != nil {
		   			klog.Errorf("Failed to delete disk %s: %v", disk, err)
		   			continue
		   		} */
	}

}
