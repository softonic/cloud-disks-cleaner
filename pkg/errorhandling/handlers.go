package errorhandling

import "k8s.io/klog"

func HandleCriticalError(err error) {
	if err != nil {
		klog.Fatalf("Critical error encountered: %v", err)
	}
}

func HandleNonCriticalError(err error, resourceName string) bool {
	if err != nil {
		klog.Errorf("Non-critical error encountered for resource %s: %v", resourceName, err)
		return true
	}
	return false
}
