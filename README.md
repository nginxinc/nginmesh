# Service Mesh with Istio and NGINX
This repository provides an implementation of a sidecar proxy based on NGINX for Istio.

## What is Service Mesh and Istio?
Please check https://istio.io for a detailed explanation of the service mesh provided by Istio.

## Production Status
The current version of nginmesh works with Istio release 0.3.0. It is not suited for production deployments.

## Architecture
The diagram below depicts how an NGINX sidecar proxy is implemented. The sidecar runs NGINX with custom modules to interface with Istio Mixer, as well as with third-party modules for tracing.

![Alt text](/images/nginx_sidecar.png?raw=true "Nginx Sidecar")

To learn more about the sidecar implementation, see [this document](istio/agent).

## Quick Start
Below are instructions to setup the Istio service mesh in a Kubernetes cluster using NGINX as a sidecar.

### Prerequisites
Make sure you have a Kubernetes cluster with alpha feature enabled. Please see [Prerequisites](https://istio.io/docs/setup/kubernetes/quick-start.html#prerequisites) for setting up a cluster.

### Installing Istio and nginmesh
Below are instructions for installing Istio with NGINX as a sidecar:
1. Download Istio release 0.3.0:
```
curl -L https://git.io/getLatestIstio | ISTIO_VERSION=0.3.0 sh -
```
2. Download nginmesh release 0.3.0:
```
curl -L https://github.com/nginmesh/nginmesh/releases/download/0.3.0/nginmesh-0.3.0.tar.gz | tar zx
```

3. Deploy Istio either with or without enabled mutual TLS (mTLS) authentication between sidecars:

a) Install Istio without enabling mTLS:
```
kubectl create -f istio-0.3.0/install/kubernetes/istio.yaml
```
b) Install Istio with mTLS:
```
kubectl create -f istio-0.3.0/install/kubernetes/istio-auth.yaml
```
4. Deploy an initializer for automatic sidecar injection:
```
kubectl apply -f nginmesh-0.3.0/install/kubernetes/istio-initializer.yaml
```

5. Ensure the following Kubernetes services are deployed: istio-pilot, istio-mixer, istio-ingress, istio-egress:
```
kubectl get svc  -n istio-system  
```
```
 NAME            CLUSTER-IP      EXTERNAL-IP       PORT(S)                       AGE
  istio-egress    10.83.247.89    <none>            80/TCP                        5h
  istio-ingress   10.83.245.171   35.184.245.62     80:32730/TCP,443:30574/TCP    5h
  istio-pilot     10.83.251.173   <none>            8080/TCP,8081/TCP             5h
  istio-mixer     10.83.244.253   <none>            9091/TCP,9094/TCP,42422/TCP   5h
```

6. Ensure the following Kubernetes pods are up and running: istio-pilot-* , istio-mixer-* , istio-ingress-* , istio-egress-* and istio-initializer-* :
```
kubectl get pods -n istio-system    
```
```
  istio-ca-3657790228-j21b9           1/1       Running   0          5h
  istio-egress-1684034556-fhw89       1/1       Running   0          5h
  istio-ingress-1842462111-j3vcs      1/1       Running   0          5h
  istio-initializer-184129454-zdgf5   1/1       Running   0          5h
  istio-pilot-2275554717-93c43        1/1       Running   0          5h
  istio-mixer-2104784889-20rm8        2/2       Running   0          5h
```
### Deploy a Sample Application
In this section we deploy the Bookinfo application, which is taken from the Istio samples. Please see [Bookinfo](https://istio.io/docs/guides/bookinfo.html)  for more details.

1. Deploy the application:
```
kubectl apply -f nginmesh-0.3.0/samples/kubernetes/bookinfo.yaml
```

2. Confirm that all application services are deployed: productpage, details, reviews, ratings.
```
kubectl get services
```
```
NAME                       CLUSTER-IP   EXTERNAL-IP   PORT(S)              AGE
details                    10.0.0.31    <none>        9080/TCP             6m
kubernetes                 10.0.0.1     <none>        443/TCP              7d
productpage                10.0.0.120   <none>        9080/TCP             6m
ratings                    10.0.0.15    <none>        9080/TCP             6m
reviews                    10.0.0.170   <none>        9080/TCP             6m
```

3. Confirm that all application pods are running --details-v1-* , productpage-v1-* , ratings-v1-* , ratings-v1-* , reviews-v1-* , reviews-v2-* and reviews-v3-* :
```
kubectl get pods
```
```
NAME                                        READY     STATUS    RESTARTS   AGE
details-v1-1520924117-48z17                 2/2       Running   0          6m
productpage-v1-560495357-jk1lz              2/2       Running   0          6m
ratings-v1-734492171-rnr5l                  2/2       Running   0          6m
reviews-v1-874083890-f0qf0                  2/2       Running   0          6m
reviews-v2-1343845940-b34q5                 2/2       Running   0          6m
reviews-v3-1813607990-8ch52                 2/2       Running   0          6m
```

4. Get the public IP of the Istio Ingress controller. If the cluster is running in an environment that supports external load balancers, run:
```
kubectl get svc -n istio-system | grep -E 'EXTERNAL-IP|istio-ingress'
```
OR
```
kubectl get ingress -o wide       
```

5. Open the Bookinfo application in a browser using the following link:
```
http://<Public-IP-of-the-Ingress-Controller>/productpage
```
### Uninstalling the Application
1. To uninstall application, run:
```
./nginmesh-0.3.0/samples/kubernetes/cleanup.sh 
```


### Uninstalling Istio
1. To uninstall the Istio core components:

a) If mTLS is disabled:
```
kubectl delete -f istio-0.3.0/install/kubernetes/istio.yaml
```

OR:

b) If mTLS is enabled:
```
kubectl delete -f istio-0.3.0/install/kubernetes/istio-auth.yaml
```

2. To uninstall the initializer, run:
```
kubectl delete -f nginmesh-0.3.0/install/kubernetes/istio-initializer.yaml
```

### Additional Examples

* **In-Depth Telemetry** [This example](https://istio.io/docs/guides/telemetry.html) demonstrates how to obtain uniform metrics, logs, traces across different applications. To run the example, you must install the telemetry services. 


## Limitations
nginmesh has the following limitations:
* TCP and gRCP traffic is not supported.
* Quota Check is not supported.
* Only Kubernetes is supported.

All sidecar-related limitations as well as supported traffic management rules are described [here](istio/agent).
