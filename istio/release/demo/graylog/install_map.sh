curl -L http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.tar.gz | tar zx
graylog_pod=$(kubectl -n graylog get pod -l service=graylog -o jsonpath='{.items[0].metadata.name}')
kubectl cp GeoLite2-City_*/GeoLite2-City.mmdb graylog/$graylog_pod:/usr/share/graylog/plugin/
rm -rf  GeoLite2-City_*
