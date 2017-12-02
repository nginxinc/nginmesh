package nginx

import (
	"reflect"
	"testing"

	"github.com/nginmesh/nginmesh/istio/agent/pilot"
)

var vh = pilot.VirtualHost{
	Name:    "first.example.com|http",
	Domains: []string{"first.example.com"},
	Routes: []*pilot.HTTPRoute{
		&pilot.HTTPRoute{
			Prefix:  "/",
			Cluster: "cluster-1",
			Headers: []pilot.Header{
				{"x", "a", false},
				{"x", "b", false},
			},
		},
		&pilot.HTTPRoute{
			Path:         "/old",
			PathRedirect: "/new",
			HostRedirect: "example.org",
		},
		&pilot.HTTPRoute{
			Prefix: "/",
			WeightedClusters: &pilot.WeightedCluster{
				Clusters: []*pilot.WeightedClusterEntry{
					&pilot.WeightedClusterEntry{"cluster-2-v1", 50},
					&pilot.WeightedClusterEntry{"cluster-2-v2", 50},
				},
			},
		},
	},
}

var services = map[string]pilot.ServiceHosts{
	"cluster-1": pilot.ServiceHosts{
		Hosts: []*pilot.ServiceHost{
			&pilot.ServiceHost{"192.168.0.1", 80, nil},
		},
	},
	"cluster-2-v1": pilot.ServiceHosts{
		Hosts: []*pilot.ServiceHost{
			&pilot.ServiceHost{"192.168.1.1", 8080, nil},
		},
	},
	"cluster-2-v2": pilot.ServiceHosts{
		Hosts: []*pilot.ServiceHost{
			&pilot.ServiceHost{"192.168.1.2", 8080, nil},
		},
	},
}

var clusters = map[string]pilot.Cluster{
	"cluster-1": pilot.Cluster{
		Name:   "cluster-1",
		LbType: pilot.LbTypeRoundRobin,
		CircuitBreaker: &pilot.CircuitBreaker{
			Default: pilot.DefaultCBPriority{
				MaxConnections: 2048,
			},
		},
		OutlierDetection: &pilot.OutlierDetection{
			ConsecutiveErrors: 10,
			IntervalMS:        5000,
		},
	},
	"cluster-2-v1": pilot.Cluster{
		Name:   "cluster-2-v1",
		LbType: pilot.LbTypeLeastRequest,
	},
	"cluster-2-v2": pilot.Cluster{
		Name:   "cluster-2-v2",
		LbType: pilot.LbTypeRandom,
	},
}

var expectedUpstreams = []Upstream{
	Upstream{
		Name: "cluster-1",
		Servers: []Server{
			Server{
				Address:     "192.168.0.1",
				Port:        "80",
				MaxConns:    2048,
				MaxFails:    10,
				FailTimeout: 5,
			},
		},
		LBMethod: LBMethodRR,
	},
	Upstream{
		Name: "cluster-2-v1",
		Servers: []Server{
			Server{
				Address:  "192.168.1.1",
				Port:     "8080",
				MaxConns: 1024,
				MaxFails: 5,
			},
		},
		LBMethod: LBMethodLeastConn,
	},
	Upstream{
		Name: "cluster-2-v2",
		Servers: []Server{
			Server{
				Address:  "192.168.1.2",
				Port:     "8080",
				MaxConns: 1024,
				MaxFails: 5,
			},
		},
		LBMethod: LBMethodRandom,
	},
}

var expectedSplits = []SplitClient{
	SplitClient{
		Variable: "$ups_from_split_clients_suffix_1",
		Distributions: []Distribution{
			Distribution{
				Weight: 50,
				Value:  "cluster-2-v1",
			},
			Distribution{
				Weight: 50,
				Value:  "cluster-2-v2",
			},
		},
	},
}

func TestCreateUpstreams(t *testing.T) {
	upstreams := createUpstreams(&vh, clusters, services, "127.0.0.1")

	if !reflect.DeepEqual(upstreams, expectedUpstreams) {
		t.Errorf("createUpstreamsAndSplitClients(...)\n got upstreams %v\nexpected upstreams %v", upstreams, expectedUpstreams)
	}
}

func TestCreateLocalLocationAndUpstream(t *testing.T) {
	expectedLoc := Location{
		Internal:    false,
		Path:        "/",
		Upstream:    "cluster-2-v2",
		MixerCheck:  true,
		MixerReport: true,
		Tracing:     true,
	}
	expectedUpstream := Upstream{
		Name: "cluster-2-v2",
		Servers: []Server{
			Server{
				Address: "127.0.0.1",
				Port:    "8080",
			},
		},
	}

	converter := NewConverter(&ConfigVariables{
		BindAddress: "192.168.1.2",
	})

	location, upstream := converter.createLocalLocationAndUpstream(&vh, clusters, services)

	if !reflect.DeepEqual(location, expectedLoc) {
		t.Errorf("createUpstreamsAndSplitClients(...); got location %v; expected location %v", location, expectedLoc)
	}
	if !reflect.DeepEqual(upstream, expectedUpstream) {
		t.Errorf("createUpstreamsAndSplitClients(...); got upstream %v; expected upstream %v", upstream, expectedUpstream)
	}

}

func TestCreateExpression(t *testing.T) {
	var tests = []struct {
		arg0     int
		arg1     []pilot.Header
		expected []Map
	}{
		{
			0,
			[]pilot.Header{
				pilot.Header{
					Name:  "x",
					Value: "1",
					Regex: false,
				},
				pilot.Header{
					Name:  "y",
					Value: "2",
					Regex: false,
				},
			},
			[]Map{
				Map{
					Source:   "$http_x",
					Variable: "$res_suffix_0",
					Params:   map[string]string{"1": "$res_suffix_0_1", "default": "0"},
				},
				Map{
					Source:   "$http_y",
					Variable: "$res_suffix_0_1",
					Params:   map[string]string{"2": "1", "default": "0"},
				},
			},
		},
		{
			1,
			[]pilot.Header{
				pilot.Header{
					Name:  "x",
					Value: "^(.*)$",
					Regex: true,
				},
			},
			[]Map{
				Map{
					Source:   "$http_x",
					Variable: "$res_suffix_1",
					Params:   map[string]string{"~^(.*)$": "1", "default": "0"},
				},
			},
		},
	}

	for _, test := range tests {
		result := createExpression("suffix", test.arg0, test.arg1)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("createExpression(...); got %v; expected %v", result, test.expected)
		}
	}
}

func TestCreateResultMap(t *testing.T) {
	expected := Map{
		Source:   "$res_suffix_0$res_suffix_1$res_suffix_2",
		Variable: "$loc_suffix",
		Params:   map[string]string{"~^1": "@loc_0", "~^01": "@loc_1", "001": "$loc_suffix_2"},
	}
	result := createResultMap("suffix", vh.Routes)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("createResultMap(...)\ngot %v\nexpected %v", result, expected)
	}
}

func TestCreateLocations(t *testing.T) {
	vhost := pilot.VirtualHost{
		Name:    "vhost0",
		Domains: []string{"example.com"},
		Routes: []*pilot.HTTPRoute{
			&pilot.HTTPRoute{
				Headers: []pilot.Header{
					pilot.Header{
						Name:  "x",
						Value: "value-x",
						Regex: false,
					},
				},
				Cluster: "cluster-1",
			},
			&pilot.HTTPRoute{
				Path:         "/old",
				PathRedirect: "/new",
				HostRedirect: "example.org",
			},
			&pilot.HTTPRoute{
				Prefix:        "/reviews",
				PrefixRewrite: "/details",
				Cluster:       "cluster-3",
				HostRewrite:   "details",
			},
			&pilot.HTTPRoute{
				Prefix:        "/reviews",
				PrefixRewrite: "/details",
				HostRewrite:   "details",
				WeightedClusters: &pilot.WeightedCluster{
					Clusters: []*pilot.WeightedClusterEntry{
						&pilot.WeightedClusterEntry{"cluster-4-1", 50},
						&pilot.WeightedClusterEntry{"cluster-4-2", 50},
					},
				},
			},
			&pilot.HTTPRoute{
				WeightedClusters: &pilot.WeightedCluster{
					Clusters: []*pilot.WeightedClusterEntry{
						&pilot.WeightedClusterEntry{"cluster-2-1", 50},
						&pilot.WeightedClusterEntry{"cluster-2-2", 50},
					},
				},
			},
		},
	}
	clusters := map[string]pilot.Cluster{
		"cluster-1": pilot.Cluster{
			Name: "cluster-1",
		},
		"cluster-2-1": pilot.Cluster{
			Name: "cluster-2-1",
		},
		"cluster-2-2": pilot.Cluster{
			Name: "cluster-2-2",
		},
		"cluster-3": pilot.Cluster{
			Name: "cluster-3",
		},
		"cluster-4-1": pilot.Cluster{
			Name: "cluster-4-1",
		},
		"cluster-4-2": pilot.Cluster{
			Name: "cluster-4-2",
		},
	}

	faults := map[string]*pilot.FilterFaultConfig{
		"cluster-1": &pilot.FilterFaultConfig{
			UpstreamCluster: "cluster-1",
			Abort: &pilot.AbortFilter{
				HTTPStatus: 400,
				Percent:    90,
			},
		},
		"cluster-2-1": &pilot.FilterFaultConfig{
			UpstreamCluster: "cluster-1",
			Abort: &pilot.AbortFilter{
				HTTPStatus: 400,
				Percent:    90,
			},
		}, "cluster-2-2": &pilot.FilterFaultConfig{
			UpstreamCluster: "cluster-1",
			Abort: &pilot.AbortFilter{
				HTTPStatus: 400,
				Percent:    90,
			},
		},
	}

	expected := []Location{
		Location{
			Path:     "@loc_0",
			Internal: true,
			Upstream: "cluster-1",
			Expressions: []*Expression{
				&Expression{
					Condition: "$ups_fault_suffix_0 != ''",
					Result:    "return 400",
				},
			},
			MixerCheck:  true,
			MixerReport: true,
			Tracing:     true,
		},
		Location{
			Path:     "@loc_1",
			Internal: true,
			Redirect: &Redirect{
				Code: 302,
				URL:  "example.org/new",
			},
			MixerCheck:  true,
			MixerReport: true,
			Tracing:     true,
		},
		Location{
			Path:     "@loc_2",
			Internal: true,
			Host:     "details",
			Upstream: "cluster-3",
			Rewrite: &Rewrite{
				Prefix:      "/reviews",
				Replacement: "/details",
			},
			MixerCheck:  true,
			MixerReport: true,
			Tracing:     true,
		},
		Location{
			Path:     "@loc_3_0",
			Internal: true,
			Upstream: "cluster-4-1",
			Host:     "details",
			Rewrite: &Rewrite{
				Prefix:      "/reviews",
				Replacement: "/details",
			},
			MixerCheck:  true,
			MixerReport: true,
			Tracing:     true,
		},
		Location{
			Path:     "@loc_3_1",
			Internal: true,
			Upstream: "cluster-4-2",
			Host:     "details",
			Rewrite: &Rewrite{
				Prefix:      "/reviews",
				Replacement: "/details",
			},
			MixerCheck:  true,
			MixerReport: true,
			Tracing:     true,
		},
		Location{
			Path:     "@loc_4_0",
			Internal: true,
			Upstream: "cluster-2-1",
			Expressions: []*Expression{
				&Expression{
					Condition: "$ups_fault_suffix_4_0 != ''",
					Result:    "return 400",
				},
			},
			MixerCheck:  true,
			MixerReport: true,
			Tracing:     true,
		},
		Location{
			Path:     "@loc_4_1",
			Internal: true,
			Upstream: "cluster-2-2",
			Expressions: []*Expression{
				&Expression{
					Condition: "$ups_fault_suffix_4_1 != ''",
					Result:    "return 400",
				},
			},
			MixerCheck:  true,
			MixerReport: true,
			Tracing:     true,
		},
		Location{
			Path:     "/",
			Upstream: "localhost:8181",
			Expressions: []*Expression{
				&Expression{
					Condition: "$loc_suffix != ''",
					Result:    "error_page 418 = $loc_suffix; return 418",
				},
			},
			Sets: []Set{
				Set{
					Variable: "$res_suffix_4",
					Value:    "1",
				},
			},
		},
	}

	expectedMaps := []Map{
		Map{
			Source:   "$http_x",
			Variable: "$res_suffix_0",
			Params: map[string]string{
				"value-x": "1",
				"default": "0",
			},
		},
		Map{
			Source:   "$uri",
			Variable: "$res_suffix_1",
			Params: map[string]string{
				"~^/old":  "1",
				"default": "0",
			},
		},
		Map{
			Source:   "$uri",
			Variable: "$res_suffix_2",
			Params: map[string]string{
				"~^/reviews": "1",
				"default":    "0",
			},
		},
		Map{
			Source:   "$uri",
			Variable: "$res_suffix_3",
			Params: map[string]string{
				"~^/reviews": "1",
				"default":    "0",
			},
		},
		Map{
			Source:   "$res_suffix_0$res_suffix_1$res_suffix_2$res_suffix_3$res_suffix_4",
			Variable: "$loc_suffix",
			Params: map[string]string{
				"~^1":    "@loc_0",
				"~^01":   "@loc_1",
				"~^001":  "@loc_2",
				"~^0001": "$loc_suffix_3",
				"00001":  "$loc_suffix_4",
			},
		},
	}

	expectedSplits := []SplitClient{
		SplitClient{
			Variable: "$ups_fault_suffix_0",
			Distributions: []Distribution{
				Distribution{
					Weight: 90,
					Value:  "abort",
				},
				Distribution{
					Weight: 10,
					Value:  "''",
				},
			},
		},
		SplitClient{
			Variable: "$loc_suffix_3",
			Distributions: []Distribution{
				Distribution{
					Weight: 50,
					Value:  "@loc_3_0",
				},
				Distribution{
					Weight: 50,
					Value:  "@loc_3_1",
				},
			},
		},
		SplitClient{
			Variable: "$loc_suffix_4",
			Distributions: []Distribution{
				Distribution{
					Weight: 50,
					Value:  "@loc_4_0",
				},
				Distribution{
					Weight: 50,
					Value:  "@loc_4_1",
				},
			},
		},
		SplitClient{
			Variable: "$ups_fault_suffix_4_0",
			Distributions: []Distribution{
				Distribution{
					Weight: 90,
					Value:  "abort",
				},
				Distribution{
					Weight: 10,
					Value:  "''",
				},
			},
		},
		SplitClient{
			Variable: "$ups_fault_suffix_4_1",
			Distributions: []Distribution{
				Distribution{
					Weight: 90,
					Value:  "abort",
				},
				Distribution{
					Weight: 10,
					Value:  "''",
				},
			},
		},
	}

	converter := NewConverter(&ConfigVariables{})

	result, maps, splits := converter.createLocationsAndMapsAndSplitClients(&vhost, clusters, "suffix", faults)

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("createLocationsAndMaps(...)\n got %#v\n expected %#v", result, expected)
	}

	if !reflect.DeepEqual(maps, expectedMaps) {
		t.Errorf("createLocationsAndMaps(...)\n got %#v\n expected %#v", maps, expectedMaps)
	}

	if !reflect.DeepEqual(splits, expectedSplits) {
		t.Errorf("createLocationsAndMaps(...)\n got %#v\n expected %#v", splits, expectedSplits)
	}
}
