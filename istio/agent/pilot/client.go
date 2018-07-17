package pilot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang/glog"
)

// Client allows getting load balancing configuration from Pilot.
type Client struct {
	endpoint       string
	httpClient     *http.Client
	serviceCluster string
	serviceNode    string
	podIP          string
	collectorAddress	string
	collectorTopic string
}

// ProxyConfig represents full load balancing configuration for a sidecar proxy.
type ProxyConfig struct {
	HTTPListeners    map[string]Listener
	TCPListeners     map[string]Listener
	HTTPRouteConfigs map[string]HTTPRouteConfig
	Services         map[string]ServiceHosts
	Clusters         map[string]Cluster
	TargetService    string
}

// NewClient creates a new Client.
func NewClient(endpoint string, httpClient *http.Client, serviceNode string, serviceCluster string, podIP string,
	collectorAddress string,collectorTopic string) *Client {
	return &Client{fmt.Sprintf("http://%v", endpoint), httpClient, serviceNode, serviceCluster, podIP,collectorAddress,collectorTopic}
}

func (c *Client) getListeners() (Listeners, error) {
	url := fmt.Sprintf("%v/v1/listeners/%v/%v", c.endpoint, c.serviceCluster, c.serviceNode)
	glog.Infof("listener url: %v", url)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("couldn't get listeners: %v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if glog.V(3) {
		glog.Infof("Response from %v: %v", url, string(body))
	}
	//glog.Infof("Response from %v: %v", url, string(body))

	return unMarshalListeners(body)
}

func unMarshalListeners(body []byte)  (Listeners, error)  {

	var res ldsResponse

	err := json.Unmarshal(body, &res)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal response from pilot: %v", err)
	}

	err = finishUnmarshallingListeners(res.Listeners)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal response from pilot: %v", err)
	}

	return res.Listeners, nil
}


func (c *Client) getHTTPRouteConfig(name string) (*HTTPRouteConfig, error) {
	url := fmt.Sprintf("%v/v1/routes/%v/%v/%v", c.endpoint, name, c.serviceCluster, c.serviceNode)
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("couldn't get routes: %v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if glog.V(3) {
		glog.Infof("Response from %v: %v", url, string(body))
	}

	var res HTTPRouteConfig

	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal response from pilot: %v", err)
	}

	return &res, nil
}

func (c *Client) getClusters() (Clusters, error) {
	url := fmt.Sprintf("%v/v1/clusters/%v/%v", c.endpoint, c.serviceCluster, c.serviceNode)
	resp, err := c.httpClient.Get(url)

	if err != nil {
		return nil, fmt.Errorf("couldn't get clusters: %v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if glog.V(3) {
		glog.Infof("Response from %v: %v", url, string(body))
	}

	res := ClusterManager{}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, fmt.Errorf("couldn't unmarshal response from pilot: %v", err)
	}

	return res.Clusters, nil
}

func (c *Client) getHostsForService(service string) (ServiceHosts, error) {
	url := fmt.Sprintf("%v/v1/registration/%v", c.endpoint, service)
	resp, err := c.httpClient.Get(url)

	res := ServiceHosts{}

	if err != nil {
		return res, fmt.Errorf("couldn't get service hosts: %v", err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if glog.V(3) {
		glog.Infof("Response from %v: %v", url, string(body))
	}

	err = json.Unmarshal(body, &res)
	if err != nil {
		return res, fmt.Errorf("couldn't unmarshal response from pilot: %v", err)
	}

	return res, nil
}

func finishUnmarshallingListeners(listeners Listeners) error {
	for _, l := range listeners {
		if len(l.Filters) == 0 {
			continue
		}
		f := l.Filters[0]
		if f.Type != "read" {
			continue
		}

		if f.Name == "http_connection_manager" {
			var httpConfig HTTPFilterConfig
			err := json.Unmarshal(f.Config, &httpConfig)
			if err != nil {
				return fmt.Errorf("couldn't unmarshal HTTPFilterConfig: %v", err)
			}

			for i := range httpConfig.Filters {
				if httpConfig.Filters[i].Type == "decoder" && httpConfig.Filters[i].Name == "mixer" {
					var mixerConfig FilterMixerV2Config
					err = json.Unmarshal(httpConfig.Filters[i].Config, &mixerConfig)
					if err != nil {
						return fmt.Errorf("couldn't unmarshal FilterMixerConfig: %v", err)
					}

					dec := json.NewDecoder(strings.NewReader(string(*mixerConfig.V2.Config)))

					err := dec.Decode(&mixerConfig.V2.ServiceConfig)
					if err != nil {
						return fmt.Errorf("couldn't unmarshal ServiceConfigs: %v", err)
					}

					httpConfig.Filters[i].FilterMixerConfig = mixerConfig.V2

				} else if httpConfig.Filters[i].Type == "decoder" && httpConfig.Filters[i].Name == "fault" {
					var faultConfig FilterFaultConfig
					err = json.Unmarshal(httpConfig.Filters[i].Config, &faultConfig)
					if err != nil {
						return fmt.Errorf("couldn't unmarshal FilterFaultConfig: %v", err)
					}
					httpConfig.Filters[i].FilterFaultConfig = &faultConfig
				}
			}

			f.HTTPFilterConfig = &httpConfig

		} else if f.Name == "tcp_proxy" {
			var tcpConfig TCPProxyFilterConfig
			err := json.Unmarshal(f.Config, &tcpConfig)
			if err != nil {
				return fmt.Errorf("couldn't unmarshal TCPProxyFilterConfig: %v", err)
			}
			f.TCPProxyFilterConfig = &tcpConfig
		}
	}

	return nil
}

// GetConfig returns configuration for a sidecar proxy from Pilot.
func (c *Client) GetConfig() ProxyConfig {
	cfg := ProxyConfig{
		HTTPListeners:    make(map[string]Listener),
		TCPListeners:     make(map[string]Listener),
		HTTPRouteConfigs: make(map[string]HTTPRouteConfig),
		Services:         make(map[string]ServiceHosts),
		Clusters:         make(map[string]Cluster),
	}

	listeners, err := c.getListeners()
	if err != nil {
		glog.Fatalf("Error getting listeners: %v", err)
	}

	for _, l := range listeners {
		if l.BindToPort == false {
			if len(l.Filters) == 0 {
				continue
			}
			f := l.Filters[0]

			if f.Type == "read" && f.Name == "http_connection_manager" {
				if f.HTTPFilterConfig != nil {
					if f.HTTPFilterConfig.RDS != nil {
						name := f.HTTPFilterConfig.RDS.RouteConfigName
						rc, err := c.getHTTPRouteConfig(name)
						if err != nil {
							glog.Warningf("Error getting a route config for listener %v: %v, skipping", l.Name, err)
							continue
						}
						cfg.HTTPRouteConfigs[name] = *rc
					} else if f.HTTPFilterConfig.RouteConfig == nil {
						glog.Warningf("Got a listener %v without  RDS or RouteConfig, skipping:", l.Name)
						continue
					}
					cfg.HTTPListeners[l.Name] = *l
				}
			} else {
				glog.Warningf("Got a listener %v with unknown filter: %v - %v, Skipping. For now, only http traffic is supported", l.Name, f.Type, f.Name)
				continue
			}
		}
	}

	clusters, err := c.getClusters()
	if err != nil {
		glog.Fatalf("Couldn't get clusters from Pilot: %v", err)
	}

	podServiceSet := make(map[string]bool)

	for _, cluster := range clusters {
		if cluster.ServiceName == "" {
			if cluster.Hosts != nil {
				var hosts []*ServiceHost
				for _, h := range cluster.Hosts {
					ip, port := parseDestination(h.URL)
					intPort, _ := strconv.Atoi(port)
					sh := ServiceHost{
						Address: ip,
						Port:    intPort,
					}
					hosts = append(hosts, &sh)
				}
				cfg.Services[cluster.Name] = ServiceHosts{Hosts: hosts}
				cfg.Clusters[cluster.Name] = *cluster
			}
			continue
		}
		hosts, err := c.getHostsForService(cluster.ServiceName)
		if err != nil {
			glog.Fatalf("Couldn't get hosts of service %v from Pilot: %v", cluster.ServiceName, err)
		}
		cfg.Services[cluster.Name] = hosts
		for _, h := range hosts.Hosts {
			if h.Address == c.podIP {
				podServiceSet[cluster.ServiceName] = true
			}
		}
		cfg.Clusters[cluster.Name] = *cluster
	}

	var podServices []string
	for svc := range podServiceSet {
		parts := strings.SplitN(svc, "|", 2)
		podServices = append(podServices, parts[0])
	}

	cfg.TargetService = strings.Join(podServices, ",")

	return cfg
}

func parseDestination(dest string) (string, string) {
	ipPort := strings.TrimPrefix(dest, "tcp://")
	parts := strings.Split(ipPort, ":")
	return parts[0], parts[1]
}
