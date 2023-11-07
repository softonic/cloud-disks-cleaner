# cloud-disks-cleaner

[![Version Widget]][Version] [![License Widget]][License] [![GoReportCard Widget]][GoReportCard] [![DockerHub Widget]][DockerHub]

[Version]: https://github.com/softonic/kube-gcp-disks-roomba/releases
[Version Widget]: https://img.shields.io/github/release/softonic/kube-gcp-disks-roomba.svg?maxAge=60
[License]: http://www.apache.org/licenses/LICENSE-2.0.txt
[License Widget]: https://img.shields.io/badge/license-APACHE2-1eb0fc.svg
[GoReportCard]: https://goreportcard.com/report/softonic/kube-gcp-disks-roomba
[GoReportCard Widget]: https://goreportcard.com/badge/softonic/kube-gcp-disks-roomba
[DockerHub]: https://hub.docker.com/r/softonic/kube-gcp-disks-roomba
[DockerHub Widget]: https://img.shields.io/docker/pulls/softonic/kube-gcp-disks-roomba.svg


Script that runs as a cronjob resource in kubernetes in GKE environments.
Removes disks from GCP that are not in use, checking first if the storage class is the default (standard).
Understanding that standard storage class has reclaimPolicy Delete.

##### Install

```
GO111MODULE=on
go build .
```

