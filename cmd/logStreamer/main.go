package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/RomanosTrechlis/logStreamer/streamer"
)

var (
	port, pport        int
	pprofInfo, console bool
	rootPath           string
	maxSize            int64
)

func init() {
	fmt.Printf("%s [INFO] Log streamer is starting...\n", streamer.PrintTime())
	// rpc server listening port
	flag.IntVar(&port, "port", 8080, "port for server to listen to requests")
	// enable/disable pprof functionality
	flag.BoolVar(&pprofInfo, "pprof", false,
		"additional server for pprof functionality")
	// enable/disable console dumps
	flag.BoolVar(&console, "console", false, "dumps log lines to console")
	// pprof port for http server
	flag.IntVar(&pport, "pport", 1111, "port for pprof server")
	// path must already exist
	flag.StringVar(&rootPath, "path", "../logs", "path for logs to be persisted")
	// the size of log files before they get renamed for storing purposes.
	size := flag.String("size", "1MB",
		"max size for individual files, -1B for infinite size")
	flag.Parse()

	i, err := streamer.LexicalToNumber(*size)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't parse size input to bytes: %v", err)
		os.Exit(2)
	}
	maxSize = i

	// prints some logo and info
	printLogo()
	infoBlock(port, pport, maxSize, rootPath, pprofInfo)
}

func main() {
	// validate path passed
	if err := streamer.CheckPath(rootPath); err != nil {
		fmt.Printf("path passed is not valid: %v\n", err)
		return
	}

	// stopAll channel listens to termination and interupt signals.
	stopAll := make(chan os.Signal, 1)
	signal.Notify(stopAll, syscall.SIGTERM, syscall.SIGINT)

	streamer := streamer.New(rootPath, port, maxSize)
	if pprofInfo {
		streamer.WithPProf(pport)
	}
	go streamer.Serve()
	<-stopAll
	streamer.Shutdown()

}
