# Platform Developer Test

This service exposes some information on the current state of the kubenertes cluster. Written in Go

## Set up

First apply the Kubernetes service config

``` shell
$ kubectl apply -f services.yaml # minikube dashboard and paste the yaml in there
```

Now build the service or run directly

``` shell
$ make build # or go run .
```

## Tests

``` shell
$ make test
```

## Endpoints

### 1. List information on all pods in the cluster

`/services` endpoint to the service that exposes all pods running in the cluster in namespace `default`:

```
GET `/services`
[
  {
    "name": "first",
    "applicationGroup": "alpha",
    "runningPodsCount": 2
  },
  {
    "name": "second",
    "applicationGroup": "beta",
    "runningPodsCount": 1
  },
  ...
]
```

### 2. Get information on a group of applications in the cluster

`/services/{group}` exposes the pods in the cluster in namespace `default` that are part of the same `applicationGroup`:

```
GET `/services/{applicationGroup}`
[
  {
    "name": "foobar",
    "applicationGroup": "<applicationGroup>",
    "runningPodsCount": 1
  },
  ...
]
```
