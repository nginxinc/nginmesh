package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang/glog"

	"github.com/nginmesh/nginmesh/istio/agent/nginx"
	"github.com/nginmesh/nginmesh/istio/agent/pilot"
)

func main() {
	proxySidecarCmd := flag.NewFlagSet("proxy sidecar", flag.ContinueOnError)
	discoveryAddress := proxySidecarCmd.String("discoveryAddress", "", "Discovery address")
	serviceCluster := proxySidecarCmd.String("serviceCluster", "", "Service cluster")
	proxySidecarCmd.String("configPath", "", "Config path")
	proxySidecarCmd.String("binaryPath", "", "Binary path")
	proxySidecarCmd.String("drainDuration", "", "Binary path")
	proxySidecarCmd.String("parentShutdownDuration", "", "Binary path")
	proxySidecarCmd.String("discoveryRefreshDelay", "", "Binary path")
	proxySidecarCmd.String("zipkinAddress", "", "Binary path")
	proxySidecarCmd.String("connectTimeout", "", "Binary path")
	proxySidecarCmd.String("statsdUdpAddress", "", "Binary path")
	proxySidecarCmd.String("proxyAdminPort", "", "Binary path")
	collectorAddress := proxySidecarCmd.String("collectorAddress","","Collector address")
	collectorTopic := proxySidecarCmd.String("collectorTopic","","Collector topic")
	verbosity := proxySidecarCmd.String("v", "", "Verbosity level")

	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Invalid number of arguments")
		fmt.Fprintf(os.Stderr, "Usage: %v proxy sidecar arguments...\n", os.Args[0])
		proxySidecarCmd.PrintDefaults()
		os.Exit(1)
	}

	if os.Args[1] != "proxy" && os.Args[2] != "sidecar" {
		fmt.Fprintf(os.Stderr, "Usage: %v proxy sidecar arguments...\n", os.Args[0])
		proxySidecarCmd.PrintDefaults()
		os.Exit(1)
	}

	proxySidecarCmd.Parse(os.Args[3:])

	// configure glog
	flag.Parse()
	flag.Lookup("logtostderr").Value.Set("true")
	flag.Lookup("v").Value.Set(*verbosity)

	if *discoveryAddress == "" {
		fmt.Fprintf(os.Stderr, "No discoveryAddress was provided")
		fmt.Fprintf(os.Stderr, "Usage: %v proxy sidecar arguments...\n", os.Args[0])
		proxySidecarCmd.PrintDefaults()
		os.Exit(1)
	}

	if *serviceCluster == "" {
		fmt.Fprintf(os.Stderr, "No serviceCluster was provided")
		fmt.Fprintf(os.Stderr, "Usage: %v proxy sidecar arguments...\n", os.Args[0])
		proxySidecarCmd.PrintDefaults()
		os.Exit(1)
	}

	podIP := os.Getenv("INSTANCE_IP")
	if podIP == "" {
		podIP = os.Getenv("POD_IP")
	}

	podName := os.Getenv("POD_NAME")
	podNamespace := os.Getenv("POD_NAMESPACE")
	podUID := fmt.Sprintf("kubernetes://%v.%v", podName, podNamespace)
	serviceNode := fmt.Sprintf("sidecar~%v~%v.%v~%v.svc.cluster.local", podIP, podName, podNamespace, podNamespace)

	configVars := nginx.ConfigVariables{
		BindAddress:    podIP,
		ServiceCluster: *serviceCluster,
		ServiceNode:    serviceNode,
	}

	if os.Getenv("DISABLE_MIXER_REPORT") == "1" {
		glog.V(2).Info("Mixer report is disabled")
		configVars.DisableMixerReport = true
	}

	if os.Getenv("DISABLE_MIXER_CHECK") == "1" {
		glog.V(2).Info("Mixer check is disabled")
		configVars.DisableMixerCheck = true
	}

	if os.Getenv("DISABLE_TRACING") == "1" {
		glog.V(2).Info("Tracing is disabled")
		configVars.DisableTracing = true
	}

	converter := nginx.NewConverter(&configVars)

	glog.Infof("Starting the agent on %v at %v", podUID, podIP)

	if err := os.Mkdir("/etc/istio/proxy/conf.d", 0755); err != nil && !os.IsExist(err) {
		glog.Fatalf("Couldn't create the conf.d folder: %v", err)
	}

	if err := os.Mkdir("/etc/istio/proxy/cache", 0755); err != nil && !os.IsExist(err) {
		glog.Fatalf("Couldn't create the cache folder: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	client := pilot.NewClient(*discoveryAddress, &http.Client{}, *serviceCluster, serviceNode, podIP,*collectorAddress,*collectorTopic)
	glog.Info("collector address: %v, topic: %v",*collectorAddress,*collectorTopic)
	pilotWatcher := pilot.NewWatcher(client, 5*time.Second)
	go pilotWatcher.Run(ctx)

	fsWatcher, err := NewFSWatcher([]string{"/etc/certs/"})
	if err != nil {
		glog.Fatalf("Failed to create FSWatcher: %v", err)
	}
	go fsWatcher.Run(ctx)

	nginxCtrl := nginx.NewController()

	stop := make(chan struct{})
	go handleSignals(stop)

mainLoop:
	for {
		select {
		case <-stop:
			glog.V(2).Info("Terminating the agent")
			err := nginxCtrl.Quit()
			if err != nil {
				glog.Errorf("Couldn't gracefully shutdown NGINX: %v", err)
			}
			break mainLoop
		case err := <-nginxCtrl.ExitStatus():
			if err != nil {
				glog.Errorf("NGINX unexpectedly exited with an error: %v", err)
			} else {
				glog.Errorf("NGINX unexpectedly exited successfully")
			}
			break mainLoop
		case proxyCfg := <-pilotWatcher.GetConfigUpdates():
			glog.V(2).Info("Configuration in Pilot has been changed. Generating new NGINX configuration")
			config := converter.Convert(proxyCfg)
			err := nginxCtrl.ApplyConfig(config)
			if err != nil {
				glog.Fatalf("Couldn't apply new configuration: %v", err)
			}
		case change := <-fsWatcher.Changes():
			glog.V(2).Infof("Reloading NGINX due to changes on the filesystem: %v", change)
			err := nginxCtrl.Reload()
			if err != nil {
				glog.Fatalf("Couldn't reload NGINX: %v", err)
			}
		}
	}

	glog.V(2).Info("Canceling the main context")
	cancel()

	// wait until NGINX gracefully quits if it is running
	err, ok := <-nginxCtrl.ExitStatus()
	if ok {
		if err != nil {
			glog.Errorf("NGINX exited with an error: %v", err)
		} else {
			glog.Info("NGINX exited successfully")
		}
	}
}

func handleSignals(stop chan struct{}) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	<-sigs
	close(stop)
}
