package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
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
	"github.com/RomanosTrechlis/go-scribe/types"
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

func fillAgentConfig(flags map[string]string) (*types.AgentConfig, error) {
	file := c.StringValue("file", "agent", flags)
	if file != "" {
		return fillAgentConfigFromFile(file)
	}
	return fillAgentConfigFromFlags(flags)
}

func fillAgentConfigFromFlags(flags map[string]string) (*types.AgentConfig, error) {
	port, err := c.IntValue("port", "agent", flags)
	if err != nil {
		return nil, fmt.Errorf("failed to get the value of 'port' flag: %v", err)
	}
	pport, err := c.IntValue("pport", "agent", flags)
	if err != nil {
		return nil, fmt.Errorf("failed to get the value of 'pport' flag: %v", err)
	}
	pprofInfo, err := c.BoolValue("pprof", "agent", flags)
	if err != nil {
		return nil, fmt.Errorf("failed to get the value of 'pprof' flag: %v", err)
	}
	mediator := c.StringValue("mediator", "agent", flags)
	maxSize, err := scribe.LexicalToNumber(c.StringValue("size", "agent", flags))
	if err != nil {
		return nil, fmt.Errorf("failed to get the value of 'size' flag: %v", err)
	}
	path := c.StringValue("path", "agent", flags)
	verbose, err := c.BoolValue("verbose", "agent", flags)
	if err != nil {
		return nil, fmt.Errorf("failed to get the value of 'verbose' flag: %v", err)
	}
	console, err := c.BoolValue("console", "agent", flags)
	if err != nil {
		return nil, fmt.Errorf("failed to get the value of 'console' flag: %v", err)
	}

	crt := c.StringValue("crt", "agent", flags)
	pk := c.StringValue("pk", "agent", flags)
	ca := c.StringValue("ca", "agent", flags)
	a := &types.AgentConfig{port, pprofInfo, console, verbose,
	mediator, pport, path, maxSize,
	types.CertificateConfig{crt, pk, ca}}
	return a, nil
}

func fillAgentConfigFromFile(file string) (*types.AgentConfig, error) {
	b, err := readConfigurationFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %v", err)
	}

	m := new(types.AgentConfig)
	yaml.Unmarshal(b, m)
	return m, nil
}

func agentHandler(flags map[string]string) error {
	conf, err := fillAgentConfig(flags)
	if err != nil {
		return err
	}

	printLogoAgent()

	id := xid.New().String()
	p.Print(fmt.Sprintf("Scribe's id: %s", id))

	// stopAll channel listens to termination and interrupt signals.
	stopAll := make(chan os.Signal, 1)
	signal.Notify(stopAll, syscall.SIGTERM, syscall.SIGINT)

	// register to mediator
	if conf.Mediator != "" {
		err := addMediator(id, conf)
		if err != nil {
			return fmt.Errorf("failed to connect to mediator %s: %v", conf.Mediator, err)
		}
	}

	// validate path passed
	if err := scribe.CheckPath(conf.LogPath); err != nil {
		fmt.Fprintf(os.Stderr, "path passed is not valid: %v\n", err)
		os.Exit(2)
	}

	s, err := scribe.New(id, conf.LogPath, conf.Port, conf.LogFileSize, conf.Mediator,
		conf.Certificate, conf.PrivateKey, conf.CertificateAuthority)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to create scribe: %v", err)
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}

	infoBlock(conf)

	defer s.Shutdown()
	go s.Serve()

	var srv *http.Server
	if conf.Profile {
		srv = profiling.Serve(conf.ProfilePort)
		defer srv.Shutdown(nil)
	}

	if conf.Verbose {
		go s.Tick(20)
	}

	cliServer, _ := gserver.New("", "", "")
	c := cliScribe{false, nil, s}
	go gserver.Serve(registerCLIScribeFunc(cliServer, c), fmt.Sprint(":4242"), cliServer)

	<-stopAll
	return nil
}

func addMediator(id string, conf *types.AgentConfig) error {
	conn, err := grpc.Dial(conf.Mediator,
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
	req := &pb.RegisterRequest{
		Id:   id,
		Addr: fmt.Sprintf("%s:%d", host, conf.Port),
	}
	var retries = 3
	var success = false
	for retries > 0 {
		r, err := cl.Register(context.Background(), req)
		if err != nil {
			retries--
			p.Print(fmt.Sprintf("Failed to register to mediator '%s. "+
				"Remaining tries: %d", conf.Mediator, retries))
			time.Sleep(1 * time.Second)
			continue
		}
		if r.GetRes() != "Success" {
			retries--
			p.Print(fmt.Sprintf("Failed to register to mediator '%s'. "+
				"Remaining tries: %d", conf.Mediator, retries))
			time.Sleep(2 * time.Second)
			continue
		}
		success = true
		break
	}
	if !success {
		return fmt.Errorf("failed to register scribe to mediator '%s'\n", conf.Mediator)
	}

	p.Print("Successfully registered to mediator")
	return nil
}

func infoBlock(conf *types.AgentConfig) {
	fmt.Println("##########################################################")
	fmt.Println("\t==>\tPort number:\t", conf.Port)
	fmt.Println("\t==>\tLog path:\t", conf.LogPath)
	fmt.Println("\t==>\tLog size:\t", conf.LogFileSize)
	fmt.Println("\t==>\tPprof server:\t", conf.Profile)
	fmt.Println("\t==>\tPprof port:\t", conf.ProfilePort)
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
