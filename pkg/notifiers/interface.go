package notifiers

import "github.com/softonic/cloud-disks-cleaner/pkg/common"

type ClientInterface interface {
	SendNotification(message common.DiskInfo)
}
