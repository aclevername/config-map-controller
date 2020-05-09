# config-map-controller

A controller that looks for an annotation `x-k8s.io/curl-me-that: mydata=data.example.com` and will append a data field `mydata` with the contents of curling `data.example.com`

## Example
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: config
  namespace: default
  x-k8s.io/curl-me-that: mydata=http://curl-a-joke.herokuapp.com
data:
  hello: "world"
```
will eventually turn into 
```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: config
  namespace: default
  x-k8s.io/curl-me-that: mydata=http://curl-a-joke.herokuapp.com
data:
  hello: "world"
  mydata: "Why was the broom late? It over swept!"
```

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
- [client-go workqueue examples](https://github.com/kubernetes/client-go/blob/master/examples/workqueue/main.go) for how to setup the boilerplate for a controller

