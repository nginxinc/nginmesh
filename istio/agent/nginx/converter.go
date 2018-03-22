// Package nginx converts Envoy load balancing configuration from Pilot to NGINX configuration.
// Additionally, the package provides means for starting/stoping NGINX and applying a new configuration.
package nginx

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/nginmesh/nginmesh/istio/agent/pilot"
	"github.com/golang/glog"
)

// ConfigVariables holds global variables used in NGINX configuration.
type ConfigVariables struct {
	BindAddress        string
	ServiceNode        string
	ServiceCluster     string
	DisableMixerReport bool
	DisableMixerCheck  bool
	DisableTracing     bool
}

// Converter converts load balancing configuration from Pilot to NGINX configuration.
type Converter struct {
	configVars *ConfigVariables
}

// NewConverter creates a new converter.
func NewConverter(configVars *ConfigVariables) *Converter {
	return &Converter{
		configVars: configVars,
	}
}

// Convert converts load balancing configuration from Pilot to NGINX configuration.
func (conv *Converter) Convert(proxyConfig pilot.ProxyConfig) Config {
	httpConfigs, destMaps := conv.convertHTTPListeners(proxyConfig)
	tcpConfigs, tcpDestMaps := conv.convertTCPListeners(proxyConfig)
	destMaps = append(destMaps, tcpDestMaps...)
	mixer := buildMixerConfig(proxyConfig, httpConfigs)

	return Config{
		HTTPConfigs: httpConfigs,
		TCPConfigs:  tcpConfigs,
		Main: Main{
			Mixer:           mixer,
			DestinationMaps: destMaps,
			PodIP:           conv.configVars.BindAddress,
			ServiceNode:     conv.configVars.ServiceNode,
			ServiceCluster:  conv.configVars.ServiceCluster,
			Tracing:         !conv.configVars.DisableTracing,
		},
	}
}

func buildMixerConfig(proxyConfig pilot.ProxyConfig, httpConfigs []HTTPConfig) *MainMixer {
	var mixer *MainMixer

	if mixerCluster, exists := proxyConfig.Clusters["mixer_server"]; exists {
		if len(mixerCluster.Hosts) > 0 {
			addr, port := parseDestination(mixerCluster.Hosts[0].URL)
			for _, cfg := range httpConfigs {
				if cfg.Mixer != nil {
					mixer = &MainMixer{
						MixerServer: addr,
						MixerPort:   port,
					}
					break
				}
			}
		}
	}

	return mixer
}

func (conv *Converter) convertHTTPListeners(proxyConfig pilot.ProxyConfig) ([]HTTPConfig, []DestinationMap) {
	var httpConfigs []HTTPConfig

	var keys []string
	for k := range proxyConfig.HTTPListeners {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	port := 20000
	var destMaps []DestinationMap

	for _, name := range keys {
		l := proxyConfig.HTTPListeners[name]

		localAddress := "127.0.0.1"
		localPort := strconv.Itoa(port)

		cfg := conv.convertHTTPListener(l, proxyConfig, localAddress, localPort)
		httpConfigs = append(httpConfigs, cfg)

		destIP, destPort := parseDestination(l.Address)
		dest := fmt.Sprintf("'~*:%v$'", destPort)
		if l.SSLContext != nil {
			dest = fmt.Sprintf("%v:%v", destIP, destPort)
		}

		dm := DestinationMap{
			Remote: dest,
			Local:  fmt.Sprintf("%v:%v", localAddress, localPort),
		}
		destMaps = append(destMaps, dm)

		port++
	}

	return httpConfigs, destMaps
}

func (conv *Converter) convertTCPListeners(proxyConfig pilot.ProxyConfig) ([]TCPConfig, []DestinationMap) {
	var tcpConfigs []TCPConfig
	var tcpKeys []string
	for k := range proxyConfig.TCPListeners {
		tcpKeys = append(tcpKeys, k)
	}
	sort.Strings(tcpKeys)

	port := 10000 // start with 10000

	var destMaps []DestinationMap

	for _, name := range tcpKeys {
		l := proxyConfig.TCPListeners[name]
		localAddress := "127.0.0.1"
		localPort := strconv.Itoa(port)

		tcpCfg, dest := conv.convertTCPListener(l, proxyConfig, localAddress, localPort)
		tcpConfigs = append(tcpConfigs, tcpCfg)

		dm := DestinationMap{
			Remote: dest,
			Local:  fmt.Sprintf("%v:%v", localAddress, localPort),
		}
		destMaps = append(destMaps, dm)

		port++
	}
	return tcpConfigs, destMaps
}

func (conv *Converter) convertHTTPListener(listener pilot.Listener, proxyConfig pilot.ProxyConfig, localAddress string, localPort string) HTTPConfig {
	var servers []VirtualServer
	var upstreams []Upstream
	var maps []Map

	f := listener.Filters[0]

	var rc pilot.HTTPRouteConfig
	if f.HTTPFilterConfig.RDS != nil {
		rc = proxyConfig.HTTPRouteConfigs[f.HTTPFilterConfig.RDS.RouteConfigName]
	} else {
		rc = *f.HTTPFilterConfig.RouteConfig
	}

	faults := make(map[string]*pilot.FilterFaultConfig)
	for _, cf := range f.HTTPFilterConfig.Filters {
		if cf.FilterFaultConfig != nil {
			faults[cf.FilterFaultConfig.UpstreamCluster] = cf.FilterFaultConfig
		}
	}

	for i, vh := range rc.VirtualHosts {
		isTargetService := false

		for _, domain := range vh.Domains {
			if domain == proxyConfig.TargetService {
				isTargetService = true
			}
		}

		var splits []SplitClient
		var locs []Location
		var ups []Upstream
		var ssl *VirtualServerSSL

		names := vh.Domains
		if len(vh.Domains) == 1 && vh.Domains[0] == "*" {
			names = []string{"_"}
		}

		if listener.SSLContext != nil {
			ssl = &VirtualServerSSL{
				Certificate:        listener.SSLContext.CertChainFile,
				Key:                listener.SSLContext.PrivateKeyFile,
				TrustedCertificate: listener.SSLContext.CaCertFile,
				// listener.SSLContext.RequireClientCertificate
			}
		}

		varSuffix := fmt.Sprintf("%v_%v", localPort, i)

		if isTargetService {
			loc, u := conv.createLocalLocationAndUpstream(vh, proxyConfig.Clusters, proxyConfig.Services)
			locs = append(locs, loc)
			ups = append(ups, u)

		} else {
			var vhmaps []Map
			locs, vhmaps, splits = conv.createLocationsAndMapsAndSplitClients(vh, proxyConfig.Clusters, varSuffix, faults)
			maps = append(maps, vhmaps...)

			ups = createUpstreams(vh, proxyConfig.Clusters, proxyConfig.Services, conv.configVars.BindAddress)
		}

		s := VirtualServer{
			Address:      localAddress,
			Port:         localPort,
			Names:        names,
			Locations:    locs,
			IsTarget:     isTargetService,
			SplitClients: splits,
			SSL:          ssl,
		}

		servers = append(servers, s)
		upstreams = append(upstreams, ups...)
	}

	var mixer HTTPMixer
	for _, cf := range f.HTTPFilterConfig.Filters {
		glog.Info("processing listner: %s",listener.Address)
		if cf.FilterMixerConfig != nil {

			filterMixerConfig := cf.FilterMixerConfig

			var sourceIp string
			var sourceUid string
			var destinationIp string
			var destinationUid string
			var destinationService string

			if filterMixerConfig.ForwardAttributes != nil {
				forwardAttributes := cf.FilterMixerConfig.ForwardAttributes.Attributes

				if forwardAttributes.SourceIp != nil {
					sourceIp = forwardAttributes.SourceIp.BytesValue
					//glog.Info("detected sourceIp: %s",sourceIp)
				}
	
				if forwardAttributes.SourceUid != nil {
					sourceUid = forwardAttributes.SourceUid.StringValue
					//glog.Info("detected source Uid: %s",sourceUid)
				}
			}

			if filterMixerConfig.MixerAttributes != nil {
				// glog.Info("has mixer attributes")
				mixerAttributes := cf.FilterMixerConfig.MixerAttributes.Attributes
				
				if mixerAttributes.DestinationIp != nil {
					sourceIp = mixerAttributes.DestinationIp.BytesValue
				}
				
				if mixerAttributes.DestinationUid != nil {
					destinationUid = mixerAttributes.DestinationUid.StringValue
				}

				
				if mixerAttributes.DestinationService != nil {
					destinationService = mixerAttributes.DestinationService.StringValue
				}
			}
			

			

			mixer = HTTPMixer{
				SourceIP:          sourceIp,
				SourceUID:         sourceUid,
				DestinationIP:     destinationIp,
				DestinationUID:    destinationUid,
				DestinationService: destinationService,
				QuotaName:          cf.FilterMixerConfig.QuotaName,
			}
			break
		}
	}

	return HTTPConfig{
		Name:           listener.Name,
		Upstreams:      upstreams,
		VirtualServers: servers,
		Mixer:          &mixer,
		Maps:           maps,
	}
}

func (conv *Converter) convertTCPListener(listener pilot.Listener, proxyConfig pilot.ProxyConfig, localAddress string, localPort string) (TCPConfig, string) {
	ip, port := parseDestination(listener.Address)

	r := listener.Filters[0].TCPProxyFilterConfig.RouteConfig.Routes[0]
	var upsServers []Server

	for _, h := range proxyConfig.Services[r.Cluster].Hosts {
		s := Server{
			Address: h.Address,
			Port:    strconv.Itoa(h.Port),
		}
		upsServers = append(upsServers, s)
	}

	if len(upsServers) == 0 {
		upsServers = append(upsServers, Server{
			Address: "127.0.0.127",
			Port:    "8181",
		})
	}

	cluster := proxyConfig.Clusters[r.Cluster]
	method := getLBMethod(&cluster)

	upstreams := []Upstream{
		Upstream{
			Name:     r.Cluster,
			Servers:  upsServers,
			LBMethod: method,
		},
	}

	server := TCPServer{
		Address:        localAddress,
		Port:           localPort,
		Upstream:       r.Cluster,
		ConnectTimeout: cluster.ConnectTimeoutMs,
	}

	return TCPConfig{
		Name:      listener.Name,
		Servers:   []TCPServer{server},
		Upstreams: upstreams,
	}, fmt.Sprintf("%v:%v", ip, port)
}

func parseDestination(dest string) (string, string) {
	ipPort := strings.TrimPrefix(dest, "tcp://")
	parts := strings.Split(ipPort, ":")
	return parts[0], parts[1]
}

func (conv *Converter) createLocalLocationAndUpstream(vh *pilot.VirtualHost, clusters map[string]pilot.Cluster, services map[string]pilot.ServiceHosts) (Location, Upstream) {
	port := 0
	name := ""
	var rp *pilot.RetryPolicy

routes:
	for _, r := range vh.Routes {
		if multiDest(r) {
			for _, c := range r.WeightedClusters.Clusters {
				hosts, exists := services[c.Name]
				if exists {
					for _, h := range hosts.Hosts {
						if h.Address == conv.configVars.BindAddress {
							port = h.Port
							name = c.Name
							rp = r.RetryPolicy
							break routes
						}
					}
				}
			}
		} else {
			hosts, exists := services[r.Cluster]
			if exists {
				for _, h := range hosts.Hosts {
					if h.Address == conv.configVars.BindAddress {
						port = h.Port
						name = r.Cluster
						rp = r.RetryPolicy
						break routes
					}
				}
			}
		}
	}

	loc := Location{
		Internal:          false,
		Path:              "/",
		Upstream:          name,
		ConnectTimeout:    clusters[name].ConnectTimeoutMs,
		ProxyNextUpstream: createProxyNextUpstream(rp),
		MixerCheck:        !conv.configVars.DisableMixerCheck,
		MixerReport:       !conv.configVars.DisableMixerReport,
		Tracing:           !conv.configVars.DisableTracing,
	}
	upsServers := []Server{
		{"127.0.0.1", strconv.Itoa(port), 0, 0, 0},
	}

	ups := Upstream{
		Name:    name,
		Servers: upsServers,
	}

	return loc, ups
}

func createProxyNextUpstream(rp *pilot.RetryPolicy) *ProxyNextUpstream {
	if rp != nil {
		return &ProxyNextUpstream{
			Condition: "error http_500 http_502 http_503 http_504", // 5xx,connect-failure,refused-stream
			Timeout:   rp.PerTryTimeoutMS,
			Tries:     rp.NumRetries,
		}
	}
	return nil
}

func createExpression(varSuffix string, index int, headers []pilot.Header) []Map {
	var maps []Map

	for i, h := range headers {
		source := fmt.Sprintf("$http_%v", h.Name)
		if h.Name == "$uri" { // special case for redirect rules
			source = "$uri"
		}

		variable := fmt.Sprintf("$res_%v_%v_%v", varSuffix, index, i)
		if i == 0 {
			variable = fmt.Sprintf("$res_%v_%v", varSuffix, index)
		}

		res := fmt.Sprintf("$res_%v_%v_%v", varSuffix, index, i+1)
		if i == len(headers)-1 {
			res = "1"
		}

		value := h.Value
		if h.Regex {
			value = "~" + value
		}

		m := Map{
			Source:   source,
			Variable: variable,
			Params:   map[string]string{fmt.Sprintf("%v", value): res, "default": "0"},
		}

		maps = append(maps, m)
	}

	return maps
}

func createResultMap(varSuffix string, routes []*pilot.HTTPRoute) Map {
	var source string

	params := make(map[string]string)

	for i, r := range routes {
		val := strings.Repeat("0", i) + "1"
		if i != len(routes)-1 {
			val = "~^" + val
		}

		res := fmt.Sprintf("@loc_%v", i)
		if multiDest(r) {
			res = fmt.Sprintf("$loc_%v_%v", varSuffix, i)
		}

		params[val] = res
	}

	for i := 0; i < len(routes); i++ {
		source = fmt.Sprintf("%v$res_%v_%v", source, varSuffix, i)
	}

	res := Map{
		Source:   source,
		Variable: fmt.Sprintf("$loc_%v", varSuffix),
		Params:   params,
	}

	return res
}

func (conv *Converter) createLocationsAndMapsAndSplitClients(vh *pilot.VirtualHost, clusters map[string]pilot.Cluster, varSuffix string, faults map[string]*pilot.FilterFaultConfig) ([]Location, []Map, []SplitClient) {
	var locs []Location
	var maps []Map

	defaultLoc := Location{
		Internal: false,
		Path:     "/",
		Upstream: "localhost:8181",
		Expressions: []*Expression{
			&Expression{
				Result:    fmt.Sprintf("error_page 418 = $loc_%v; return 418", varSuffix),
				Condition: fmt.Sprintf("$loc_%v != ''", varSuffix),
			},
		},
	}

	var splits []SplitClient

	for i, r := range vh.Routes {
		if multiDest(r) {
			split := createSplitClient(varSuffix, i, r)
			splits = append(splits, split)

			for ci, c := range r.WeightedClusters.Clusters {
				cluster := clusters[c.Name]
				var exprs []*Expression

				loc := conv.createLocation(&cluster, fmt.Sprintf("%v_%v", i, ci))

				if f, exists := faults[c.Name]; exists {
					if f.Abort != nil {
						varName := fmt.Sprintf("$ups_fault_%v_%v_%v", varSuffix, i, ci)
						split, expr := createAbort(f, varName)
						splits = append(splits, split)
						exprs = append(exprs, &expr)
					}
				}

				loc.Expressions = exprs
				loc.ProxyNextUpstream = createProxyNextUpstream(r.RetryPolicy)

				if r.PrefixRewrite != "" {
					loc.Host = r.HostRewrite
					loc.Rewrite = &Rewrite{
						Prefix:      r.Prefix,
						Replacement: r.PrefixRewrite,
					}
				}

				locs = append(locs, loc)
			}
			// add a match based on the URI
			if r.PrefixRewrite != "" {
				r.Headers = append(r.Headers, pilot.Header{
					Name:  "$uri",
					Regex: true,
					Value: "^" + r.Prefix,
				})
			}
		} else {

			var loc Location

			if r.Cluster != "" {

				cluster := clusters[r.Cluster]
				loc = conv.createLocation(&cluster, fmt.Sprintf("%v", i))

				var exprs []*Expression
				if f, exists := faults[r.Cluster]; exists {
					if f.Abort != nil {
						varName := fmt.Sprintf("$ups_fault_%v_%v", varSuffix, i)
						split, expr := createAbort(f, varName)
						splits = append(splits, split)
						exprs = append(exprs, &expr)
					}
				}
				loc.Expressions = exprs
				loc.ProxyNextUpstream = createProxyNextUpstream(r.RetryPolicy)

				if r.PrefixRewrite != "" {
					loc.Host = r.HostRewrite
					loc.Rewrite = &Rewrite{
						Prefix:      r.Prefix,
						Replacement: r.PrefixRewrite,
					}
					// add a match based on the URI
					r.Headers = append(r.Headers, pilot.Header{
						Name:  "$uri",
						Regex: true,
						Value: "^" + r.Prefix,
					})
				}
			} else {
				if r.PathRedirect != "" {
					loc = Location{
						Internal: true,
						Redirect: &Redirect{
							Code: 302,
							URL:  r.HostRedirect + r.PathRedirect,
						},
						Path:        fmt.Sprintf("@loc_%v", i),
						MixerCheck:  !conv.configVars.DisableMixerCheck,
						MixerReport: !conv.configVars.DisableMixerReport,
						Tracing:     !conv.configVars.DisableTracing,
					}
					// add a match based on the URI
					r.Headers = append(r.Headers, pilot.Header{
						Name:  "$uri",
						Regex: true,
						Value: "^" + r.Path,
					})
				}
				// TO-DO other cases?
			}

			locs = append(locs, loc)
		}

		if len(r.Headers) > 0 {
			resMaps := createExpression(varSuffix, i, r.Headers)
			maps = append(maps, resMaps...)
		} else {
			defaultLoc.Sets = append(defaultLoc.Sets, Set{fmt.Sprintf("$res_%v_%v", varSuffix, i), "1"})
		}
	}

	m := createResultMap(varSuffix, vh.Routes)
	maps = append(maps, m)

	locs = append(locs, defaultLoc)

	return locs, maps, splits
}

func (conv *Converter) createLocation(cluster *pilot.Cluster, locSuffix string) Location {
	var ssl *LocationSSL
	if cluster.SSLContext != nil {
		var name string
		if len(cluster.SSLContext.VerifySubjectAltName) > 0 {
			name = cluster.SSLContext.VerifySubjectAltName[0]
		}
		ssl = &LocationSSL{
			Certificate:        cluster.SSLContext.CertChainFile,
			Key:                cluster.SSLContext.PrivateKeyFile,
			TrustedCertificate: cluster.SSLContext.CaCertFile,
			Name:               name,
		}
	}

	return Location{
		Internal:       true,
		Upstream:       cluster.Name,
		Path:           fmt.Sprintf("@loc_%v", locSuffix),
		SSL:            ssl,
		ConnectTimeout: cluster.ConnectTimeoutMs,
		MixerCheck:     !conv.configVars.DisableMixerCheck,
		MixerReport:    !conv.configVars.DisableMixerReport,
		Tracing:        !conv.configVars.DisableTracing,
	}
}

func createUpstreams(vh *pilot.VirtualHost, clusters map[string]pilot.Cluster, services map[string]pilot.ServiceHosts, bindAddress string) []Upstream {
	upstreams := make(map[string]Upstream)

	for _, r := range vh.Routes {
		if multiDest(r) {
			for _, c := range r.WeightedClusters.Clusters {
				cluster := clusters[c.Name]
				ups := clusterToUpstream(&cluster, services[c.Name], bindAddress)
				if _, exists := upstreams[ups.Name]; !exists {
					upstreams[ups.Name] = ups
				}
			}
		} else if r.Cluster != "" {
			cluster := clusters[r.Cluster]
			ups := clusterToUpstream(&cluster, services[r.Cluster], bindAddress)
			if _, exists := upstreams[ups.Name]; !exists {
				upstreams[ups.Name] = ups
			}
		}
	}

	var keys []string
	for k := range upstreams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var res []Upstream

	for _, k := range keys {
		res = append(res, upstreams[k])
	}

	return res
}

func multiDest(r *pilot.HTTPRoute) bool {
	return r.WeightedClusters != nil && len(r.WeightedClusters.Clusters) > 0
}

func createSplitClient(varSuffix string, index int, r *pilot.HTTPRoute) SplitClient {
	var dists []Distribution
	for i, c := range r.WeightedClusters.Clusters {
		d := Distribution{
			Weight: c.Weight,
			Value:  fmt.Sprintf("@loc_%v_%v", index, i),
		}
		dists = append(dists, d)
	}

	varName := fmt.Sprintf("$loc_%v_%v", varSuffix, index)

	split := SplitClient{
		Variable:      varName,
		Distributions: dists,
	}
	return split
}

func createAbort(f *pilot.FilterFaultConfig, varName string) (SplitClient, Expression) {
	expr := Expression{
		Condition: fmt.Sprintf("%v != ''", varName),
		Result:    fmt.Sprintf("return %v", f.Abort.HTTPStatus),
	}
	split := SplitClient{
		Variable: varName,
		Distributions: []Distribution{
			Distribution{
				Weight: f.Abort.Percent,
				Value:  "abort",
			},
			Distribution{
				Weight: 100 - f.Abort.Percent,
				Value:  "''",
			},
		},
	}

	return split, expr
}

func clusterToUpstream(cluster *pilot.Cluster, hosts pilot.ServiceHosts, bindAddress string) Upstream {
	var upsServers []Server

	local := false
	port := 0
	if len(hosts.Hosts) > 0 {
		maxConns := 1024 / len(hosts.Hosts)
		if cluster.CircuitBreaker != nil {
			maxConns = cluster.CircuitBreaker.Default.MaxConnections / len(hosts.Hosts)
		}
		maxFails := 5
		var failTimeout int64
		if cluster.OutlierDetection != nil {
			maxFails = cluster.OutlierDetection.ConsecutiveErrors
			failTimeout = cluster.OutlierDetection.IntervalMS / 1000
		}

		for _, h := range hosts.Hosts {
			if h.Address == bindAddress || h.Address == "127.0.0.1" {
				local = true
				port = h.Port
				break
			}

			upsServers = append(upsServers, Server{
				Address:     h.Address,
				Port:        strconv.Itoa(h.Port),
				MaxConns:    maxConns,
				MaxFails:    maxFails,
				FailTimeout: failTimeout,
			})
		}
	}

	if local {
		upsServers = []Server{
			Server{
				Address: "127.0.0.1",
				Port:    strconv.Itoa(port),
			},
		}
	}

	if len(upsServers) == 0 {
		upsServers = append(upsServers, Server{
			Address: "127.0.0.127",
			Port:    "8181",
		})
	}

	method := getLBMethod(cluster)

	ups := Upstream{
		Name:     cluster.Name,
		Servers:  upsServers,
		LBMethod: method,
	}

	return ups
}

func getLBMethod(cluster *pilot.Cluster) string {
	method := LBMethodRR

	switch cluster.LbType {
	case pilot.LbTypeRandom:
		method = LBMethodRandom
	case pilot.LbTypeLeastRequest:
		method = LBMethodLeastConn
	}

	return method
}
