# SnappCloud Status Backend

This project aims to provide the backend for `snappcloud-status-frontend` application to observer applications healthness in SnappCloud.

## Deployment

We normally deploy our production applications via `argoCD` and staging applications via `helm chart` as described below.

### Installation

``` bash
oc project your-desired-namespace # set your namespace

cmd=<cmd> # server
override=<override> # production, staging

./deployments/tearup.sh $cmd $override
```

### Uninstallation

``` bash
oc project your-desired-namespace # set your namespace

cmd=<cmd> # server

./deployments/teardown.sh $cmd
```
