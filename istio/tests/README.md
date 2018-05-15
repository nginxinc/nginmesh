### E2E test suite

This E2E test suites run series of functional and performance test to ensure nginMesh is working correctly and satisfy performance requirements.

### Prerequisites

Make sure below requirements are met:

| Version | Name | Details |
| --- | ------ | ------ |
|1.9|Kubernetes cluster|Without alpha feature, [link](https://istio.io/docs/setup/kubernetes/quick-start.html#google-kubernetes-engine)|
|0.6.0|Istio|[link](https://istio.io/docs/setup/kubernetes/quick-start.html)|
|0.6.0|nginMesh|[link](https://github.com/nginmesh/nginmesh/blob/master/README.md)|
|1.5.0|Bookinfo Application|[link](https://github.com/istio/istio/blob/master/samples/bookinfo/src)|
|0.9.2|Mamba|[link](https://github.com/nestorsalceda/mamba)|
|4.1.0|Wrecker|[link](https://github.com/wg/wrk)|
|2.7|Python|[link](https://www.python.org)|
|11.9.0|Pipenv|[link](https://docs.pipenv.org/)|
|1.10.0|Kubectl|[link](https://kubernetes.io/docs/tasks/tools/install-kubectl/)|

### Run Test 

1. Change directory to /cases directory in nginMesh repo:
```
cd tests/cases
```

2. Install Pipenv to create and manage a virtual environment for project:

```
python -m pip install pipenv
```

3. Install python dependencies in virtual environment using Pipenv:

```
./install.sh
```

4. Make sure required python dependencies are installed successfully:

```
pipenv graph
```
```
expects==0.8.0
grequests==0.3.0
  - gevent [required: Any, installed: 1.2.2]
    - greenlet [required: >=0.4.10, installed: 0.4.13]
  - requests [required: Any, installed: 2.18.4]
    - certifi [required: >=2017.4.17, installed: 2018.1.18]
    - chardet [required: >=3.0.2,<3.1.0, installed: 3.0.4]
    - idna [required: >=2.5,<2.7, installed: 2.6]
    - urllib3 [required: <1.23,>=1.21.1, installed: 1.22]
mamba==0.9.2
  - clint [required: Any, installed: 0.5.1]
    - args [required: Any, installed: 0.1.0]
  - coverage [required: Any, installed: 4.5.1]
``` 

5. Set app_namespace environment variable where Bookinfo application deployed in Kubernetes cluster: 

```
export app_namespace=default
```
Note: It will use "default" namespace as default, if not set.

6. Run all spec tests for Bookinfo application:

```
pipenv run mamba --format documentation .
```
```
              _                           _
   _ __   __ _(_)_ __  _ __ ___   ___  ___| |__
  | `_ \ / _  | |  _ \|  _   _ \ / _ \/ __| |_ \
  | | | | (_| | | | | | | | | | |  __/\__ \ | | |
  |_| |_|\__, |_|_| |_|_| |_| |_|\___||___/_| |_|
         |___/

Testing basic functionality
  Starting test
 | V1 Hit=2 | V2 Hit=4 | V3 Hit=4 | Total Hit=10 |
  44 requests in 1.01s, 232.67KB read
Requests/sec:     43.58
Transfer/sec:    230.45KB
    ✓ it Testing basic functionality (1.4042 seconds)

Testing route all requests to V1
  Set environment
    ✓ it Add routing rule (5.6762 seconds)
  Starting test
 | V1 Hit=10 | V2 Hit=0 | V3 Hit=0 | Total Hit=10 |
  56 requests in 1.02s, 250.66KB read
Requests/sec:     54.97
Transfer/sec:    246.05KB
    ✓ it Testing route all requests to V1 (1.3384 seconds)
  Clean environment
    ✓ it Delete routing rule (5.5004 seconds)

Testing route all requests to V3
  Set environment
    ✓ it Add routing rule (5.4385 seconds)
  Starting test
 | V1 Hit=0 | V2 Hit=0 | V3 Hit=10 | Total Hit=10 |
  39 requests in 1.02s, 224.25KB read
Requests/sec:     38.37
Transfer/sec:    220.65KB
    ✓ it Testing route all requests to V3 (1.4494 seconds)
  Clean environment
    ✓ it Delete routing rule (5.3219 seconds)

Testing route all requests to V1 and V3
  Set environment
    ✓ it Add routing rule (5.4349 seconds)
  Starting test
 | V1 Hit=6 | V2 Hit=0 | V3 Hit=4 | Total Hit=10 |
  48 requests in 1.02s, 251.79KB read
Requests/sec:     47.10
Transfer/sec:    247.07KB
    ✓ it Testing route all requests to V1 and V3 (1.3605 seconds)
  Clean environment
    ✓ it Delete routing rule (5.3254 seconds)

Testing destination-weight, route to V1-75%, V2-25%
  Set environment
    ✓ it Add routing rule (5.4592 seconds)
  Starting test
 | V1 Hit=7 | V2 Hit=3 | V3 Hit=0 | Total Hit=10 |
  57 requests in 1.02s, 265.35KB read
Requests/sec:     56.16
Transfer/sec:    261.43KB
    ✓ it Bookinfo destination-weight test, route to V1-75%, V2-25% (1.3863 seconds)
  Clean environment
    ✓ it Delete routing rule (5.3272 seconds)

Testing route all requests to V2 and V3
  Set environment
    ✓ it Add routing rule (5.4347 seconds)
  Starting test
 | V1 Hit=0 | V2 Hit=4 | V3 Hit=6 | Total Hit=10 |
  36 requests in 1.03s, 207.05KB read
Requests/sec:     34.99
Transfer/sec:    201.25KB
    ✓ it Testing route all requests to V2 and V3 (1.4698 seconds)
  Clean environment
    ✓ it Delete routing rule (5.3268 seconds)

Testing HTTP Redirect
  Set environment
    ✓ it Add routing rule (5.4378 seconds)
  Starting test
  180 requests in 1.02s, 26.02KB read
Requests/sec:    176.51
Transfer/sec:     25.51KB
    ✓ it Testing HTTP Redirect (2.2186 seconds)
  Clean environment
    ✓ it Delete routing rule (7.0557 seconds)

Testing HTTP Retry
  Set environment
    ✓ it Add routing rule (5.4625 seconds)
  Starting test
Total Retry Hit=10
  36 requests in 1.02s, 207.00KB read
Requests/sec:     35.16
Transfer/sec:    202.14KB
    ✓ it Testing HTTP Retry (16.0645 seconds)
  Clean environment
    ✓ it Delete routing rule (5.4433 seconds)

Testing route "jason" user to V2
  Set environment
    ✓ it Add routing rule (5.4366 seconds)
  Starting test
 | V1 Hit=0 | V2 Hit=10 | V3 Hit=0 | Total Hit=10 |
  57 requests in 1.03s, 304.88KB read
Requests/sec:     55.57
Transfer/sec:    297.21KB
    ✓ it Testing route "jason" user to V2 (2.6802 seconds)
  Clean environment
    ✓ it Delete routing rule (5.4189 seconds)

Testing Kafka messages
 | V1 Hit=3 | V2 Hit=4 | V3 Hit=3 | Total Hit=10 |
  29 requests in 1.02s, 154.05KB read
Requests/sec:     28.41
Transfer/sec:    150.90KB

Starting test
Processed a total of 10 messages
  ✓ it Testing Kafka functionality (6.5686 seconds)

Zipkin tracing functionality
  Set environment
Forwarding from 127.0.0.1:9411 -> 9411
    ✓ it Add Zipkin tracing feature (20.7318 seconds)
  Starting Test
Handling connection for 9411
 | V1 Hit=3 | V2 Hit=3 | V3 Hit=4 | Total Hit=10 |
  39 requests in 1.03s, 207.74KB read
Requests/sec:     38.00
Transfer/sec:    202.42KB
    ✓ it Bookinfo Zipkin tracing feature (3.4614 seconds)
  Clean Environment
    ✓ it Delete Zipkin tracing feature (14.1853 seconds)

29 examples ran in 164.8122 seconds
```
7. To run selectively, please input one or multiple test cases:
```
pipenv run mamba --format documentation 1_bd_spec.py 2_route_all_v1_spec.py

```
```
               _                           _
   _ __   __ _(_)_ __  _ __ ___   ___  ___| |__
  | `_ \ / _  | |  _ \|  _   _ \ / _ \/ __| |_ \
  | | | | (_| | | | | | | | | | |  __/\__ \ | | |
  |_| |_|\__, |_|_| |_|_| |_| |_|\___||___/_| |_|
         |___/

Testing basic functionality
  Starting test
 | V1 Hit=3 | V2 Hit=3 | V3 Hit=4 | Total Hit=10 |
  68 requests in 1.02s, 361.75KB read
Requests/sec:     66.82
Transfer/sec:    355.44KB
    ✓ it Testing basic functionality (1.3689 seconds)

Testing route all requests to V1
  Set environment
    ✓ it Add routing rule (5.6165 seconds)
  Starting test
 | V1 Hit=10 | V2 Hit=0 | V3 Hit=0 | Total Hit=10 |
  67 requests in 1.01s, 299.87KB read
Requests/sec:     66.16
Transfer/sec:    296.09KB
    ✓ it Testing route all requests to V1 (1.3311 seconds)
  Clean environment
    ✓ it Delete routing rule (5.4916 seconds)

4 examples ran in 15.9301 seconds
```
