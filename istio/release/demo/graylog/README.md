
# Demo nginMesh streaming using Graylog

Graylog is a powerful log management and analysis tool that has many use cases, from monitoring to debugging applications.

It has 3 main components:

**Server nodes:** Serves as a worker that receives and processes messages, and communicates with all other non-server components. Its performance is CPU dependent. <br>
**Elasticsearch nodes:** Stores all of the logs/messages. Its performance is RAM and disk I/O dependent.<br>
**MongoDB:** Stores metadata and does not experience much load.

## Architecture

![Alt text](images/graylog.png?raw=true "Graylog Architecture") 

Please, check [link](http://docs.graylog.org/en/2.4/) for documentation.

## Quick Start
Below are instructions to quickly install and configure  Graylog v2.4.3.

### Prerequisites
Make sure below requirements are met:
  
  | Version | Name | Details |
  | --- | ------ | ------ |
  |1.9|Kubernetes cluster|Without alpha feature, [link](https://istio.io/docs/setup/kubernetes/quick-start.html#google-kubernetes-engine)|
  |0.7.0|Istio|[link](https://istio.io/docs/setup/kubernetes/quick-start.html)|
  |0.7.0|nginMesh|[link](https://github.com/nginmesh/nginmesh/blob/master/README.md)|
  |1.5.0|Bookinfo Application|[link](https://github.com/istio/istio/blob/master/samples/bookinfo/src)|
  |1.1.0|Kafka|[link](https://kafka.apache.org/downloadsc)|

### Install Graylog
1. Install graylog deployment in graylog namespace:
```
 kubectl create -f graylog.yaml
```
2. Make sure following pods are up and running:

```
kubectl get pods -n graylog
```
```
NAME                            READY     STATUS    RESTARTS   AGE
elasticsearch-97c476698-7tmpd   1/1       Running   0          1m
graylog-c4d976795-vfhpf         1/1       Running   0          1m
mongo-6bb464754d-d6fd8          1/1       Running   0          1m
```
3. Make sure following services are up and running: 
```
kubectl get svc -n graylog
```
```
NAME            TYPE           CLUSTER-IP     EXTERNAL-IP       PORT(S)                          AGE
elasticsearch   ClusterIP      None           <none>            55555/TCP                        2m
graylog         LoadBalancer   10.55.242.76   100.100.100.100   9000:31927/TCP,12201:30371/TCP   2m
mongo           ClusterIP      None           <none>            55555/TCP                        2m

```

4. Activate port-forwarding for running graylog pod:
```
./graylog-portforward.sh
```

5. Access to Graylog Dashboard from browser using default credentials:
 ```
username: admin
password: somesaltpassword
```
```
http://127.0.0.1:9000/
```
![Alt text](images/1_login.png?raw=true "Login")
Note: Check graylog deployment file for username/password passed as environment variable. 


### Configure Kafka

1. Select Content Packs tab from System menu:

![Alt text](images/2_content_packs.png?raw=true "Content Packs")

2. Upload  [nginmesh_kafka_package.json](nginmesh_kafka_pacskage.json) file which contains all configuration related to Kafka input/extractor/dashboard:

![Alt text](images/3_upload_package.png?raw=true "Upload Packs")

3. Apply content of package:

![Alt text](images/4_apply_content.png?raw=true "Apply Content ")

4. Generate few requests towards sample application deployed and monitor dashboard widgets in Dashboards menu:

![Alt text](images/5_dashboard.png?raw=true "Dashboard ")

### Configure Geo-Location plugin

1. Install Map database provided by MaxMind:
```
./install_map.sh
```
Note: Please, refer for plugin [details](http://docs.graylog.org/en/2.4/pages/geolocation.html).

2. Enable Geo-Location processor and set path to "/usr/share/graylog/plugin/GeoLite2-City.mmdb" in System/Configurations menu:

![Alt text](images/6_geoloc_proc.png?raw=true "GeoLoc Processor ")

3. Enable Message processors in below order in System/Configurations menu:

![Alt text](images/7_message_proc.png?raw=true "Message Processor ")


### Uninstalling the Graylog
1. To uninstall Graylog deployment, run:
``` 
kubectl delete -f graylog.yaml
```

