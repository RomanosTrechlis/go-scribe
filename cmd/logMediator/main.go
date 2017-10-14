// +build go1.8

package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/RomanosTrechlis/logScribe/mediator"
	"github.com/RomanosTrechlis/logScribe/profiling"
)

var (
	port, pport int
	pprofInfo   bool
)

var (
	cert = "certs/server.crt"
	key  = "certs/server.key"
	ca   = "certs/CertAuth.crt"
)

func init() {
	// rpc server listening port
	flag.IntVar(&port, "port", 8000, "port for mediator server to listen to requests")
	// enable/disable pprof functionality
	flag.BoolVar(&pprofInfo, "pprof", false,
		"additional server for pprof functionality")
	// pprof port for http server
	flag.IntVar(&pport, "pport", 1111, "port for pprof server")
	flag.StringVar(&cert, "crt", "", "host's certificate for secured connections")
	flag.StringVar(&key, "pk", "", "host's private key")
	flag.StringVar(&ca, "ca", "", "certificate authority's certificate")
	flag.Parse()

	printLogo()
}

func main() {
	// stopAll channel listens to termination and interupt signals.
	stopAll := make(chan os.Signal, 1)
	signal.Notify(stopAll, syscall.SIGTERM, syscall.SIGINT)

	m, err := mediator.New(port, cert, key, ca)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to start a new mediator: %v", err)
		os.Exit(2)
	}
	defer m.Shutdown()
	go m.Serve()

	var srv *http.Server
	if pprofInfo {
		srv = profiling.Serve(pport)
		defer srv.Shutdown(nil)
	}
	<-stopAll
}
