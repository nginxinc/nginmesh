
# Service Mesh with Istio and NGINX

This repo provides an implementation of sidecar proxy for Istio using NGINX open source version.

## What is Service Mesh and Istio?

Please see https://istio.io for detail explanation of service mesh provided by Istio.  
Combination of Nginx and Istio provides best service mesh for deploying micro-services.

## Production Status

This version of nginmesh work with with 0.2.12 release of Istio.
Please see below for Istio features we support.  Nginmesh is not production ready yet.  


<TBD>

## Architecture

Please see diagram below to see how Nginx Sidecar Proxy is implemented as of 0.2.12 version.
The sidecar run NGINX open source version with custom module to interface to Istio Mixer.

![Alt text](/images/nginx_sidecar.png?raw=true "Nginx Sidecar")

## Quick start
Below are instructions to setup the Istio service mesh in a Kubernetes cluster using NGINX as a sidecar.
 

### Prerequisities 

Make sure you have  kubernetes cluster with alpha feature enabled. Please, refer to [Prerequisites](https://istio.io/docs/setup/kubernetes/quick-start.html#prerequisites) in Istio project for setting up cluster for your environment.

### Installation Istio and Nginmesh
Below are instructions for installing Istio with NGINX as a sidecar in application pods.

1.  Download Istio release 0.2.12, by running below command:

```
curl -L https://git.io/getLatestIstio | ISTIO_VERSION=0.2.12 sh -
```

2. Download Nginmesh release 0.2.12, by running below command:

```
curl -L https://github.com/nginmesh/nginmesh/releases/download/0.2.12-RC2/nginmesh-0.2.12.tar.gz | tar zx
```

3. Create Istio deployment without authentication:
```
kubectl create -f istio-0.2.12/install/kubernetes/istio.yaml
```

4. Deploy automatic sidecar injection initializer:
```
kubectl apply -f nginmesh-0.2.12/install/kubernetes/istio-initializer.yaml
```

5. Ensure the following Kubernetes services are deployed: istio-pilot, istio-mixer, istio-ingress, istio-egress:

```
kubectl get svc  -n istio-system  
```
which produces the following output:
```
  NAME            CLUSTER-IP      EXTERNAL-IP       PORT(S)                       AGE
  istio-egress    10.83.247.89    <none>            80/TCP                        5h
  istio-ingress   10.83.245.171   35.184.245.62     80:32730/TCP,443:30574/TCP    5h
  istio-pilot     10.83.251.173   <none>            8080/TCP,8081/TCP             5h
  istio-mixer     10.83.244.253   <none>            9091/TCP,9094/TCP,42422/TCP   5h
```


6. Ensure the corresponding Kubernetes pods are up and running: istio-pilot-* , istio-mixer-* , istio-ingress-* , istio-egress-* and istio-initializer-* :
```
kubectl get pods -n istio-system    
```
which produces the following output:
```
  istio-ca-3657790228-j21b9           1/1       Running   0          5h
  istio-egress-1684034556-fhw89       1/1       Running   0          5h
  istio-ingress-1842462111-j3vcs      1/1       Running   0          5h
  istio-initializer-184129454-zdgf5   1/1       Running   0          5h
  istio-pilot-2275554717-93c43        1/1       Running   0          5h
  istio-mixer-2104784889-20rm8        2/2       Running   0          5h
```

### Deploy Application
The sample app is copied from Istio project without modification. Please, refer to [Bookinfo](https://istio.io/docs/guides/bookinfo.html) for more details.  

Note: We only support deployment using Kubernetes initializer. 

1. Deploy the application containers:

```
kubectl apply -f nginmesh-0.2.12/samples/kubernetes/bookinfo.yaml
```

2. Confirm all application services are correctly defined and running: productpage, details, reviews,ratings.
```
kubectl get services
```
which produces the following output:
```
NAME                       CLUSTER-IP   EXTERNAL-IP   PORT(S)              AGE
details                    10.0.0.31    <none>        9080/TCP             6m
kubernetes                 10.0.0.1     <none>        443/TCP              7d
productpage                10.0.0.120   <none>        9080/TCP             6m
ratings                    10.0.0.15    <none>        9080/TCP             6m
reviews                    10.0.0.170   <none>        9080/TCP             6m
```
3. Confirm all application pods are correctly defined and running: details-v1-* , productpage-v1-* , ratings-v1-* , ratings-v1-* , reviews-v1-* , reviews-v2-* and reviews-v3-* :

```
kubectl get pods
```
which produces the following output:
```
NAME                                        READY     STATUS    RESTARTS   AGE
details-v1-1520924117-48z17                 2/2       Running   0          6m
productpage-v1-560495357-jk1lz              2/2       Running   0          6m
ratings-v1-734492171-rnr5l                  2/2       Running   0          6m
reviews-v1-874083890-f0qf0                  2/2       Running   0          6m
reviews-v2-1343845940-b34q5                 2/2       Running   0          6m
reviews-v3-1813607990-8ch52                 2/2       Running   0          6m
```
4. If cluster is running in an environment that supports external load balancers, the IP address of ingress can be obtained by the following command:
```
kubectl get svc -n istio-system | grep -E 'EXTERNAL-IP|istio-ingress'
```
OR
```
kubectl get ingress -o wide       
```
5. Set Variable to Ingress address obtained in Step 4:
```
export GATEWAY_URL=104.196.5.186:80
```
6. To confirm that the BookInfo application is up and running:

a) Run the following curl command and check received response code is 200:

```
curl -o /dev/null -s -w "%{http_code}\n" http://${GATEWAY_URL}/productpage
```

OR:

b) Open in browser Bookinfo application and make sure successfully running:
```
http://${GATEWAY_URL}/productpage
```
### Cleanup Application

1.  Uninstall application, run the following shell script:
```
./samples/kubernetes/cleanup.sh 
```

2.  Make sure Application pods and services lists are empty:
```
kubectl get pods
kubectl get svc
```

### Cleanup Istio

1. Uninstall Istio from Kubernetes environment, run the following commands:

a) If no authentication enabled:
```
kubectl delete -f istio-0.2.12/install/kubernetes/istio.yaml
```

OR:

b) If authentication enabled:
```
kubectl delete -f istio-0.2.12/install/kubernetes/istio-auth.yaml
```
2. Uninstall Initializor, run the following commands:
```
kubectl delete -f nginmesh-0.2.12/install/kubernetes/istio-initializer.yaml
```

3. Make sure Istio pods and services lists are empty:
```
kubectl get pods -n istio-system
kubectl get svc  -n istio-system 
```

### Optional: 

[In-Depth Telemetry](https://istio.io/docs/guides/telemetry.html) This sample demonstrates how to obtain uniform metrics, logs, traces across different services using NGiNX sidecar. Additionally, for quick install of telemetry services, please check this [link](https://github.com/nginmesh/nginmesh/blob/release-doc-0.2.12/istio/release/install/kubernetes/README.md).

[Intelligent Routing](https://istio.io/docs/guides/intelligent-routing.html) Refer to Michael README.md ? Difference in delay with Istio.


