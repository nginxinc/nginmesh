# NGINX dynamic module for Istio Mixer 


## Dependencies

Build system use Docker to generate the module binary.

## Compatibility

* 1.11.x (last tested with 1.13.5)


## Synopsis

```nginx

 http   {
 
 	 mixer_server istio-mixer.istio-system;
     mixer_port   9091;

 
	 server {
	 
	      mixer_source_ip     10.0.0.0;
          mixer_source_uid    kubernetes://productpage-v1-2213572757-758cs.beta1;
          mixer_source_service productpage.beta1.svc.cluster.local;
          mixer_destination_service abc.ns.svc.cluster.local;
          mixer_destination_uid details;
         
         
          location /  {
              mixer_report on;
              proxy_pass http://service1;
          }
         
            
			
	 }
		
 }	

```


## Directives

### mixer_server

| -   | - |
| --- | --- |
| **Syntax**  | **mixer_server** <mixer_url_or_ip> |
| **Default** | - |
| **Context** | http |

`Description:` Specify the mixer server address


### mixer_port

| -   | - |
| --- | --- |
| **Syntax**  | **mixer_port** <port_number> |
| **Default** | - |
| **Context** | http |

`Description:` Specify the mixer server port


### mixer_source_ip

| -   | - |
| --- | --- |
| **Syntax**  | **mixer_source_ip** <ip_address> |
| **Default** | - |
| **Context** | server, location  |
| **Mixer attribute** | source.ip  |

`Description:` Standard mixer attribute **Client IP address**

### mixer_source_uid

| -   | - |
| --- | --- |
| **Syntax**  | **mixer_source_uid** <kubernetes client service id> |
| **Default** | - |
| **Context** | server, location  |
| **Mixer attribute** | source.uid  |

`Description:` Standard mixer attribute **Platform-specific unique identifier for the client instance of the source service**

### mixer_source_service

| -   | - |
| --- | --- |
| **Syntax**  | **mixer_source_service** <kubernetes client service name> |
| **Default** | - |
| **Context** | server, location  |
| **Mixer attribute** | source.service  |

`Description:` Standard mixer attribute **The fully qualified name of the service that the client belongs to**


### mixer_source_port

| -   | - |
| --- | --- |
| **Syntax**  | **mixer_source_port** <kubernetes client service name> |
| **Default** | - |
| **Context** | server, location  |
| **Mixer attribute** | source.service  |

`Description:` Standard mixer attribute **The fully qualified name of the service that the client belongs to**


### mixer_destination_service

| -   | - |
| --- | --- |
| **Syntax**  | **mixer_destination_service** <kubernetes destination service name> |
| **Default** | - |
| **Context** | server, location  |
| **Mixer attribute** | destination.servicee  |

`Description:` Standard mixer attribute **The fully qualified name of the service that the server belongs to**

### mixer_destination_uid

| -   | - |
| --- | --- |
| **Syntax**  | **mixer_destination_service** <kubernetes destination service uid> |
| **Default** | - |
| **Context** | server, location  |
| **Mixer attribute** | destination.uid  |

`Description:` Standard mixer attribute **Platform-specific unique identifier for the server instance of the destination service.**



## Installation

1. Clone the git repository

  ```
  shell> git clone git@github.com:nginmesh/ngx-istio-mixer.git
  ```

2. Build the dynamic module

  ```
  shell> make build-base;make build-module
  ```

  This copies the generated .so file into module/release directory

3. Test deploy in the local K8 cluster

  First deploy test nginx server in the cluster
  ```make test-k8-only```  
  Then test using ```make test-http-report'''


## Running unit test

```bash
make test-unit
```
