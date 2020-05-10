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

# Questions
1\. How would I deploy this kubernetes?
  - I would use a tool such as [helm](https://helm.sh) or [kapp](https://get-kapp.io) to provide a package that can be used to install 
    the controller as Deployment into a k8s cluster.
    
    
6\. In the context of your controller, what is the observed state and what is the desired state?
  - The observed state is the current state of a the world, in the case of this controller is the state of any given
   ConfigMap. The desired state is that if a ConfigMap contains the annotation  `x-kv8s.io/curl-me-that` with a valid key/url 
   assigned then there should be a data field on the ConfigMap with the key and associated data
   
    
7\. The content returned when curling URLs may be always different. How is it going to affect your controllers?
  - Since the output of curling the URL may be different the controller cannot make any assertions 
    around what the value of the data key should be, only that if the annotation is there and the key exists
    then it is in a desired state. This means a manual edit of the value of the data key would
    go unnoticed by the controller

