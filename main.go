// main.go
package main

import (
	"os"

	"github.com/softonic/cloud-disks-cleaner/pkg/gcp"
	"github.com/softonic/cloud-disks-cleaner/pkg/k8s"
	"github.com/softonic/cloud-disks-cleaner/pkg/usage"
	"google.golang.org/api/compute/v1"
	_ "k8s.io/apimachinery/pkg/api/resource"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/klog"
)

func init() {
	klog.InitFlags(nil)
}

func processUnusedDisks(gcpChecker usage.Checker, k8sChecker usage.Checker) ([]string, error) {

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

func removeUnusedDisks(gcpDeleter *gcp.GCPDeleter, disks []string) {

	for _, disk := range disks {
		klog.Infof("Delete disk %s", disk)
		/* 		err := gcpDeleter.DeleteResource("disk", disk)
		   		if err != nil {
		   			klog.Errorf("Failed to delete disk %s: %v", disk, err)
		   			continue
		   		} */
	}

}

func main() {

	projectID := os.Getenv("PROJECT_ID")
	zone := os.Getenv("ZONE")

	if projectID == "" || zone == "" {
		klog.Fatalf("PROJECT_ID and ZONE environment variables are required")
	}

	// GCP init

	client, err := gcp.NewClient()
	if err != nil {
		klog.Fatalf("Failed to create gcp clientset: %v", err)
		return
	}

	gcpChecker, err := gcp.NewGCPChecker(*client, projectID, zone)
	if err != nil {
		klog.Errorf("Failed to create GCP checker: %v", err)
	}

	gcpDeleter, err := gcp.NewGCPDeleter(*client, projectID, zone)
	if err != nil {
		klog.Errorf("Failed to instantiate gcp deleter: %v", err)
	}

	// k8s init

	clientset, err := k8s.NewClient(false)
	if err != nil {
		klog.Fatalf("Failed to create k8s clientset: %v", err)
		return
	}

	k8sChecker, err := k8s.NewK8sChecker(*clientset)
	if err != nil {
		klog.Errorf("Failed to create k8s checker: %v", err)
	}

	// app

	disks, err := processUnusedDisks(gcpChecker, k8sChecker)
	if err != nil {
		klog.Errorf("Failed to process unused disks: %v", err)
	}

	removeUnusedDisks(gcpDeleter, disks)

}
