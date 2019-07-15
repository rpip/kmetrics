# Kubernetes Platform Metrics

This service exposes some information on the current state of the kubernetes cluster. Written in Go

## Set up

First apply the Kubernetes service config

``` shell
$ kubectl apply -f services.yaml # minikube dashboard and paste the yaml in there
```

Now build the service or run directly

``` shell
$ make build # or go run .
```

## Configuration

Uses environment variables for config in line with the [12 Factor App](https://12factor.net).

Manage this with the `.env` file in the repo.


## Tests

``` shell
$ make test
```

## Endpoints

### 1. List information on all pods in the cluster

`/services` endpoint to the service that exposes all pods running in the cluster in namespace `default`:

Request: `GET /services`


Response:

```
[
  {
    "name": "blissful-goodall-deployment",
    "applicationGroup": "beta",
    "runningPodsCount": 1
  },
  {
    "name": "confident-cartwright-deployment",
    "applicationGroup": "beta",
    "runningPodsCount": 1
  },
  {
    "name": "happy-colden-deployment",
    "applicationGroup": "",
    "runningPodsCount": 1
  },
  {
    "name": "quirky-raman-deployment",
    "applicationGroup": "gamma",
    "runningPodsCount": 1
  },
  {
    "name": "stoic-sammet-deployment",
    "applicationGroup": "alpha",
    "runningPodsCount": 2
  }
]
```

### 2. Get information on a group of applications in the cluster

`/services/{group}` exposes the pods in the cluster in namespace `default` that are part of the same `applicationGroup`:

Request: `GET /services/{applicationGroup}`

Response:

```
[
  {
    "name": "blissful-goodall-deployment",
    "applicationGroup": "beta",
    "runningPodsCount": 1
  },
  {
    "name": "confident-cartwright-deployment",
    "applicationGroup": "beta",
    "runningPodsCount": 1
  }
]
```
