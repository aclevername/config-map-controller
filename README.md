# config-map-controller

A controller that looks for an annotation `x-k8s.io/curl-me-that: mydata=data.example.com` and will append a data field `mydata` with the contents of curling `data.example.com`


# Test
run `make test` to run tests
