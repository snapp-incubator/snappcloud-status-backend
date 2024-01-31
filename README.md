# Helm Charts for SnappCloud Status page Services

## Prerequisites

Before you begin, ensure you have Helm installed.

## Usage

To use the Helm charts for each project, follow these steps:

1. Add the repository:

```shell
helm repo add snappcloud-status-page https://snapp-incubator.github.io/snappcloud-status-backend/
```

If you had already added this repo earlier, run `helm repo update` to retrieve the latest versions of the packages.

2. You can access to helm charts available within the repository via:

```shell
helm search repo snappcloud-status-page
```

3. Install each Helm chart:

```shell
helm install my-release snappcloud-status-page/<chart-name>
```

4. To uninstall each helm chart:

```shell
helm delete my-release
```
