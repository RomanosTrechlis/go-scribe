package main

import (
	"os"

	"github.com/RomanosTrechlis/go-icls/cli"
)

var version = "undefined"

var c *cli.CLI

func init() {
	c = cli.New()
	agent := c.New("agent", agentShortHelp, agentLongHelp, agentHandler)
	agent.IntFlag("port", "", 8080, "port for server to listen to requests", false)
	agent.BoolFlag("pprof", "", "additional server for pprof functionality", false)
	agent.BoolFlag("console", "", "dumps log lines to console", false)
	agent.BoolFlag("verbose", "", "prints regular handled request count", false)
	agent.StringFlag("mediator", "", "", "mediators address if exists, i.e 127.0.0.1:8080", false)
	agent.IntFlag("pport", "", 1111, "port for pprof server", false)
	agent.StringFlag("path", "", "../../logs", "path for logs to be persisted", false)
	agent.StringFlag("size", "", "1MB", "max size for individual files, -1B for infinite size", false)
	agent.StringFlag("crt", "", "", "host's certificate for secured connections", false)
	agent.StringFlag("pk", "", "", "host's private key", false)
	agent.StringFlag("ca", "", "", "certificate authority's certificate", false)

	med := c.New("mediator", mediatorShortHelp, mediatorLongHelp, mediatorHandler)
	med.StringFlag("file", "", "", "configuration file path", false)
	med.IntFlag("port", "", 8000, "port for mediator server to listen to requests", false)
	med.BoolFlag("pprof", "", "additional server for pprof functionality", false)
	med.IntFlag("pport", "", 2222, "port for pprof server", false)
	med.StringFlag("crt", "", "", "host's certificate for secured connections", false)
	med.StringFlag("pk", "", "", "host's private key", false)
	med.StringFlag("ca", "", "", "certificate authority's certificate", false)
}

func main() {
	args := os.Args
	if len(args) == 1 {
		return
	}
	line := ""
	for _, a := range args[1:] {
		line += a + " "
	}
	c.Execute(line)
}
