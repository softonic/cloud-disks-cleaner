package notifiers

import (
	"fmt"

	"github.com/ashwanthkumar/slack-go-webhook"
	"github.com/softonic/cloud-disks-cleaner/pkg/common"
	"k8s.io/klog"
)

type SlackClient struct {
	WebhookUrl string
	Channel    string
	Username   string
	IconEmoji  string
}

func NewSlackClient() *SlackClient {
	return &SlackClient{
		WebhookUrl: "https://hooks.slack.com/services/...",

		Channel:   "#k8s-clean-disks-GCP-unused",
		Username:  "k8s-clean-disks-GCP-unused",
		IconEmoji: ":ghost:",
	}
}

func (c SlackClient) SendNotification(disk common.DiskInfo) {

	attachment := slack.Attachment{}
	attachment.AddField(slack.Field{Title: "Disks to be removed", Value: fmt.Sprintf("%v", disk.DiskName)})
	attachment.AddField(slack.Field{Title: "pvcName", Value: disk.PVCName})
	attachment.AddField(slack.Field{Title: "Namespace", Value: disk.Namespace})
	attachment.AddAction(slack.Action{Type: "button", Text: "Check in the console", Url: "https://console.cloud.google.com/compute/disks?project=kubertonic", Style: "primary"})

	payload := slack.Payload{
		Text:        "Disk Removed in GCP",
		Username:    c.Username,
		Channel:     c.Channel,
		IconEmoji:   c.IconEmoji,
		Attachments: []slack.Attachment{attachment},
	}

	err := slack.Send(c.WebhookUrl, "", payload)
	if err != nil {
		klog.Error("\nError: ", err)
	}
}
