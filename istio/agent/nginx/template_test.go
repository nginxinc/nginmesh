package nginx

import (
	"bytes"
	"testing"
	"text/template"
)

func TestTemplate(t *testing.T) {
	cfg := HTTPConfig{
		Name: "8080",
		Upstreams: []Upstream{
			Upstream{
				Name: "details",
				Servers: []Server{
					Server{
						Address:     "192.168.1.2",
						Port:        "8080",
						MaxConns:    15,
						MaxFails:    5,
						FailTimeout: 20000,
					},
				},
			},
			Upstream{
				Name: "details2",
				Servers: []Server{
					Server{
						Address: "192.168.1.2",
						Port:    "8080",
					},
				},
				LBMethod: LBMethodLeastConn,
			},
		},
		VirtualServers: []VirtualServer{
			VirtualServer{
				Address: "127.0.0.1",
				Port:    "2001",
				Names:   []string{"details", "details.local"},
				SSL: &VirtualServerSSL{
					Certificate:        "cert.pem",
					Key:                "key.pem",
					TrustedCertificate: "ca.pem",
				},
				SplitClients: []SplitClient{
					SplitClient{
						Variable: "$ups_from_split_clients_suffix_1",
						Distributions: []Distribution{
							Distribution{
								Weight: 50,
								Value:  "details",
							},
							Distribution{
								Weight: 50,
								Value:  "details2",
							},
						},
					},
				},
				Locations: []Location{
					Location{
						Path:           "@loc_0",
						Internal:       true,
						Upstream:       "cluster-1",
						ConnectTimeout: 1000,
						ProxyNextUpstream: &ProxyNextUpstream{
							Condition: "error timeout",
							Timeout:   1000,
							Tries:     5,
						},
						MixerCheck:  true,
						MixerReport: true,
						Tracing:     true,
					},
					Location{
						Path:           "@loc_1",
						Internal:       true,
						Upstream:       "cluster-2",
						ConnectTimeout: 1500,
					},
					Location{
						Path:     "@loc_2",
						Internal: true,
						Redirect: &Redirect{
							Code: 302,
							URL:  "example.com/foo",
						},
					},
					Location{
						Path:     "@loc_3",
						Internal: true,
						Host:     "details",
						Rewrite: &Rewrite{
							Prefix:      "/details",
							Replacement: "/moreinfo",
						},
						Upstream: "cluster-3",
					},
					Location{
						Path:     "@loc_4",
						Internal: true,
						Upstream: "$ups_from_split_clients_suffix_1",
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
								Variable: "$res_suffix_1",
								Value:    "@loc_1",
							},
						},
					},
				},
			},
		},
		Mixer: &HTTPMixer{
			DestinationIP:      "10.24.0.62",
			DestinationService: "productpage.default.svc.cluster.local",
			DestinationUID:     "kubernetes://productpage",
			SourceIP:           "10.24.0.62",
			SourceUID:          "kubernetes://productpage",
		},
		Maps: []Map{
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
	}

	tmpl, err := template.New("config.tmpl").Parse(httpTemplate)
	if err != nil {
		t.Errorf("Couldn't parse the template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, cfg)
	t.Log(string(buf.Bytes()))
	if err != nil {
		t.Errorf("Couldn't execute the template: %v", err)
	}
}

func TestTemplateTCP(t *testing.T) {
	cfg := TCPConfig{
		Name: "192.168.1.3-8888",
		Upstreams: []Upstream{
			Upstream{
				Name: "backend",
				Servers: []Server{
					Server{
						Address: "192.168.10.10",
						Port:    "7777",
					},
				},
			},
			Upstream{
				Name: "backend-2",
				Servers: []Server{
					Server{
						Address: "192.168.10.10",
						Port:    "7777",
					},
				},
				LBMethod: LBMethodRandom,
			},
		},
		Servers: []TCPServer{
			TCPServer{
				Address:        "192.168.1.3",
				Port:           "8888",
				Upstream:       "backend",
				ConnectTimeout: 1000,
			},
		},
	}
	tmpl, err := template.New("tcpconfig.tmpl").Parse(tcpTemplate)
	if err != nil {
		t.Errorf("Couldn't parse the template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, cfg)
	t.Log(string(buf.Bytes()))
	if err != nil {
		t.Errorf("Couldn't execute the template: %v", err)
	}
}

func TestMainTemplate(t *testing.T) {
	mainCfg := Main{
		DestinationMaps: []DestinationMap{
			{"192.168.1.2:1122", "127.0.0.1:1000"},
			{"192.168.1.3:3366", "127.0.0.1:1001"},
		},
		Mixer: &MainMixer{
			MixerServer: "istio-mixer.istio-system",
			MixerPort:   "9091",
		},
		PodIP:          "10.24.0.62",
		ServiceNode:    "sidecar~10.20.2.52~productpage-v1-1262296934-bmrl8.default~default.svc.cluster.local",
		ServiceCluster: "productpage",
		Tracing:        true,
	}

	tmpl, err := template.New("main.tmpl").Parse(mainTemplate)
	if err != nil {
		t.Errorf("Couldn't parse the template: %v", err)
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, mainCfg)
	t.Log(string(buf.Bytes()))
	if err != nil {
		t.Errorf("Couldn't execute the template: %v", err)
	}
}
