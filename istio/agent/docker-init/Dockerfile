FROM debian:stretch-slim
RUN apt-get update && apt-get install -y iptables
ADD prepare_proxy.sh /
ENTRYPOINT ["/prepare_proxy.sh"]
