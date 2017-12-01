
# Directory structure
Files required for [installing Istio on a Kubernetes cluster](https://github.com/istio/istio/tree/master/install/kubernetes):

* [istio.yaml](https://github.com/istio/istio/blob/master/install/kubernetes/istio.yaml) - Use this file for installation without authentication enabled.
* [istio-auth.yaml](https://github.com/istio/istio/blob/master/install/kubernetes/istio-auth.yaml) - Use this file for installation with authentication enabled.
* [addons](https://github.com/istio/istio/blob/master/install/kubernetes/addons) - Directory contains optional components (Prometheus, Grafana, Service Graph, Zipkin, Zipkin to Stackdriver).

 Additional files to install NGiNX as a sidecar:
 
 * [updateVersion.sh](https://github.com/nginmesh/nginmesh/blob/release-doc-0.2.12/istio/release/updateVersion.sh) Use this file to regenerate Initializer file.
* [istio-initializer.yaml](https://github.com/nginmesh/nginmesh/blob/release-doc-0.2.12/istio/release/install/kubernetes/istio-initializer.yaml) - Use this file for installation of istio initializer for transparent injection.
* [templates](https://github.com/nginmesh/nginmesh/blob/release-doc-0.2.12/istio/release/install/kubernetes/templates) - Directory contains the template used to generate Initializer file.

