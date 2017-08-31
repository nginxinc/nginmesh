
# Service Mesh with Istio and NGINX

## The Demo

Please see README.md in samples/bookinfo directory

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