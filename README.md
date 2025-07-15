# RBAC Collector
A Kubernetes worker designed to run inside a cluster, collect RBAC (Role-Based Access Control) information, and log it in a structured format.

## Features
* Collects Roles, RoleBindings, ClusterRoles, and ClusterRoleBindings.
* Structured, human-readable logging with zerolog.
* Configurable collection interval.
* Automated build, test, and release pipeline using GitHub Actions.
* Container image vulnerability scanning with Trivy.
* Helm chart for easy deployment.

## CI/CD
The pipeline is configured to:

* On Pull Requests: Build the Docker image, run the Trivy vulnerability scan, and lint the Helm chart.
* On Push to main: Build and push the latest Docker image.
* On New v* Tag: Build and push a versioned Docker image, package the Helm chart, and create a new GitHub Release with the chart as an asset.