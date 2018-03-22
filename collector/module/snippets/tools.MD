# Snippets


## Mixer

port forward mixer server to localhost

```bash
kubectl port-forward $(kubectl get pod -l istio=mixer -n istio-system -o jsonpath='{.items[0].metadata.name}')  -n istio-system 9091:9091
```

watch logs
```bash
kubectl logs -f $(kubectl get pod -l istio=mixer -n istio-system -o jsonpath='{.items[0].metadata.name}')  -n istio-system -c mixer
```