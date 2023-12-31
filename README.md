# cloud-disks-cleaner

[![Version Widget]][Version] [![License Widget]][License] [![GoReportCard Widget]][GoReportCard] [![DockerHub Widget]][DockerHub]

[Version]: https://github.com/softonic/cloud-disks-cleaner/releases
[Version Widget]: https://img.shields.io/github/release/softonic/cloud-disks-cleaner.svg?maxAge=60
[License]: http://www.apache.org/licenses/LICENSE-2.0.txt
[License Widget]: https://img.shields.io/badge/license-APACHE2-1eb0fc.svg
[GoReportCard]: https://goreportcard.com/report/softonic/cloud-disks-cleaner
[GoReportCard Widget]: https://goreportcard.com/badge/softonic/cloud-disks-cleaner
[DockerHub]: https://hub.docker.com/r/softonic/cloud-disks-cleaner
[DockerHub Widget]: https://img.shields.io/docker/pulls/softonic/cloud-disks-cleaner.svg






Cloud disks cleaner

Brief description of what the project does.
Description

This project is designed to run within a Kubernetes cluster, preferably on Google Kubernetes Engine (GKE). Its primary purpose is to identify and remove unused Google Cloud Platform (GCP) disks. This includes disks that are not attached to any virtual machines and do not have an associated Persistent Volume (PV) or Persistent Volume Claim (PVC) in Kubernetes. The project helps to keep the GCP environment clean by removing unused resources, which can contribute to cost reduction and more efficient resource management.

Script that runs as a cronjob resource in kubernetes in GKE environments.
Removes disks from GCP that are not in use. 
If PV does not exists in kubernetes, we can remove the disk in GCP.
If PV does exists, but pvc is not bound to the pv, we can remove the disk in GCP.

How It Works

The project follows these steps:

    Load Configuration: Reads the necessary configuration to run the application.
    Initialize Services: Sets up the necessary services to interact with GCP and Kubernetes.
    Run the Application: Performs the following operations:
        Lists disks in GCP and checks if they are being used.
        Checks if these disks have associated PVs or PVCs in Kubernetes.
        Deletes the disks that are not in use in GCP.

Installation

(Provide here the steps to install the project, e.g., cloning the repository, installing dependencies, etc.)
Usage


```
GO111MODULE=on
go build .
```



Describe how to run or use the application, including any necessary configuration parameters.


# Example of how to run the application
# Command to run the application (if applicable)


Contribution

Instructions for those who wish to contribute to the project.
License

Include details about the license under which the project is distributed.