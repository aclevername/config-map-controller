# config-map-controller

A controller that looks for an annotation `x-kv8s.io/curl-me-that: mydata=data.example.com` and will append a data field `mydata` with the contents of curling `data.example.com`

## Tutorial
### Start controller
1.`go build main.go && ./main --kubeconfig $KUBECONFIG` where `$KUBECONFIG` is an environment variable pointing to your kubeconfig
### Create configmap
In a seperate terminal run
1. `kubectl create -f fixtures/config-map-valid.yml` 
1. `kubectl get configmap valid -o yaml` and you should see a `mydata` key added with a joke as its value

If something goes wrong, for example an invalid annotation `x-kv8s.io/curl-me-that: mydata=totally not a url` you can
view the error by running `kubectl describe configmap invalid`. For example:
1. `kubectl create -f fixtures/config-map-invalid.yml` 
1. `kubectl get configmap invalid -o yaml` and you should see it hasn't changed
1. `kubectl describe configmap invalid` to see the error


# Requirements
- go `1.13.8` to build and run
- Access to a kubernetes cluster `v1.17.0` (tested with kind `v0.7.0`) to run `acceptance-tests` against
- [counterfeiter](https://github.com/maxbrunsfeld/counterfeiter) to generate fakes (only required if editing codebase)
- [ginkgo](https://onsi.github.io/ginkgo/) installed to run the tests

# Running the tests

run `make test-unit` to only run the unit tests

run `make test-acceptance` to only run the unit tests. Requires `KUBECONFIG` environment 
variable exported pointing at a valid kubeconfig

run `make test` to run all tests. Requires `KUBECONFIG` environment 
variable exported pointing at a valid kubeconfig
 
#  

# Acknowledgments
- [Bitnami introduction to controllers](https://engineering.bitnami.com/articles/a-deep-dive-into-kubernetes-controllers.html) for explaining the basics of controllers
- [client-go examples](https://github.com/kubernetes/client-go/blob/master/examples/) for how to use the `client-go`
- [Community guide to writing controllers](https://github.com/kubernetes/community/blob/master/contributors/devel/sig-api-machinery/controllers.md) for how to setup the workqueue to watch for resource changes

