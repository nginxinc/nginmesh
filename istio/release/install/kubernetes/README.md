# Optional addons installation

1. Enable [Zipkin](https://istio.io/docs/tasks/telemetry/distributed-tracing.html#accessing-the-dashboard):
```
kubectl apply -f istio-0.2.12/install/kubernetes/addons/zipkin.yaml
kubectl port-forward -n istio-system $(kubectl get pod -n istio-system -l app=zipkin -o jsonpath='{.items[0].metadata.name}') 9411:9411 &
```
Then open in browser: http://localhost:9411

2. Enable [Graphana](https://istio.io/docs/tasks/telemetry/using-istio-dashboard.html):
```
kubectl apply -f istio-0.2.12/install/kubernetes/addons/grafana.yaml
kubectl -n istio-system port-forward $(kubectl -n istio-system get pod -l app=grafana -o jsonpath='{.items[0].metadata.name}') 3000:3000 &
```
Then open in browser: http://localhost:3000/dashboard/db/istio-dashboard 

3. Enable [Prometheus](https://istio.io/docs/tasks/telemetry/querying-metrics.html):
```
kubectl apply -f istio-0.2.12/install/kubernetes/addons/prometheus.yaml
kubectl -n istio-system port-forward $(kubectl -n istio-system get pod -l app=prometheus -o jsonpath='{.items[0].metadata.name}') 9090:9090 &
```
Then open in browser: http://localhost:9090/graph 

4. Enable [ServiceGraph](https://istio.io/docs/tasks/telemetry/servicegraph.html):
```
kubectl apply -f istio-0.2.12/install/kubernetes/addons/servicegraph.yaml
kubectl -n istio-system port-forward $(kubectl -n istio-system get pod -l app=servicegraph -o jsonpath='{.items[0].metadata.name}') 8088:8088 &
```
Then open in browser:  http://localhost:8088/dotviz

# Uninstall
```
kubectl delete -f istio-0.2.12/install/kubernetes/addons/zipkin.yaml #-- Delete Zipkin
kubectl delete -f istio-0.2.12/install/kubernetes/addons/grafana.yaml #-- Delete Graphana
kubectl delete -f istio-0.2.12/install/kubernetes/addons/prometheus.yaml #-- Delete Prometheus
kubectl delete -f istio-0.2.12/install/kubernetes/addons/servicegraph.yaml #Delete ServiceGraph
killall kubectl #-- Remove any kubectl port-forward processes that may be running
```

#  Verify uninstall
```
kubectl get pods -n istio-system #-- Istio pods should be deleted
kubectl get svc  -n istio-system #-- Istio services should be deleted
```
