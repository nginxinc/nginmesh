# The Sidecar Implementation

The NGINX sidecar is implemented using the agent application, which runs alongside NGINX and ensures that NGINX configuration always matches the current state of the Istio environment. In this document we discuss the implementation in details.

## What is the Agent?

The agent is an application written in Go. It is deployed within the same container as NGINX. The agent has the following responsibilities:
1. *Controlling NGINX*, which includes starting/stopping NGINX and applying new configuration (reloading). If NGINX suddenly terminates (crashes), the agent terminates as well.
1. *Monitoring Pilot for configuration changes.* The agent constantly monitors Pilot for changes in load balancing configuration for the service instance which NGINX belongs too.
1. *Converting Pilot (Envoy) configuration to NGINX configuration.* When configuration is updated in Pilot, the agent converts the updated configuration and applies it to NGINX.
1. *Monitoring the Istio Auth certificates and keys on the file system.* The agent monitors the folder where Istio Auth deploys the certificates and keys which are used for mutual TLS authentication between NGINX proxies. When those files are updated by Istio Auth, the agent reloads NGINX to apply the updates.

## Supported Traffic Management Features

Currently, not all Istio traffic management features (routing rules and destination policies) are supported. The two tables below show the support status of each feature.

### Routing Rules

| Rule | Is It Supported? |
| --- | --- |
| MatchCondition | Yes |
| DestinationWeight | Yes |
| HTTPRedirect | Yes |
| HTTPRewrite | Yes |
| HTTPTimeout | Yes |
| HTTPRetry | Yes |
| HTTPFaultInjection | Yes, but with limitations: Delay is not supported. |


### Destination Policies


| Policies | Is it Supported? |
| --- | --- |
| LoadBalancing - SimpleLBPolicy | Yes |
| CircuitBreaker - SimpleCircuitBreakerPolicy | Yes, but with limitations: only `maxConnections`, `httpConsecutiveErrors` , `httpDetectionInterval` and `sleepWindow` are supported; `sleepWindow` must be equal to `httpDetectionInterval` |

## Limitations

The NGINX sidecar has the following limitations:
* The agent does not support configuring NGINX as an Ingress controller or an Egress controller.
* The agent works only when Istio is deployed on Kubernetes.
* When mutual TLS authentication is enabled, an NGINX-client does not validate the identity of an NGINX-server -- it cannot validate the service account presented in the NGINX-server certificate.

## Building the Sidecar 

To build the NGINX sidecar, you need to have the following software installed on your machine:
* Docker
* Make

Also, you need to have access to a docker registry for uploading docker images. 

During the build, the following images are built and uploaded to your docker registry:
* *proxy_debug*, which comes with the agent and NGINX.
* *proxy_init*, which is used for configuring IPtables rules for transparently injecting an NGINX proxy from the `proxy_debug` image into an application pod.  

Additionally, the `tracing_builder` image is built. It is used during the build for compiling NGINX third-party modules required for tracing. The image is not uploaded to the docker registry. 

To build and upload the images, run the following command:
```
$ make clean
$ make release REPO=<your-docker-repo>
```

The `<your-docker-repo>/proxy_debug:<release-tag>` and `<your-docker-repo>/proxy_init:<release-tag>` images will be uploaded to your docker registry. Not that the `<release-tag>` is set in the `Makefile`.

