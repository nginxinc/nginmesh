FROM ubuntu:16.04

ARG KUBE_VERSION
#v1.9.4
ARG NGINMESH_VERSION
#0.6.0
ARG WRK_VERSION
#3.1.0
COPY cases /bookinfo_spec
WORKDIR /bookinfo_spec

# ADD KUBECTL
ADD https://storage.googleapis.com/kubernetes-release/release/$KUBE_VERSION/bin/linux/amd64/kubectl /usr/local/bin/kubectl

# ADD Istio
# ADD  https://git.io/getLatestIstio  /bookinfo_spec

# ADD dependencies
RUN apt-get update && \
    apt-get install -y --no-install-recommends -y build-essential  \
      git curl ca-certificates python python-setuptools python-pip libssl-dev nano
# ADD nginMesh
RUN  curl -L# https://github.com/nginmesh/nginmesh/releases/download/$NGINMESH_VERSION/nginmesh-$NGINMESH_VERSION.tar.gz  | tar zx
# ADD WRK
RUN (mkdir /opt/wrk && cd /opt/wrk && curl -L# https://github.com/wg/wrk/archive/$WRK_VERSION.tar.gz | tar zx --strip 1 && make && mv wrk /bin)
# CLEAN
RUN apt-get clean && apt-get -y remove build-essential curl && apt-get -y autoremove && \
rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*
RUN set -x && \
    chmod +x /usr/local/bin/kubectl
RUN ./install.sh
# CMD ["mamba","--format=documentation ."]
CMD ["tail","-f","/dev/null"]