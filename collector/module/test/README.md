# Set up Mixer for testing

##  Set up Istio for developer

Follow instruction https://github.com/istio/istio/tree/master/devel

Build mixer.  This module works with 0.16 version, so be sure to use 0.16 tag.

`cd $(ISTIO)/mixer
git checkout 0.1.6
bazel clean
bazel build ...
`

This should build mixc and mixs binary in the go path

## Set up book

Set up BookInfo app as here:  https://istio.io/docs/samples/bookinfo.html

## Set up local forward to remote k8 instance

`kubectl get pods`

`kubectl port-forward istio-mixer-2450814972-f2m3w 9091:9091`

This will port forward local port 9091 to remote k8 mixer.

## Test mixer

Test that mixer can be reached

`./bazel-bin/cmd/client/mixc report  -a target.service=reviews.default.svc.cluster.local   --string_attributes request.headers=content-length:0`


 