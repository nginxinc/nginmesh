package pilot

import (
	"encoding/json"
	"testing"
)

var listeners = []byte(`{
	"listeners": [
	 {
	  "address": "tcp://0.0.0.0:15001",
	  "name": "virtual",
	  "filters": [],
	  "bind_to_port": true,
	  "use_original_dst": true
	 },
	 {
	  "address": "tcp://10.23.250.12:80",
	  "name": "tcp_10.23.250.12_80",
	  "filters": [
	   {
		"type": "read",
		"name": "tcp_proxy",
		"config": {
		 "stat_prefix": "tcp",
		 "route_config": {
		  "routes": [
		   {
			"cluster": "out.bffcf45bb7cf7eacfc01ecd666ddac27979567c5",
			"destination_ip_list": [
			 "10.23.250.12/32"
			]
		   }
		  ]
		 }
		}
	   }
	  ],
	  "bind_to_port": false
	 },
	 {
		"address": "tcp://0.0.0.0:8080",
		"name": "http_0.0.0.0_8080",
		"filters": [
		 {
		  "type": "read",
		  "name": "http_connection_manager",
		  "config": {
		   "codec_type": "auto",
		   "stat_prefix": "http",
		   "generate_request_id": true,
		   "tracing": {
			"operation_name": "ingress"
		   },
		   "rds": {
			"cluster": "rds",
			"route_config_name": "8080",
			"refresh_delay_ms": 30000
		   },
		   "filters": [
			{
			 "type": "decoder",
			 "name": "mixer",
			 "config": {
			  "mixer_attributes": {
			   "destination.ip": "10.20.2.13",
			   "destination.uid": "kubernetes://productpage"
			  },
			  "forward_attributes": {
			   "source.ip": "10.20.2.13",
			   "source.uid": "kubernetes://productpage"
			  },
			  "quota_name": "RequestCount"
			 }
			},
			{
				"type": "decoder",
				"name": "fault",
				"config": {
				 "abort": {
				  "abort_percent": 10,
				  "http_status": 400
				 },
				 "upstream_cluster": "out.8046187fa0ebd79882fa2d601888ca7b89aa7c37"
				}
			},
			{
				"type": "decoder",
				"name": "fault",
				"config": {
				 "abort": {
				  "abort_percent": 90,
				  "http_status": 400
				 },
				 "upstream_cluster": "out.8046187fa0ebd79882fa2d601888ca7b89aa7c37"
				}
			},
			{
			 "type": "decoder",
			 "name": "router",
			 "config": {}
			}
		   ],
		   "access_log": [
			{
			 "path": "/dev/stdout"
			}
		   ]
		  }
		 }
		],
		"bind_to_port": false
	   }
	]
   }`)

func TestUnmarshal(t *testing.T) {
	var res ldsResponse
	err := json.Unmarshal(listeners, &res)
	if err != nil {
		t.Error(err)
	}
	err = finishUnmarshallingListeners(res.Listeners)
	if err != nil {
		t.Error(err)
	}
	if res.Listeners[1].Filters[0].TCPProxyFilterConfig == nil {
		t.Error("TCPProxyFilterConfig is nil for the tcp filter")
	}
	if res.Listeners[2].Filters[0].HTTPFilterConfig == nil {
		t.Error("HTTPFilterConfig is nil for the http filter")
	}
	if filter := res.Listeners[2].Filters[0].HTTPFilterConfig.Filters[0]; filter.FilterMixerConfig == nil {
		t.Errorf("FilterMixerConfig is nil in %v", filter)
	}
	if filter := res.Listeners[2].Filters[0].HTTPFilterConfig.Filters[1]; filter.FilterFaultConfig == nil {
		t.Errorf("FilterMixerConfig is nil in %v", filter)
	}
}
