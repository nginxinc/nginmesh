package main

import (
	"flag"
	"fmt"
	"os"
)

var filename = flag.String("f", "", "Input Kubernetes resource filename")
var initImage = flag.String("init-image", "nginmesh/istio-nginx-init:0.16-beta", "Image of the init container")
var proxyImage = flag.String("proxy-image", "nginmesh/istio-nginx-sidecar:0.16-beta", "Image of the proxy container")
var proxyUID = flag.Int64("proxy-uid", 104, "UID of the proxy process")
var proxyListenPort = flag.Int("proxy-port", 15001, "The port on which the proxy listens")

func main() {
	flag.Parse()

	if *filename == "" {
		fmt.Fprintln(os.Stderr, "No input file was provided")
		flag.Usage()
		os.Exit(1)
	}

	in, err := os.Open(*filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't open the input file %v:%v", *filename, err)
		os.Exit(1)
	}

	params := params{
		InitImage:       *initImage,
		ProxyImage:      *proxyImage,
		SidecarProxyUID: *proxyUID,
		ProxyListenPort: *proxyListenPort,
	}

	err = injectIntoResourceFile(&params, in, os.Stdout)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Couldn't inject containers into the file: %v", err)
		os.Exit(1)
	}
}
