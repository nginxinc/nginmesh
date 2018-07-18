package pilot

import (
	"testing"
	"io/ioutil"
	"fmt"
	"os"
)


func TestUnmarshal(t *testing.T) {

	pwd, _ := os.Getwd()
	body, ferr := ioutil.ReadFile(pwd+"/test/0.7.1/listener.json")
	if ferr != nil {
		fmt.Printf("could not load file")
		t.Error(ferr)
		return
	}


	listeners, err := unMarshalListeners(body)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(len(listeners))
	if len(listeners) != 23 {
		t.Error("there should been 23 listener founded")
	}

	for _, listener := range listeners {
		if listener.Address == "tcp://0.0.0.0:15007" {

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
			// fmt.Printf("Filters: %+v\n", httpFilter)
			fmt.Printf("config: %+v\n", string(httpFilter.Config))
			if httpFilter.FilterMixerConfig == nil {
				t.Error("no mixer config founded")
			}

			mixerConfig := httpFilter.FilterMixerConfig
			fmt.Printf("mixerConfig: %+v\n", mixerConfig)
			fmt.Printf("mixerConfig: %+v\n", mixerConfig.DestinationService)
			fmt.Printf("mixer attributes: %+v\n", mixerConfig.ForwardAttributes.Attributes.SourceLabels.StringMapValue.Entries)
			fmt.Printf("service configs: %+v\n", mixerConfig.ServiceConfig[mixerConfig.DestinationService])

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

			ForwardAttributes := mixerConfig.ForwardAttributes
			if ForwardAttributes.Attributes == nil {
				t.Error("no attributes for forward attributes founded")
			}

			fmt.Printf("Source IP: %+v\n", ForwardAttributes.Attributes.SourceIp.BytesValue)

			return
		}


	}

	// t.Error("no listener 10.40.1.19:9080 founded")

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