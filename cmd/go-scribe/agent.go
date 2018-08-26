package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	pb "github.com/RomanosTrechlis/go-scribe/api"
	p "github.com/RomanosTrechlis/go-scribe/internal/util/format/print"
	"github.com/RomanosTrechlis/go-scribe/internal/util/gserver"
	"github.com/RomanosTrechlis/go-scribe/internal/util/net"
	"github.com/RomanosTrechlis/go-scribe/profiling"
	"github.com/RomanosTrechlis/go-scribe/scribe"
	"github.com/rs/xid"
	"golang.org/x/net/context"
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

func agentHandler(flags map[string]string) error {
	printLogoAgent()

	id := xid.New().String()
	p.Print(fmt.Sprintf("Scribe's id: %s", id))

	// stopAll channel listens to termination and interrupt signals.
	stopAll := make(chan os.Signal, 1)
	signal.Notify(stopAll, syscall.SIGTERM, syscall.SIGINT)

	// register to mediator
	mediator := c.StringValue("mediator", "agent", flags)
	if mediator != "" {
		err := addMediator(id, "agent", flags)
		if err != nil {
			return fmt.Errorf("failed to connect to mediator %s: %v", c.StringValue("mediator", "agent", flags), err)
		}
	}

	port, _ := c.IntValue("port", "agent", flags)

	// validate path passed
	if err := scribe.CheckPath(c.StringValue("path", "agent", flags)); err != nil {
		fmt.Fprintf(os.Stderr, "path passed is not valid: %v\n", err)
		os.Exit(2)
	}
	maxSize, err := scribe.LexicalToNumber(c.StringValue("size", "agent", flags))
	if err != nil {
		fmt.Fprintf(os.Stderr, "couldn't parse size input to bytes: %v", err)
		os.Exit(2)
	}
	s, err := scribe.New(id, c.StringValue("path", "agent", flags), port, maxSize, c.StringValue("mediator", "agent", flags),
		c.StringValue("crt", "agent", flags), c.StringValue("pk", "agent", flags), c.StringValue("ca", "agent", flags))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create scribe: %v", err)
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}

	infoBlock("agent", flags)

	defer s.Shutdown()
	go s.Serve()

	var srv *http.Server
	pprofInfo, _ := c.BoolValue("pprof", "agent", flags)
	pport, _ := c.IntValue("pport", "agent", flags)
	if pprofInfo {
		srv = profiling.Serve(pport)
		defer srv.Shutdown(nil)
	}

	verbose, _ := c.BoolValue("verbose", "agent", flags)
	if verbose {
		go s.Tick(20)
	}

	cliServer, _ := gserver.New("", "", "")
	c := cliScribe{false, nil, s}
	go gserver.Serve(registerCLIScribeFunc(cliServer, c), fmt.Sprint(":4242"), cliServer)

	<-stopAll
	return nil
}

func addMediator(id, cmd string, flags map[string]string) error {
	conn, err := grpc.Dial(c.StringValue("mediator", cmd, flags),
		grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	cl := pb.NewRegisterClient(conn)
	host, err := net.GetIPAddress()
	if err != nil {
		return fmt.Errorf("failed to get ip/hostname from system: %v", err)
	}
	port, _ := c.IntValue("port", cmd, flags)
	req := &pb.RegisterRequest{
		Id:   id,
		Addr: fmt.Sprintf("%s:%d", host, port),
	}
	var retries = 3
	var success = false
	for retries > 0 {
		r, err := cl.Register(context.Background(), req)
		if err != nil {
			retries--
			p.Print(fmt.Sprintf("Failed to register to mediator '%s. "+
				"Remaining tries: %d", c.StringValue("mediator", cmd, flags), retries))
			time.Sleep(1 * time.Second)
			continue
		}
		if r.GetRes() != "Success" {
			retries--
			p.Print(fmt.Sprintf("Failed to register to mediator '%s'. "+
				"Remaining tries: %d", c.StringValue("mediator", cmd, flags), retries))
			time.Sleep(2 * time.Second)
			continue
		}
		success = true
		break
	}
	if !success {
		return fmt.Errorf("failed to register scribe to mediator '%s'\n", c.StringValue("mediator", cmd, flags))
	}

	p.Print("Successfully registered to mediator")
	return nil
}

func infoBlock(cmd string, flags map[string]string) {
	port, _ := c.IntValue("port", cmd, flags)
	pport, _ := c.IntValue("pport", cmd, flags)
	pprofInfo, _ := c.BoolValue("pprof", cmd, flags)
	rootPath := c.StringValue("path", cmd, flags)
	size := c.StringValue("size", cmd, flags)

	fmt.Println("##########################################################")
	fmt.Println("\t==>\tPort number:\t", port)
	fmt.Println("\t==>\tLog path:\t", rootPath)
	maxSize, _ := scribe.LexicalToNumber(size)
	fmt.Println("\t==>\tLog size:\t", maxSize)
	fmt.Println("\t==>\tPprof server:\t", pprofInfo)
	fmt.Println("\t==>\tPprof port:\t", pport)
	fmt.Println("##########################################################")
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
