package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	med "github.com/RomanosTrechlis/go-scribe/mediator"
	"github.com/RomanosTrechlis/go-scribe/profiling"
)

const (
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

type mediator struct {
	port      int
	pprofInfo bool
	pport     int

	cert string
	key  string
	ca   string
}

func (cmd *mediator) Name() string      { return "mediator" }
func (cmd *mediator) Args() string      { return "" }
func (cmd *mediator) ShortHelp() string { return agentShortHelp }
func (cmd *mediator) LongHelp() string  { return agentLongHelp }
func (cmd *mediator) Hidden() bool      { return false }

func (cmd *mediator) Register(fs *flag.FlagSet) {
	// rpc server listening port
	fs.IntVar(&cmd.port, "port", 8000, "port for mediator server to listen to requests")
	// enable/disable pprof functionality
	fs.BoolVar(&cmd.pprofInfo, "pprof", false, "additional server for pprof functionality")
	// pprof port for http server
	fs.IntVar(&cmd.pport, "pport", 1111, "port for pprof server")
	fs.StringVar(&cmd.cert, "crt", "", "host's certificate for secured connections")
	fs.StringVar(&cmd.key, "pk", "", "host's private key")
	fs.StringVar(&cmd.ca, "ca", "", "certificate authority's certificate")
}

func (cmd *mediator) Run(ctx *ctx, args []string) error {
	printLogo()

	// stopAll channel listens to termination and interupt signals.
	stopAll := make(chan os.Signal, 1)
	signal.Notify(stopAll, syscall.SIGTERM, syscall.SIGINT)

	m, err := med.New(cmd.port, cmd.cert, cmd.key, cmd.ca)
	if err != nil {
		return fmt.Errorf("failed to start a new mediator: %v", err)
	}
	defer m.Shutdown()
	go m.Serve()

	var srv *http.Server
	if cmd.pprofInfo {
		srv = profiling.Serve(cmd.pport)
		defer srv.Shutdown(nil)
	}
	<-stopAll
	return nil
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
