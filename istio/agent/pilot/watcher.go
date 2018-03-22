// Package pilot provides a client for getting Envoy load balancing configuration from Pilot.
// Additionally, it provides a watcher that periodically pulls configuration from Pilot
// and notifies in case the configuration has been updated.
package pilot

import (
	"context"
	"reflect"
	"time"

	"github.com/golang/glog"
)

// Watcher watches load balancing configuration in Pilot and sends the new configuration to
// a special channel when the configuration was updated.
type Watcher struct {
	client          *Client
	refreshInterval time.Duration
	configCh        chan ProxyConfig
}

// NewWatcher creates a NewWatcher.
func NewWatcher(client *Client, refreshInterval time.Duration) *Watcher {
	return &Watcher{
		client:          client,
		refreshInterval: refreshInterval,
		configCh:        make(chan ProxyConfig),
	}
}

// Run starts the watcher.
func (w *Watcher) Run(ctx context.Context) {
	var oldConfig ProxyConfig
	for {
		select {
		case <-ctx.Done():
			glog.V(2).Infof("Terminating Pilot Watcher")
			return
		case <-time.After(w.refreshInterval):
			cfg := w.client.GetConfig()
			glog.V(2).Info("%+v\n", cfg)
			if !reflect.DeepEqual(cfg, oldConfig) {
				glog.V(2).Info("Configuration in Pilot has been changed")
				w.configCh <- cfg
				oldConfig = cfg
			} else {
				glog.V(2).Info("Configuration in Pilot has not been changed")
			}
		}
	}
}

// GetConfigUpdates returns a chanell where the updated configuration is sent to.
func (w *Watcher) GetConfigUpdates() <-chan ProxyConfig {
	return w.configCh
}
