# List of handy scripts to work with nginmesh/istio


## Find mixer pod in the istio-system namespace
```bash
kubectl get pod -l istio=mixer -n istio-system -o jsonpath='{.items[0].metadata.name}'
```

## Logs from mixer pod
```bash
kubectl logs -f $(kubectl get pod -l istio=mixer -n istio-system -o jsonpath='{.items[0].metadata.name}')  -n istio-system -c mixer
```

## mixer address
```
istio-mixer.istio-system:9093
```



