// main.go
package main

import (
	"context"
	"errors"
	"os"

	"github.com/softonic/cloud-disks-cleaner/pkg/errorhandling"
	"github.com/softonic/cloud-disks-cleaner/pkg/gcp"
	"github.com/softonic/cloud-disks-cleaner/pkg/kubernetes"
	"github.com/softonic/cloud-disks-cleaner/pkg/usage"

	"github.com/softonic/cloud-disks-cleaner/internal/app"

	_ "k8s.io/apimachinery/pkg/api/resource"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"k8s.io/klog"
)

func init() {
	klog.InitFlags(nil)
}

func main() {

	// Load Configuration
	projectID, zone, err := loadConfiguration()
	if err != nil {
		errorhandling.HandleCriticalError(err) // Llama a la función de manejo de errores críticos.
	}

	// Initialize services
	gcpChecker, gcpDeleter, k8sChecker, err := initializeServices(projectID, zone)
	if err != nil {
		errorhandling.HandleCriticalError(err)
	}

	// Execute the app
	runApplication(gcpChecker, gcpDeleter, k8sChecker)

}

func loadConfiguration() (projectID string, zone string, err error) {
	projectID = os.Getenv("PROJECT_ID")
	zone = os.Getenv("ZONE")
	if projectID == "" || zone == "" {
		return "", "", errors.New("PROJECT_ID and ZONE environment variables are required")
	}
	return projectID, zone, nil
}

func initializeServices(projectID string, zone string) (usage.Checker, *gcp.GCPDeleter, usage.Checker, error) {
	// GCP init

	// client, err := gcp.NewClient()
	// if err != nil {
	// 	return nil, nil, nil, err
	// }

	ctx := context.Background()

	computeService, err := gcp.NewComputeService(ctx)
	if err != nil {
		klog.Fatalf("Failed to create GCP compute service: %v", err)
	}

	gcpChecker, err := gcp.NewGCPChecker(computeService, projectID, zone)
	if err != nil {
		klog.Errorf("Failed to create GCP checker: %v", err)
		return nil, nil, nil, err
	}

	gcpDeleter, err := gcp.NewGCPDeleter(computeService, projectID, zone)
	if err != nil {
		return nil, nil, nil, err
	}

	// k8s init

	k8sConfig := "/Users/santiago.nunezcacho/.kube/config"

	clientset, err := kubernetes.NewKubernetesService(false, k8sConfig)
	if err != nil {
		return nil, nil, nil, err
	}

	k8sChecker, err := kubernetes.NewK8sChecker(clientset)
	if err != nil {
		//klog.Errorf("Failed to create k8s checker: %v", err)
		return nil, nil, nil, err
	}

	return gcpChecker, gcpDeleter, k8sChecker, nil

}

func runApplication(gcpChecker usage.Checker, gcpDeleter *gcp.GCPDeleter, k8sChecker usage.Checker) {
	// app

	disks, err := app.ProcessUnusedDisks(gcpChecker, k8sChecker)
	if err != nil {
		klog.Errorf("Failed to process unused disks: %v", err)
	}

	app.RemoveUnusedDisks(gcpDeleter, disks)

}
