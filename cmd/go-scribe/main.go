package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/RomanosTrechlis/go-icls/cli"
	pb "github.com/RomanosTrechlis/go-scribe/api"
	med "github.com/RomanosTrechlis/go-scribe/mediator"
	"github.com/RomanosTrechlis/go-scribe/profiling"
	"github.com/RomanosTrechlis/go-scribe/scribe"
	p "github.com/RomanosTrechlis/go-scribe/internal/util/format/print"
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

	mediatorShortHelp = `starts a scribe mediator service`
	mediatorLongHelp  = `go-scribe mediator starts a scribe mediator service.

Mediator is a functionality of go-scribe that instead of
writing logs, it delegates this responsibility to a registered
scribe agent.
The algorithm deciding the responsible agent is very simple.

It, also, supports 2-way-SSL authentication by passing from the
flags the certificate, the private key, and the certificate
authority filename.

There is also support for profiling the server it runs by
passing the pprof flag and the pport to access it.
	`
)

var c *cli.CLI

func init() {
	c = cli.New()
	agent := c.New("agent", agentShortHelp, agentLongHelp, agentHandler)
	agent.IntFlag("port", "", 8080, "port for server to listen to requests", false)
	agent.BoolFlag("pprof", "", false, "additional server for pprof functionality", false)
	agent.BoolFlag("console", "", false, "dumps log lines to console", false)
	agent.BoolFlag("verbose", "", false, "prints regular handled request count", false)
	agent.StringFlag("mediator", "", "", "mediators address if exists, i.e 127.0.0.1:8080", false)
	agent.IntFlag("pport", "", 1111, "port for pprof server", false)
	agent.StringFlag("path", "", "../../logs", "path for logs to be persisted", false)
	agent.StringFlag("size", "", "1MB", "max size for individual files, -1B for infinite size", false)
	agent.StringFlag("crt", "", "", "host's certificate for secured connections", false)
	agent.StringFlag("pk", "", "", "host's private key", false)
	agent.StringFlag("ca", "", "", "certificate authority's certificate", false)

	med := c.New("mediator", mediatorShortHelp, mediatorLongHelp, mediatorHandler)
	med.IntFlag("port", "", 8000, "port for mediator server to listen to requests", false)
	med.BoolFlag("pprof", "", false, "additional server for pprof functionality", false)
	med.IntFlag("pport", "", 2222, "port for pprof server", false)
	med.StringFlag("crt", "", "", "host's certificate for secured connections", false)
	med.StringFlag("pk", "", "", "host's private key", false)
	med.StringFlag("ca", "", "", "certificate authority's certificate", false)
}

func agentHandler(flags map[string]string) error {
	printLogoAgent()

	id := xid.New().String()
	p.Print(fmt.Sprintf("Scribe's id: %s", id))

	// stopAll channel listens to termination and interupt signals.
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
		fmt.Fprintf(os.Stderr,"path passed is not valid: %v\n", err)
		os.Exit(2)
	}
	maxSize, err := scribe.LexicalToNumber(c.StringValue("size", "agent", flags))
	if err != nil {
		fmt.Fprintf(os.Stderr,"couldn't parse size input to bytes: %v", err)
		os.Exit(2)
	}
	s, err := scribe.New(c.StringValue("path", "agent", flags), port, maxSize, c.StringValue("mediator", "agent", flags),
		c.StringValue("crt", "agent", flags), c.StringValue("pk", "agent", flags), c.StringValue("ca", "agent", flags))
	if err != nil {
		fmt.Fprintf(os.Stderr,"failed to create scribe: %v", err)
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

	<-stopAll
	return nil
}

func mediatorHandler(flags map[string]string) error {
	printLogo()

	// stopAll channel listens to termination and interupt signals.
	stopAll := make(chan os.Signal, 1)
	signal.Notify(stopAll, syscall.SIGTERM, syscall.SIGINT)

	port, _ := c.IntValue("port", "mediator", flags)
	pport, _ := c.IntValue("pport", "mediator", flags)
	m, err := med.New(port, c.StringValue("crt", "mediator", flags), c.StringValue("pk", "mediator", flags), c.StringValue("ca", "mediator", flags))
	if err != nil {
		return fmt.Errorf("failed to start a new mediator: %v", err)
	}
	defer m.Shutdown()
	go m.Serve()

	var srv *http.Server
	pprofInfo, err := c.BoolValue("pprof", "mediator", flags)
	if pprofInfo {
		srv = profiling.Serve(pport)
		defer srv.Shutdown(nil)
	}
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
	host, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("failed to get hostname from system: %v", err)
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

func printLogo() {
	fmt.Println(" ___      _______  _______    __   __  _______  ______   ___   _______  _______  _______  ______   ")
	fmt.Println("|   |    |       ||       |  |  |_|  ||       ||      | |   | |   _   ||       ||       ||    _ |  ")
	fmt.Println("|   |    |   _   ||    ___|  |       ||    ___||  _    ||   | |  |_|  ||_     _||   _   ||   | ||  ")
	fmt.Println("|   |    |  | |  ||   | __   |       ||   |___ | | |   ||   | |       |  |   |  |  | |  ||   |_||_ ")
	fmt.Println("|   |___ |  |_|  ||   ||  |  |       ||    ___|| |_|   ||   | |       |  |   |  |  |_|  ||    __  |")
	fmt.Println("|       ||       ||   |_| |  | ||_|| ||   |___ |       ||   | |   _   |  |   |  |       ||   |  | |")
	fmt.Println("|_______||_______||_______|  |_|   |_||_______||______| |___| |__| |__|  |___|  |_______||___|  |_|")
}

func main() {
	args := os.Args
	if len(args) == 1 {
		return
	}
	line := ""
	for _, a :=  range args[1:] {
		line += a + " "
	}
	c.Execute(line)
}
