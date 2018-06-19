# Ingenious
_Ingenious_ is a photo-sharing demo app created by NGINX to show the Microservices approach to application development. The app is designed to allow the user to login to a personalized account, and then store, view and delete their own pictures. It also includes a blog in which users can view the latest news and updates within the application.

The _Ingenious_ application is built with microservices and utilizes their inherent benefits to generate a robust, stable, independent, polyglot environment.


![Microservice Reference Architecture diagram of services](/images/diagram-microservices-reference-architecture-850x600.png)

The _Ingenious_ application employs seven different services in order to create its functionality.

Pages is the foundational service built in PHP upon which the other services provide functionality. Pages makes calls directly to User Manager, Album Manager, Content Service, and Uploader.

User Manager is built completely using Python and backed by DynamoDB. It's use is to store and modify user information, allowing the app a login system. Login is done with Google and Facebook through OAuth, but also includes a system for local login when testing the system.

Album Manager is built using Ruby and backed by MySQL, and allows the user to upload albums of multiple images at once. Album Manager makes calls to the Uploader service and therefore the Resizer to upload and modify images specified by the user.

The Uploader service is built using Javascript and is used to upload images to an S3 bucket. Uploader then makes calls to the Resizer service with the previously generated id of the image within S3, and Resizer then makes copies of the image with size "Large", "Medium", and "Thumbnail".

Content Service is built in Go and backed by RethinkDB. The Content Service provides, retrieves, and displays content for the NGINX _Ingenious_ application

Auth Proxy is a Python app that utilizes Redis' capabilities as a caching server. Making direct connections to both Pages and User Manager, Auth Proxy is used to validate the user's identity. It also serves as the gateway into the application, acting as the only public-facing service within the application.


## Quick Start
Below are instructions to setup the Istio service mesh in a Kubernetes cluster using NGINX as a sidecar.

### Prerequisites
Make sure you have a Kubernetes cluster with alpha feature enabled. Please see [Prerequisites](https://istio.io/docs/setup/kubernetes/quick-start.html#prerequisites) for setting up a cluster and for installation of Istio and nginmesh latest version refer to [nginmesh](https://github.com/nginxinc/nginmesh#installing-istio-and-nginmesh) repo.




### Deploy an Ingenious Application
In this section we deploy the Ingenious application. Please see [Ingenious](https://github.com/nginxinc/mra-ingenious)  for more details.

1. Deploy the application:
```
kubectl apply -f ingenious.yaml
```

2. Get the public IP of the Istio Ingress controller and Fake-s3 service. If the cluster is running in an environment that supports external load balancers, run:
```
kubectl get svc -n istio-system | grep -E 'EXTERNAL-IP|istio-ingress'
```
```
kubectl get svc -n default | grep -E 'EXTERNAL-IP|fake-s3'
```

3. Modify your _hosts_ file in Client machine to include the public IP of fake-s3 service. It should look like:
    ```
    104.196.56.175   fake-s3 http://fake-s3
    ```

4. Confirm that all application services are deployed: album-manager, auth-proxy, content-service, fake-s3, pages, photoresizer, photouploader, user-manager:
```
kubectl get services
```
```
NAME              TYPE           CLUSTER-IP     EXTERNAL-IP      PORT(S)        AGE
album-manager     ClusterIP      10.3.241.120   <none>           80/TCP         13h
auth-proxy        LoadBalancer   10.3.242.142   <none>           80:32700/TCP   13h
content-service   ClusterIP      10.3.255.176   <none>           80/TCP         13h
fake-s3           LoadBalancer   10.3.246.77    104.196.56.175   80:31244/TCP   13h
kubernetes        ClusterIP      10.3.240.1     <none>           443/TCP        18d
pages             ClusterIP      10.3.246.154   <none>           80/TCP         13h
photoresizer      ClusterIP      10.3.244.43    <none>           80/TCP         13h
photouploader     ClusterIP      10.3.245.255   <none>           80/TCP         13h
user-manager      ClusterIP      10.3.253.152   <none>           80/TCP         13h
```

5. Confirm that all application pods are running -- album-manager-* , auth-proxy-* , content-service-* , fake-s3-* , pages-* , photoresizer-* , user-manager-* and photouploader-* :
```
kubectl get pods
```
```
NAME                               READY     STATUS    RESTARTS   AGE
album-manager-5865f88699-mbcj4     2/2       Running   0          13h
auth-proxy-69d6995446-5h2pk        2/2       Running   0          13h
content-service-7968b4b584-nfpr6   2/2       Running   0          13h
fake-s3-5bd8ff94d6-rnpxl           2/2       Running   0          13h
pages-747f4dc4b7-q2w57             2/2       Running   0          13h
photoresizer-68cd8c5df4-vbpk5      2/2       Running   0          3h
photouploader-6bdbff759d-fsszg     2/2       Running   0          3h
user-manager-694cc64879-zr6ck      2/2       Running   0          3h
```



6. Open the Ingenious application in a browser using the following link:
```
http://<Public-IP-of-the-Ingress-Controller>/
```
### Uninstalling the Application
1. To uninstall application, run:
```
kubectl delete -f ingenious.yaml
```



