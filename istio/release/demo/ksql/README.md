# KSQL Demo

### Install Elastic Search

1. Use Helm to install Elastic Search
```
 helm repo add incubator http://storage.googleapis.com/kubernetes-charts-incubator
 helm install --name elastic incubator/elasticsearch --namespace elastic --set data.persistence.enabled=false,master.persistence.enabled=false --set rbac.create=true
```
This install elastic 5.4 in the namespace 'elastic'
2. Run following script to setup Kafka. It is installed in 'kafka' namespace.  It is possible to use your existing kafka installation.
```
nginmesh-0.6.0/install/kafka/install.sh
```