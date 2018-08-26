package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/RomanosTrechlis/go-scribe/internal/util/gserver"
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

func mediatorHandler(flags map[string]string) error {
	printLogo()

	// stopAll channel listens to termination and interrupt signals.
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
	pprofInfo, _ := c.BoolValue("pprof", "mediator", flags)
	if pprofInfo {
		srv = profiling.Serve(pport)
		defer srv.Shutdown(nil)
	}

	cliServer, _ := gserver.New("", "", "")
	c := cliScribe{true, m, nil}
	go gserver.Serve(registerCLIScribeFunc(cliServer, c), fmt.Sprint(":4242"), cliServer)

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
