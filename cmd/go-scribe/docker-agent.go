package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	pb "github.com/RomanosTrechlis/go-scribe/api"
	"github.com/RomanosTrechlis/go-scribe/profiling"
	"github.com/RomanosTrechlis/go-scribe/scribe"
	p "github.com/RomanosTrechlis/go-scribe/util/format/print"
	rpc "github.com/RomanosTrechlis/go-scribe/util/gserver"
	"github.com/rs/xid"
	"google.golang.org/grpc"
)

const (
	dockerAgentShortHelp = `starts a scribe agent in docker container`
	dockerAgentLongHelp  = `go-scribe docker-agent starts a scribe agent for docker container.

The configuration for docker-agent comes from enviroment variables.

It sets up a gRPC server to receives logging
requests.

Agent supports a 2-way-SSL authentication by passing the
certificate, the private key, and the certificate authority
file names.

There is also support for profiling the server it runs by
passing the pprof flag and the pport to access it.
	`
)

type dockerAgent struct {
	port      int
	pprofInfo bool
	console   bool
	verbose   bool
	mediator  string
	pport     int
	rootPath  string
	size      string
	dbServer  string
	dbName    string

	cert string
	key  string
	ca   string
}

func (cmd *dockerAgent) Name() string      { return "docker-agent" }
func (cmd *dockerAgent) Args() string      { return "" }
func (cmd *dockerAgent) ShortHelp() string { return dockerAgentShortHelp }
func (cmd *dockerAgent) LongHelp() string  { return dockerAgentLongHelp }
func (cmd *dockerAgent) Hidden() bool      { return false }

func (cmd *dockerAgent) Register(fs *flag.FlagSet) {}

func (cmd *dockerAgent) Run(ctx *ctx, args []string) error {
	cmd.port, _ = strconv.Atoi(os.Getenv("AGENT_PORT"))
	cmd.pprofInfo, _ = strconv.ParseBool(os.Getenv("AGENT_PPROF"))
	cmd.dbName = os.Getenv("AGENT_DB_NAME")
	cmd.dbServer = os.Getenv("AGENT_DB_SERVER")

	cmd.cert = os.Getenv("AGENT_CERT")
	cmd.key = os.Getenv("AGENT_KEY")
	cmd.ca = os.Getenv("AGENT_CA")


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

	s, err := cmd.createScribe()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
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

func (cmd *dockerAgent) createScribe() (*scribe.LogScribe, error) {
	if cmd.dbServer != "" && cmd.dbName != "" {
		gRPC, err := rpc.New(cmd.cert, cmd.key, cmd.ca)
		if err != nil {
			return nil, fmt.Errorf("failed to create a gRPC server: %v", err)
		}
		s, err := scribe.NewScribe(cmd.port, gRPC, true, cmd.dbServer, cmd.dbName, "", 0)
		if err != nil {
			return nil, fmt.Errorf("failed to create scribe: %v", err)
		}
		return s, nil
	}

	// validate path passed
	if err := scribe.CheckPath(cmd.rootPath); err != nil {
		return nil, fmt.Errorf("path passed is not valid: %v\n", err)
	}
	maxSize, err := scribe.LexicalToNumber(cmd.size)
	if err != nil {
		return nil, fmt.Errorf("couldn't parse size input to bytes: %v", err)
	}
	s, err := scribe.New(cmd.rootPath, cmd.port, maxSize, cmd.mediator, cmd.cert, cmd.key, cmd.ca)
	if err != nil {
		return nil, fmt.Errorf("failed to create scribe: %v", err)
	}
	return s, nil
}

func (cmd *dockerAgent) addMediator(id string) error {
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
