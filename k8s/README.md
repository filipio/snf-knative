# Knative on kubernetes cluster

## Setup

Setup is using knative-operator, based upon [knative-operator-deployment](https://knative.dev/docs/install/operator/knative-with-operators/)

1.  Install knative operator  
    `kubectl apply -f operator.yaml`  
    Check if operator is ready  
     `kubectl get deployment knative-operator`

2.  Install knative-serving  
    `kubectl apply -f serving.yaml`  
    Wait for all knative-serving deployments to be ready  
     `kubectl get deployment -n knative-serving --watch`  
    Check the status of Knative Serving Custom Resource  
    `kubectl get KnativeServing knative-serving -n knative-serving`

3.  Configure DNS using [sslip](https://sslip.io/)  
    `kubectl apply -f serving-default-domain.yaml`

4.  Install knative-eventing  
     `kubectl apply -f eventing.yaml`  
     Wait for all knative-eventing deployments to be ready  
     `kubectl get deployment -n knative-eventing --watch`  
     Check the status of Knative Eventing Custom Resource  
    `kubectl get KnativeEventing knative-eventing -n knative-eventing`

5.  Modify `image` value in `hello-world.yaml` if you want to use knative function created from scratch using instructions in [../src/hello-python/README.md](../src/hello-python/README.md). Image pattern for dockerhub is  
    `docker.io/<dockerhub-username>/hello-python`

6.  Create hello-world app function:  
    `kubectl apply -f hello-world.yaml`

## Access knative service

Execute  
`kubectl get ksvc`  
and use URL from output to access endpoints. Following methods are available:  
`GET <URL>` (with optional query parameters)  
`POST <URL>` (with optional json body)  
`GET <URL>/health/readiness`  
`GET <URL>/health/liveness`  
If you are using vscode, install `rest client` plugin and execute requests available in `test.http`. If not, check mentioned file to see how to execute requests with other client (curl, postman, ...)

## Clean up

1. `kubectl delete KnativeServing knative-serving -n knative-serving`

2. `kubectl delete KnativeEventing knative-eventing -n knative-eventing`

3. `kubectl delete -f operator.yaml`
