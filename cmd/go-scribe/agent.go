package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/RomanosTrechlis/go-scribe/api"
	"github.com/RomanosTrechlis/go-scribe/profiling"
	"github.com/RomanosTrechlis/go-scribe/scribe"
	p "github.com/RomanosTrechlis/go-scribe/util/format/print"
	"github.com/rs/xid"
	"google.golang.org/grpc"
)

const (
	agentShortHelp = `starts a scribe agent`
	agentLongHelp  = `go-scribe agent starts a scribe agent.

Before starting the agent the logging path must exists.
Then it sets up a gRPC server to receives logging
requests.

Agent supports a 2-way-SSL authentication by passing the
certificate, the private key, and the certificate authority
file names.

There is also support for profiling the server it runs by
passing the pprof flag and the pport to access it.
	`
)

type agent struct {
	port      int
	pprofInfo bool
	console   bool
	verbose   bool
	mediator  string
	pport     int
	rootPath  string
	size      string

	cert string
	key  string
	ca   string
}

func (cmd *agent) Name() string      { return "agent" }
func (cmd *agent) Args() string      { return "" }
func (cmd *agent) ShortHelp() string { return agentShortHelp }
func (cmd *agent) LongHelp() string  { return agentLongHelp }
func (cmd *agent) Hidden() bool      { return false }

func (cmd *agent) Register(fs *flag.FlagSet) {
	fs.IntVar(&cmd.port, "port", 8080, "port for server to listen to requests")
	fs.BoolVar(&cmd.pprofInfo, "pprof", false, "additional server for pprof functionality")
	fs.BoolVar(&cmd.console, "console", false, "dumps log lines to console")
	fs.BoolVar(&cmd.verbose, "verbose", false, "prints regular handled request count")
	fs.StringVar(&cmd.mediator, "mediator", "", "mediators address if exists, i.e 127.0.0.1:8080")
	fs.IntVar(&cmd.pport, "pport", 1111, "port for pprof server")
	fs.StringVar(&cmd.rootPath, "path", "../logs", "path for logs to be persisted")
	fs.StringVar(&cmd.size, "size", "1MB", "max size for individual files, -1B for infinite size")
	// certificate files
	fs.StringVar(&cmd.cert, "crt", "", "host's certificate for secured connections")
	fs.StringVar(&cmd.key, "pk", "", "host's private key")
	fs.StringVar(&cmd.ca, "ca", "", "certificate authority's certificate")
}

func (cmd *agent) Run(ctx *ctx, args []string) error {
	// validate path passed
	if err := scribe.CheckPath(cmd.rootPath); err != nil {
		return fmt.Errorf("path passed is not valid: %v\n", err)
	}

	maxSize, err := scribe.LexicalToNumber(cmd.size)
	if err != nil {
		return fmt.Errorf("couldn't parse size input to bytes: %v", err)
	}

	printLogoAgent()
	infoBlock(cmd.port, cmd.pport, maxSize, cmd.rootPath, cmd.pprofInfo)

	id := xid.New().String()
	p.Print(fmt.Sprintf("Scribe's id: %s", id))

	// stopAll channel listens to termination and interupt signals.
	stopAll := make(chan os.Signal, 1)
	signal.Notify(stopAll, syscall.SIGTERM, syscall.SIGINT)

	// register to mediator
	if cmd.mediator != "" {
		err := cmd.addMediator(id)
		if err != nil {
			return fmt.Errorf("failed to connect to mediator %s: %v", cmd.mediator, err)
		}
	}

	s, err := scribe.New(cmd.rootPath, cmd.port, maxSize, cmd.mediator, cmd.cert, cmd.key, cmd.ca)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(2)
	}
	defer s.Shutdown()
	go s.Serve()

	var srv *http.Server
	if cmd.pprofInfo {
		srv = profiling.Serve(cmd.pport)
		defer srv.Shutdown(nil)
	}

	if cmd.verbose {
		go s.Tick(20)
	}

	<-stopAll
	return nil
}

func (cmd *agent) addMediator(id string) error {
	conn, err := grpc.Dial(cmd.mediator,
		grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	c := pb.NewRegisterClient(conn)
	host, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname from system: %v", err)
	}
	req := &pb.RegisterRequest{
		Id:   id,
		Addr: fmt.Sprintf("%s:%d", host, cmd.port),
	}
	var retries = 3
	var success = false
	for retries > 0 {
		r, err := c.Register(context.Background(), req)
		if err != nil {
			retries--
			p.Print(fmt.Sprintf("Failed to register to mediator '%s. "+
				"Remaining tries: %d", cmd.mediator, retries))
			time.Sleep(1 * time.Second)
			continue
		}
		if r.GetRes() != "Success" {
			retries--
			p.Print(fmt.Sprintf("Failed to register to mediator '%s'. "+
				"Remaining tries: %d", cmd.mediator, retries))
			time.Sleep(2 * time.Second)
			continue
		}
		success = true
		break
	}
	if !success {
		return fmt.Errorf("failed to register scribe to mediator '%s'\n", cmd.mediator)
	}

	p.Print("Successfully registered to mediator")
	return nil
}

func printLogoAgent() {
	fmt.Println()
	fmt.Println("██╗      ██████╗  ██████╗     ███████╗ ██████╗██████╗ ██╗██████╗ ███████╗")
	fmt.Println("██║     ██╔═══██╗██╔════╝     ██╔════╝██╔════╝██╔══██╗██║██╔══██╗██╔════╝")
	fmt.Println("██║     ██║   ██║██║  ███╗    ███████╗██║     ██████╔╝██║██████╔╝█████╗  ")
	fmt.Println("██║     ██║   ██║██║   ██║    ╚════██║██║     ██╔══██╗██║██╔══██╗██╔══╝  ")
	fmt.Println("███████╗╚██████╔╝╚██████╔╝    ███████║╚██████╗██║  ██║██║██████╔╝███████╗")
	fmt.Println("╚══════╝ ╚═════╝  ╚═════╝     ╚══════╝ ╚═════╝╚═╝  ╚═╝╚═╝╚═════╝ ╚══════╝")
}

func infoBlock(port, pport int, maxSize int64, path string, pprofInfo bool) {
	// info block
	fmt.Println("##########################################################")
	fmt.Println("\t==>\tPort number:\t", port)
	fmt.Println("\t==>\tLog path:\t", path)
	fmt.Println("\t==>\tLog size:\t", maxSize)
	fmt.Println("\t==>\tPprof server:\t", pprofInfo)
	fmt.Println("\t==>\tPprof port:\t", pport)
	fmt.Println("##########################################################")
}
