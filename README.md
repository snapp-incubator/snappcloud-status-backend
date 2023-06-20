# SnappCloud Status Backend

This project aims to provide the backend for `snappcloud-status-frontend` application to observe applications healthness in SnappCloud.

## Run

1. port-forward `thanos-query-frontend-http` service

    ```bash
    oc project openshift-monitoring
    kubectl port-forward service/thanos-query-frontend-http 9090:9090
    ```

2. run server command

    ```bash
    go run main.go server
    ```

## Deployment

We normally deploy our production applications via `argoCD` and staging applications via `helm chart` as described below.

### Installation

``` bash
oc project your-desired-namespace # set your namespace

cmd=<cmd> # server

./deployments/tearup.sh $cmd
```

### Uninstallation

``` bash
oc project your-desired-namespace # set your namespace

cmd=<cmd> # server

./deployments/teardown.sh $cmd
```
