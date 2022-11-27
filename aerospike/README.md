# Aerospike

## Prerequisites

- `go` version v1.19+
- `docker` version 20.10.21+
- `kubectl` version v1.21+
- [`gcloud` CLI](https://cloud.google.com/sdk/docs/install), version 407.0+
- `minikube` version v1.27.1+
- Access to a Kubernetes v1.23+ cluster with the following permissions: read/create/delete pods and read/create namespaces.

## Kubernetes Cluster Setup

A Kubernetes cluster can be created with `minikube` locally or created remotely with Google Cloud GKE.

#### Minikube
To create a `minikube` cluster locally, run:
```
$ minikube start
```

#### Google Cloud GKE
To create a Kubernetes cluster remotely on GKE run the following commands from the root of this repo:

Prerequisites:
- Create a GCP account
- Create a GCP project
- Enable the Kubernetes Engine API in GCP

Authenticate to Google Cloud from the command line:
```
$ make auth
```

Create a GKE Kubernetes cluster:
```
$ make build
```

Configure `.kube/config` to connect to the Kubernetes cluster from a local terminal:
```
$ make connect
```

## Run the program

From the root of this repo, build the program executable:
```
$ go build
```

Usage:
```
$ ./aerospike --help
Usage of ./aerospike:
  -delete
    	Executes delete operations. Default false.
  -kubeconfig string
    	absolute path to the kubeconfig file (default "$HOME/.kube/config")
```

Executing the program:
```
$ ./aerospike
```

Executing the program will perform all of the following operations one after another:
- Lists all namespaces, writes each namespace name to stdout
- Creates a new namespace named `app`
- Creates a hello world pod in the new namespace
- Lists all pods with the label "k8s-app=kube-dns", writes each pod name and namespace to stdout

Once the hello world pod is created, it can be accessed like so:
```
$ kubectl port-forward pod/hello-world 8080:8080 -n app

// in another shell make a request against the hello world app
$ curl localhost:8080
```

Clean up after the program:
```
$ ./aerospike -delete
```

Executing the delete command deletes the hello world pod.


## Kubernetes Cluster Teardown

#### Minikube
To delete the `minikube` cluster, run:
```
minikube delete --all
```

#### Google Cloud GKE
To delete the Kubernetes Cluster running on GKE, run:
```
$ make delete
```
