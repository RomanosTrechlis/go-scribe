// +build go1.8

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"net/http"

	pb "github.com/RomanosTrechlis/logScribe/api"
	"github.com/RomanosTrechlis/logScribe/profiling"
	"github.com/RomanosTrechlis/logScribe/scribe"
	"github.com/rs/xid"
	"google.golang.org/grpc"

	p "github.com/RomanosTrechlis/logScribe/util/format/print"
)

var (
	port, pport        int
	pprofInfo, console bool
	rootPath           string
	maxSize            int64
	mediator           string
	verbose 		   bool
)

var (
	cert = "certs/server.crt"
	key  = "certs/server.key"
	ca   = "certs/CertAuth.crt"
)

func init() {
	flag.IntVar(&port, "port", 8080, "port for server to listen to requests")
	flag.BoolVar(&pprofInfo, "pprof", false, "additional server for pprof functionality")
	flag.BoolVar(&console, "console", false, "dumps log lines to console")
	flag.BoolVar(&verbose, "verbose", false, "prints regular handled request count")
	flag.StringVar(&mediator, "mediator", "","mediators address if exists, i.e 127.0.0.1:8080")
	flag.IntVar(&pport, "pport", 1111, "port for pprof server")
	flag.StringVar(&rootPath, "path", "../logs", "path for logs to be persisted")
	size := flag.String("size", "1MB", "max size for individual files, -1B for infinite size")

	// certificate files
	flag.StringVar(&cert, "crt", "", "host's certificate for secured connections")
	flag.StringVar(&key, "pk", "", "host's private key")
	flag.StringVar(&ca, "ca", "", "certificate authority's certificate")
	flag.Parse()

	i, err := scribe.LexicalToNumber(*size)
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't parse size input to bytes: %v", err)
		os.Exit(2)
	}
	maxSize = i

	// prints some logo and info
	printLogo()
	infoBlock(port, pport, maxSize, rootPath, pprofInfo)
}

func addMediator() {
	conn, err := grpc.Dial(mediator,
		grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewRegisterClient(conn)
	host, err := os.Hostname()
	if err != nil {
		return
	}
	req := &pb.RegisterRequest{
		Id:   id,
		Addr: fmt.Sprintf("%s:%d", host, port),
	}
	var retries = 3
	var success = false
	for retries > 0 {
		r, err := c.Register(context.Background(), req)
		if err != nil {
			retries--
			p.Print(fmt.Sprintf("Failed to register to mediator '%s. "+
				"Remaining tries: %d", mediator, retries))
			time.Sleep(1 * time.Second)
			continue
		}
		if r.GetRes() != "Success" {
			retries--
			p.Print(fmt.Sprintf("Failed to register to mediator '%s'. "+
				"Remaining tries: %d", mediator, retries))
			time.Sleep(2 * time.Second)
			continue
		}
		success = true
		break
	}
	if !success {
		fmt.Fprintf(os.Stderr, "failed to register scribe to mediator '%s'\n",
			mediator)
		os.Exit(2)
	}

	p.Print("Successfully registered to mediator")
}

var id string

func main() {
	// validate path passed
	if err := scribe.CheckPath(rootPath); err != nil {
		fmt.Printf("path passed is not valid: %v\n", err)
		return
	}

	id = xid.New().String()
	p.Print(fmt.Sprintf("Scribe's id: %s", id))

	// stopAll channel listens to termination and interupt signals.
	stopAll := make(chan os.Signal, 1)
	signal.Notify(stopAll, syscall.SIGTERM, syscall.SIGINT)

	// register to mediator
	if mediator != "" {
		addMediator()
	}

	s, err := scribe.New(rootPath, port, maxSize, mediator, cert, key, ca)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	defer s.Shutdown()
	go s.Serve()

	var srv *http.Server
	if pprofInfo {
		srv = profiling.Serve(pport)
		defer srv.Shutdown(nil)
	}

	if verbose {
		go s.Tick(20)
	}

	<-stopAll
}
