# Demo nginMesh streaming using KSQL 

This demo focuses on building real-time analytics of users inside nginMesh enabled cluster. Please, refer to this [link](https://github.com/confluentinc/ksql/tree/master/ksql-clickstream-demo) for details of used demo.

## Quick Start

### Prerequisites

Make sure below requirements are met:
  
  | Version | Name | Details |
  | --- | ------ | ------ |
  |1.9|Kubernetes cluster|Without alpha feature, [link](https://istio.io/docs/setup/kubernetes/quick-start.html#google-kubernetes-engine)|
  |0.7.0|Istio|[link](https://istio.io/docs/setup/kubernetes/quick-start.html)|
  |0.7.0|nginMesh|[link](https://github.com/nginmesh/nginmesh/blob/master/README.md)|
  |1.5.0|Bookinfo Application|[link](https://github.com/istio/istio/blob/master/samples/bookinfo/src)|
  |1.1.0|Kafka|[link](https://kafka.apache.org/downloadsc)|
  |2.9.0|Helm|[link](https://docs.helm.sh/using_helm/)|
  |4.1.0| Wrk|[link](https://github.com/wg/wrk)| 

### Install 

1. Install Kafka connect in the namespace 'kafka':
```
kubectl create -f ../../install/kafka/connect.yml
```

2. Install KSQL server in the namespace 'kafka':
```
kubectl create -f ../../install/kafka/ksql.yml
```

3. Download and install KSQL client from this [link](https://www.confluent.io/download/). Update $PATH variable either in .bash_profile or .bashrc to include /bin directory of KSQL.

4. Install Elasticsearch  in the namespace 'elastic':
```
./install-elastic.sh
```

5. Install Grafana in the namespace 'kafka':
```
./install-grafana.sh
```
6. Make sure all pods and services are up and running:
```
kubectl get pods -n kafka
```
```
NAME                             READY     STATUS    RESTARTS   AGE
grafana-8679d8f6b9-xvn6p         1/1       Running   0          14h
kconnect                         1/1       Running   0          14h
ksql-cli                         1/1       Running   0          3d
my-kafka-kafka-0                 1/1       Running   1          11d
my-kafka-kafka-1                 1/1       Running   0          3d
my-kafka-kafka-2                 1/1       Running   0          11d
my-kafka-zookeeper-0             1/1       Running   0          11d
my-kafka-zookeeper-1             1/1       Running   1          11d
my-kafka-zookeeper-2             1/1       Running   0          11d
testclient                       1/1       Running   0          11d
```
```
kubectl get pods -n elastic
```
```
NAME                                            READY     STATUS    RESTARTS   AGE
elastic-elasticsearch-client-5c8f946c87-79s5g   1/1       Running   2          3d
elastic-elasticsearch-client-5c8f946c87-vvkc2   1/1       Running   3          3d
elastic-elasticsearch-data-0                    1/1       Running   0          3d
elastic-elasticsearch-data-1                    1/1       Running   1          3d
elastic-elasticsearch-master-0                  1/1       Running   0          3d
elastic-elasticsearch-master-1                  1/1       Running   1          3d
elastic-elasticsearch-master-2                  1/1       Running   0          3d
```

7. Run following script to activate port-forwarding for ksql, connect, elasticsearch and grafana to localhost:

```
./all-portforward.sh
```

8. Run following script to create nginMesh stream and tables which will be pushed to Elasticsearch and Grafana:
```
./install_metrics.sh
```

### Visualize metrics in Grafana

1. Run following script to get password for grafana admin user:
```
./grafana-password.sh
```

2. Access to Grafana Dashboard in [http://localhost:3000](http://localhost:3000/) from browser using retrieved credentials.

3. Run below script to import nginMesh dashboard to Grafana:
```
./grafana-upload-dashboard.sh
```

4. Generate requests towards sample application deployed and make sure all widgets show data accordingly. Below script could be used in Istio default Bookinfo application case:
```
./gen_load.sh
```

![Alt text](images/dashboard.png?raw=true "Grafana Dashboard")

### Unistall 

1. Unistall Kafka connect:
```
kubectl delete -f ../../install/kafka/connect.yml
```

2. Unistall KSQL server:
```
kubectl delete -f ../../install/kafka/ksql.yml
```

3. Uninstall Elasticsearch:
```
helm del --purge elastic;
```

4. Uninstall Grafana:
```
helm del --purge grafana;
```

5. Deactivate all port-forwardings:
```
pkill -f port-forward
```
