# NGINX Streaming Architecture with Istio Service Mesh
This repository provides an implementation of a NGINX based service mesh (nginMesh).  nginMesh is compatible with Istio.  It leverages NGINX as a sidecar proxy. 

## What is Service Mesh and Istio?
Please check https://istio.io for a detailed explanation of the service mesh.  

## Production Status
The current version of nginMesh is designed to work with Istio release 0.6.0. It should not be used in production environments.  

## Architecture
The diagram below depicts how an NGINX sidecar proxy is implemented. The sidecar uses open source version of NGINX with first-party modules as well as third-party modules for tracing.
In 0.6.0 release, nginMesh leverages Kafka for delivery of mesh metrics. 
mTLS authentication between sidecars is planned in upcoming release.

![Alt text](/images/nginx_sidecar.png?raw=true "NGINX Sidecar")

To learn more about the sidecar implementation, see [this document](istio/agent).

## Quick Start
Below are instructions to quickly install and configure nginMesh.  Currently, only Kubernetes environment is supported.

### Prerequisites
Make sure you have a Kubernetes cluster with at least 1.9 or greater due to fact only automatic sidecar injection is supported and no alpha feature is enabled. Please see [Prerequisites](https://istio.io/docs/setup/kubernetes/quick-start.html) for setting up a kubernetes cluster.

### Installing Istio and nginMesh
nginMesh requires installation of Istio first.

1. Download and install Istio 0.6.0:
```
curl -L https://git.io/getLatestIstio | ISTIO_VERSION=0.6.0 sh -
```
2. Download nginMesh release 0.6.0:
```
curl -L https://github.com/nginmesh/nginmesh/releases/download/0.6.0/nginmesh-0.6.0.tar.gz | tar zx
```

3. Deploy Istio between sidecars:

```
kubectl create -f istio-0.6.0/install/kubernetes/istio.yaml
```


4. Ensure the following Kubernetes services are deployed: istio-pilot, istio-mixer, istio-ingress:
```
kubectl get svc  -n istio-system  
```
```
 NAME            CLUSTER-IP      EXTERNAL-IP       PORT(S)                       AGE
  istio-ingress   10.83.245.171   35.184.245.62     80:32730/TCP,443:30574/TCP    5h
  istio-pilot     10.83.251.173   <none>            8080/TCP,8081/TCP             5h
  istio-mixer     10.83.244.253   <none>            9091/TCP,9094/TCP,42422/TCP   5h
```

5. Ensure the following Kubernetes pods are up and running: istio-pilot-* , istio-mixer-* , istio-ingress-*  and istio-initializer-* :
```
kubectl get pods -n istio-system    
```
```
  istio-ca-3657790228-j21b9           1/1       Running   0          5h
  istio-ingress-1842462111-j3vcs      1/1       Running   0          5h
  istio-pilot-2275554717-93c43        1/1       Running   0          5h
  istio-mixer-2104784889-20rm8        2/2       Running   0          5h
```

6. Automatic sidecar:
To set up sidecar injection, please run following script which will install Istio webhook with nginMesh customization.
```
nginmesh-0.6.0/install/kubernetes/install-sidecar.sh
```

7. Verify that istio-injection label is not labeled for the default namespace :
```
kubectl get namespace -L istio-injection
```
```
NAME           STATUS        AGE       ISTIO-INJECTION
default        Active        1h        
istio-system   Active        1h        
kube-public    Active        1h        
kube-system    Active        1h
```

### Kafka deployment using Helm

1. Install the binary release of the Helm client depending on your OS.  Please follow [Setup guide](https://docs.helm.sh/using_helm/#quickstart) for detailed instructions.

2. Install Helm server(Tiller) in Kubernetes cluster:
```
helm init
```
3. Make sure Tiller is up and running:

```
kubectl get pods --namespace kube-system
```
```
NAME                                                 READY     STATUS    RESTARTS   AGE
tiller-deploy-f44659b6c-p48hf                        1/1       Running   0          51m
```

4. Run the following script to setup Kafka. It will be installed in 'kafka' namespace.  It is also possible to use existing kafka installation.

```
nginmesh-0.6.0/install/kafka/install.sh
```

Note: In GKE environment you may need to grant permission to default serviceaccount for cluster-wide access before install:

```
kubectl create clusterrolebinding add-on-cluster-admin --clusterrole=cluster-admin --serviceaccount=kube-system:default

kubectl create serviceaccount --namespace kube-system tiller
kubectl create clusterrolebinding tiller-cluster-rule --clusterrole=cluster-admin --serviceaccount=kube-system:tiller
kubectl patch deploy --namespace kube-system tiller-deploy -p '{"spec":{"template":{"spec":{"serviceAccount":"tiller"}}}}'
```

5. Wait for a while and make sure all kafka and zookeeper pods are up and runnning:

```
kubectl get pods -n kafka
```
```
NAME                   READY     STATUS    RESTARTS   AGE
my-kafka-kafka-0       1/1       Running   1          15m
my-kafka-kafka-1       1/1       Running   0          13m
my-kafka-kafka-2       1/1       Running   0          12m
my-kafka-zookeeper-0   1/1       Running   0          15m
my-kafka-zookeeper-1   1/1       Running   0          14m
my-kafka-zookeeper-2   1/1       Running   0          14m
testclient             1/1       Running   0          15m
```
6. Check the deployed release status:

```
helm list
```
```
NAME    	         REVISION	         UPDATED                       	STATUS  	    CHART      	  NAMESPACE
my-kafka	           1           Tue Mar 27 18:45:18 2018	          DEPLOYED  	 kafka-0.4.7  	   kafka
```

7. Set up  topic named "nginmesh" by running below script:

```
nginmesh-0.6.0/tools/kafka-add-topics.sh nginmesh
```
8. View created topic by running below script:

```
nginmesh-0.6.0/tools/kafka-list-topics.sh nginmesh
```
```
nginmesh
```

### Deploy a Sample Application
In this section we deploy the Bookinfo application, which is taken from the Istio samples. Please see [Bookinfo](https://istio.io/docs/guides/bookinfo.html)  for more details.

1. Label the default namespace with istio-injection=enabled:

```
kubectl label namespace default istio-injection=enabled
```

2. Deploy the application:

```
kubectl apply -f  nginmesh-0.6.0/samples/bookinfo/kube/bookinfo.yaml
```

3. Confirm that all application services are deployed: productpage, details, reviews, ratings:

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

4. Confirm that all application pods are running --details-v1-* , productpage-v1-* , ratings-v1-* , reviews-v1-* , reviews-v2-* and reviews-v3-* :
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

5. Get the public IP of the Istio Ingress controller. If the cluster is running in an environment that supports external load balancers:

```
kubectl get svc -n istio-system | grep -E 'EXTERNAL-IP|istio-ingress'
```

OR

```
kubectl get ingress -o wide       
```

6. Open the Bookinfo application in a browser using the following link:
```
http://<Public-IP-of-the-Ingress-Controller>/productpage
```
### Uninstalling the Application
1. To uninstall application, run:

```
./nginmesh-0.6.0/samples/bookinfo/kube/cleanup.sh
```


### Uninstalling Istio
1. To uninstall the Istio core components:

```
kubectl delete -f istio-0.6.0/install/kubernetes/istio.yaml
```


2. To uninstall the initializer, run:

```
nginmesh-0.6.0/install/kubernetes/delete-sidecar.sh
```

### Uninstalling Kafka

1. Uninstall Kafka:

```
nginmesh-0.6.0/install/kafka/uninstall.sh
``` 

2. Delete Tiller deployment:

```
kubectl delete deployment tiller-deploy --namespace kube-system
```

 


## Limitations
nginMesh has the following limitations:
* TCP and gRPC traffic is not supported.
* Quota Check is not supported.
* Only Kubernetes is supported.

All sidecar-related limitations and supported traffic management rules are described [here](istio/agent).
