package notifiers

import "github.com/softonic/k8s-clean-disks-GCP-unused/pkg/common"

type ClientInterface interface {
	SendNotification(message common.DiskInfo)
}
