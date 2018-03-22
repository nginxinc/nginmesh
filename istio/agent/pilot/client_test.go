package pilot

import (
	"testing"
	"io/ioutil"
	"fmt"
	"os"
)


func TestUnmarshal(t *testing.T) {

	pwd, _ := os.Getwd()
	body, ferr := ioutil.ReadFile(pwd+"/test/0.3.0/listener.json")
	if ferr != nil {
		fmt.Printf("could not load file")
		t.Error(ferr)
		return
	}


	listeners, err := unMarshalListeners(body)
	if err != nil {
		t.Error(err)
	}

	if len(listeners) !=18 {
		t.Error("there should been 18 listener founded")
	}

	for _, listener := range listeners {
		if listener.Address == "tcp://10.40.1.19:9080" {
			
			filters := listener.Filters
			if len(filters) == 0  {
				t.Error("no filter founded")
				return
			}

			filter := filters[0]
			httpFilterConfig := filter.HTTPFilterConfig
			tcpFilterConfig := filter.TCPProxyFilterConfig

			if httpFilterConfig == nil  {
				t.Error("no http filter detected")
				return
			}

			if tcpFilterConfig != nil  {
				t.Error("wrong tcp filter")
			}

			httpFilters := httpFilterConfig.Filters

			if len(httpFilters) == 0 {
				t.Error("no http filter conf founded")
			}

			httpFilter := httpFilters[0]

			if httpFilter.FilterMixerConfig == nil {
				t.Error("no mixer config founded")
			}

			mixerConfig := httpFilter.FilterMixerConfig

			if mixerConfig.MixerAttributes == nil {
				t.Error("no mixer attributes founded")
			}

			if mixerConfig.ForwardAttributes == nil {
				t.Error("no mixer forward attributes founded")
			}

			mixerAttributres := mixerConfig.MixerAttributes
			if mixerAttributres.Attributes == nil {
				t.Error("no attributes for mixer attributes founded")
			}

			mixerAttributeDetail := mixerAttributres.Attributes

			if mixerAttributeDetail.DestinationIp == nil {
				t.Error("no destination ip founded")
			}

			if mixerAttributeDetail.DestinationIp.BytesValue != "AAAAAAAAAAAAAP//CigBEw==" {
				t.Error("destination ip not same")
			}

			if mixerAttributeDetail.DestinationUid == nil {
				t.Error("no destination Uid founded")
			}

			if mixerAttributeDetail.DestinationUid.StringValue != "kubernetes://productpage-v1-5fb67b856-6r5f2.default" {
				t.Error("destination uid not same")
			}
		
			return
		}


	}

	t.Error("no listener 10.40.1.19:9080 founded")
	
	/*
	if listeners[1].Filters[0].TCPProxyFilterConfig == nil {
		t.Error("TCPProxyFilterConfig is nil for the tcp filter")
	}
	if listeners[2].Filters[0].HTTPFilterConfig == nil {
		t.Error("HTTPFilterConfig is nil for the http filter")
	}
	if filter := listeners[2].Filters[0].HTTPFilterConfig.Filters[0]; filter.FilterMixerConfig == nil {
		t.Errorf("FilterMixerConfig is nil in %v", filter)
	}
	if filter := listeners[2].Filters[0].HTTPFilterConfig.Filters[1]; filter.FilterFaultConfig == nil {
		t.Errorf("FilterMixerConfig is nil in %v", filter)
	}
	*/
}
