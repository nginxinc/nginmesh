
#!/usr/bin/env bash
#curl -X "GET" "http://localhost:3000/api/dashboards/db/click-stream-analysis" \
#        -H "Content-Type: application/json" \
#	     --user admin:admin


echo "Loading Grafana nginMesh Dashboard"
grafana_password=$(kubectl get secret --namespace kafka grafana -o jsonpath="{.data.admin-password}" | base64 --decode)

curl -X "POST" "http://localhost:3000/api/dashboards/db" \
	    -H "Content-Type: application/json" \
	     --user admin:$grafana_password \
	     --data-binary @grafana-dashboard.json
