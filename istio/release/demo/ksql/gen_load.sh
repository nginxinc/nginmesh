GW_IP=$(kubectl get svc -n istio-system | grep -E 'istio-ingress' | awk '{ print $4 }')
wrk -t2 -c10 -d10s "http://$GW_IP/productpage"