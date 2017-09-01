# Running BookInfo on Nginx Proxy

This app configuration only works with 0.16 release of Istio.

Much of the contents are borrowed from Istio sample app documentation: https://istio.io/docs/samples/bookinfo.html
with customization for Nginx proxy.

## Before You begin

* Please have at least 4 worker nodes with at least 4G of memory
* Setup Istio by following instructions in the https://istio.io/docs/tasks/installing-istio.html


## Start the application

1.  Bring up the application containers

    ```bash
     ./deploy.sh
    ```
    The above commands launches four microservices and create the gateway ingress resources.
    The reviews microservice has 3 versions: v1, v2, v3
    
    (Optional) The script "deploy.sh" call Mac OSX version of the nginx inject utility which generates sidecar configuration for application.
    If you are using other OS, please checkout out https://github.com/nginmesh/istio-nginx-inject and build the binary again.
    
1.  Confirm all services and pods are correctly defined and running:

    ```bash
    kubectl get services
    ```    
    
    which produces similar output:
    
    ```bash
       NAME                       CLUSTER-IP   EXTERNAL-IP   PORT(S)              AGE
       details                    10.0.0.31    <none>        9080/TCP             6m
       istio-ingress              10.0.0.122   <pending>     80:31565/TCP         8m
       istio-pilot                10.0.0.189   <none>        8080/TCP             8m
       istio-mixer                10.0.0.132   <none>        9091/TCP,42422/TCP   8m
       kubernetes                 10.0.0.1     <none>        443/TCP              14d
       productpage                10.0.0.120   <none>        9080/TCP             6m
       ratings                    10.0.0.15    <none>        9080/TCP             6m
       reviews                    10.0.0.170   <none>        9080/TCP             6m
       ```
       
    and
    
    ```bash
       kubectl get pods
       ```
       
       which produces
       
       ```bash
       NAME                                        READY     STATUS    RESTARTS   AGE
       details-v1-1520924117-48z17                 2/2       Running   0          6m
       istio-ingress-3181829929-xrrk5              1/1       Running   0          8m
       istio-pilot-175173354-d6jm7                 2/2       Running   0          8m
       istio-mixer-3883863574-jt09j                2/2       Running   0          8m
       productpage-v1-560495357-jk1lz              2/2       Running   0          6m
       ratings-v1-734492171-rnr5l                  2/2       Running   0          6m
       reviews-v1-874083890-f0qf0                  2/2       Running   0          6m
       reviews-v2-1343845940-b34q5                 2/2       Running   0          6m
       reviews-v3-1813607990-8ch52                 2/2       Running   0          6m
       ```
       
       
1. Determine the gateway ingress URL:

   ```bash
   kubectl get ingress -o wide
   ```
   
   ```bash
   NAME      HOSTS     ADDRESS                 PORTS     AGE
   gateway   *         130.211.10.121          80        1d
   ```

   If your Kubernetes cluster is running in an environment that supports external load balancers,
   and the Istio ingress service was able to obtain an External IP, the ingress resource ADDRESS will be equal to the
   ingress service external IP.

   ```bash
   export GATEWAY_URL=130.211.10.121:80
   ```
   
   > Sometimes when the service is unable to obtain an external IP, the ingress ADDRESS may display a list
   > of NodePort addresses. In this case, you can use any of the addresses, along with the NodePort, to access the ingress. 
   > If, however, the cluster has a firewall, you will also need to create a firewall rule to allow TCP traffic to the NodePort.
   > In GKE, for instance, you can create a firewall rule using the following command:
   > ```bash
   > gcloud compute firewall-rules create allow-book --allow tcp:$(kubectl get svc istio-ingress -o jsonpath='{.spec.ports[0].nodePort}')
   > ```

   If your deployment environment does not support external load balancers (e.g., minikube), the ADDRESS field will be empty.
   In this case you can use the service NodePort instead:
   
   ```bash
   export GATEWAY_URL=$(kubectl get po -l istio=ingress -o 'jsonpath={.items[0].status.hostIP}'):$(kubectl get svc istio-ingress -o 'jsonpath={.spec.ports[0].nodePort}')
   ```
       
1. Confirm that the BookInfo application is running with the following `curl` command:

   ```bash
   curl -o /dev/null -s -w "%{http_code}\n" http://${GATEWAY_URL}/productpage
   ```
   ```bash
   200
   ```
   
## Cleanup

When you're finished experimenting with the BookInfo sample, you can uninstall it as follows:

1. Delete the routing rules and terminate the application pods

   ```bash
   ./cleanup.sh
   ```

1. Confirm shutdown

   ```bash
   istioctl get route-rules   #-- there should be no more routing rules
   kubectl get pods           #-- the BookInfo pods should be deleted
   ```

## Content-based Routing

Because the BookInfo sample deploys 3 versions of the reviews microservice, we need to set a default route. 
Otherwise if you access the application several times, you'll notice that sometimes the output contains star ratings. 
This is because without an explicit default version set, Istio will route requests to all available versions of a service in a random fashion.

> Note: This task assumes you don't have any routes set yet. If you've already created conflicting route rules for the sample, you'll need to use `replace` rather than `create` in one or both of the following commands.

1. Set the default version for all microservices to v1.

    ```bash
   istioctl create -f route-rule-all-v1.yaml
   ```

   Since rule propagation to the proxies is asynchronous, you should wait a few seconds for the rules
   to propagate to all pods before attempting to access the application.
   
   
1. Open the BookInfo URL (http://$GATEWAY_URL/productpage) in your browser

   You should see the BookInfo application productpage displayed.
   Notice that the `productpage` is displayed with no rating stars since `reviews:v1` does not access the ratings service.


1. Route a specific user to `reviews:v2`

   Lets enable the ratings service for test user "jason" by routing productpage traffic to
   `reviews:v2` instances.

   ```bash
   istioctl create -f route-rule-reviews-test-v2.yaml
   ```

   Confirm the rule is created:

   ```bash
   istioctl get route-rule reviews-test-v2
   ```
   ```yaml
   destination: reviews.default.svc.cluster.local
   match:
     httpHeaders:
       cookie:
         regex: ^(.*?;)?(user=jason)(;.*)?$
   precedence: 2
   route:
   - tags:
       version: v2
   ```

1. Log in as user "jason" at the `productpage` web page.

   You should now see ratings (1-5 stars) next to each review. Notice that if you log in as
   any other user, you will continue to see `reviews:v1`.