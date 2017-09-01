
# Service Mesh with Istio and NGINX

This repo provides an implementation of sidecar proxy for Istio using NGINX open source version.

## What is Service Mesh and Istio?

Please see https://istio.io for detail explanation of service mesh provided by Istio.  
Combination of Nginx and Istio provides best service mesh for deploying micro-services.

## Production Status

This is the preview version of the Nginx proxy design to run with 0.16 release of Istio.
Due to arch differences between 0.16 and 0.2 of Istio, some feature such as routing rules and auth are not implemented in the preview; they are under development and will be available in the 0.2 release.


## Architecture

Please see diagram below to see how Nginx Sidecar Proxy is implemented as of 0.16 version.
The sidecar run NGINX open source version with custom module to interface to Istio Mixer.

![Alt text](/images/nginx_sidecar.png?raw=true "Nginx Sidecar")

## The Demo

Please see README.md in samples/bookinfo directory.  

The sample app is copied from Istio project without modification.  
App deployment is done using custom sidecar inject tool.  Please see below for building tool for your OS.

You can use following kubectl snippet to go inside NGNX proxy container.  This will enter 'productpage' service proxy.

```bash
kubectl exec -it  $(kubectl get pod -l app=productpage  -o jsonpath='{.items[0].metadata.name}') -c proxy  /bin/bash
```

## Build the nginx-inject

If you would like to build the nginx-inject tool used in the demo for injecting containers into Kubernetes yaml files, make sure you have Docker engine and make installed on your system: 

Clone the repo and run one of the following command from this folder to build the tool:

If you use macOS:
```
$ make darwin
```

If you use Linux:
```
$ make linux
```

The tool will be available in the `build` folder.