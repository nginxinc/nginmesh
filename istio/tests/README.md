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
 | V1 Hit=4 | V2 Hit=2 | V3 Hit=4 | Total Hit=10 |
  52 requests in 1.02s, 284.10KB read
Requests/sec:     51.07
Transfer/sec:    279.03KB
    ✓ it Testing basic functionality (1.4355 seconds)

Testing route all requests to V1
  Set environment
    ✓ it Add routing rule (5.6113 seconds)
  Starting test
 | V1 Hit=10 | V2 Hit=0 | V3 Hit=0 | Total Hit=10 |
  63 requests in 1.01s, 281.98KB read
Requests/sec:     62.25
Transfer/sec:    278.61KB
    ✓ it Testing route all requests to V1 (1.3310 seconds)
  Clean environment
    ✓ it Delete routing rule (5.5135 seconds)

Testing route all requests to V3
  Set environment
    ✓ it Add routing rule (5.4327 seconds)
  Starting test
 | V1 Hit=0 | V2 Hit=0 | V3 Hit=10 | Total Hit=10 |
  77 requests in 1.01s, 442.68KB read
Requests/sec:     76.23
Transfer/sec:    438.27KB
    ✓ it Testing route all requests to V3 (1.3356 seconds)
  Clean environment
    ✓ it Delete routing rule (5.3283 seconds)

Testing route all requests to V1 and V3
  Set environment
    ✓ it Add routing rule (5.4259 seconds)
  Starting test
 | V1 Hit=3 | V2 Hit=0 | V3 Hit=7 | Total Hit=10 |
  71 requests in 1.02s, 359.79KB read
Requests/sec:     69.57
Transfer/sec:    352.56KB
    ✓ it Testing route all requests to V1 and V3 (1.3478 seconds)
  Clean environment
    ✓ it Delete routing rule (5.3143 seconds)

Testing destination-weight, route to V1-75%, V2-25%
  Set environment
    ✓ it Add routing rule (5.4597 seconds)
  Starting test
 | V1 Hit=6 | V2 Hit=4 | V3 Hit=0 | Total Hit=10 |
  77 requests in 1.01s, 375.28KB read
Requests/sec:     76.26
Transfer/sec:    371.68KB
    ✓ it Bookinfo destination-weight test, route to V1-75%, V2-25% (1.3235 seconds)
  Clean environment
    ✓ it Delete routing rule (5.3187 seconds)

Testing route all requests to V2 and V3
  Set environment
    ✓ it Add routing rule (5.4383 seconds)
  Starting test
 | V1 Hit=0 | V2 Hit=1 | V3 Hit=9 | Total Hit=10 |
  69 requests in 1.02s, 396.80KB read
Requests/sec:     67.62
Transfer/sec:    388.85KB
    ✓ it Testing route all requests to V2 and V3 (1.3597 seconds)
  Clean environment
    ✓ it Delete routing rule (5.3194 seconds)

Testing HTTP Redirect
  Set environment
    ✓ it Add routing rule (5.4329 seconds)
  Starting test
  179 requests in 1.02s, 25.87KB read
Requests/sec:    174.77
Transfer/sec:     25.26KB
    ✓ it Testing HTTP Redirect (2.2081 seconds)
  Clean environment
    ✓ it Delete routing rule (5.3181 seconds)

Testing HTTP Retry
  Set environment
    ✓ it Add routing rule (5.4503 seconds)
  Starting test
Total Retry Hit=10
  77 requests in 1.02s, 442.68KB read
Requests/sec:     75.45
Transfer/sec:    433.79KB
    ✓ it Testing HTTP Retry (16.0318 seconds)
  Clean environment
    ✓ it Delete routing rule (5.3102 seconds)

Testing route "jason" user to V2
  Set environment
    ✓ it Add routing rule (5.4323 seconds)
  Starting test
 | V1 Hit=0 | V2 Hit=10 | V3 Hit=0 | Total Hit=10 |
  76 requests in 1.01s, 403.91KB read
Requests/sec:     75.47
Transfer/sec:    401.08KB
    ✓ it Testing route "jason" user to V2 (2.4717 seconds)
  Clean environment
    ✓ it Delete routing rule (5.3364 seconds)

25 examples ran in 129.3463 seconds
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
