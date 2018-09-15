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
	"gopkg.in/yaml.v2"
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

type MediatorConfig struct {
	Port        int  `yaml:"port"`
	Profile     bool `yaml:"profile"`
	ProfilePort int  `yaml:"profile_port"`

	Certificate          string `yaml:"certificate"`
	PrivateKey           string `yaml:"private_key"`
	CertificateAuthority string `yaml:"certificate_authority"`
}

func fillMediatorConfig(flags map[string]string) (*MediatorConfig, error) {
	file := c.StringValue("file", "mediator", flags)
	if file != "" {
		return fillMediatorConfigFromFile(file)
	}
	return fillMediatorConfigFromFlags(flags)
}

func fillMediatorConfigFromFlags(flags map[string]string) (*MediatorConfig, error) {
	port, err := c.IntValue("port", "mediator", flags)
	if err != nil {
		return nil, fmt.Errorf("failed to get the value of 'port' flag: %v", err)
	}
	pport, err := c.IntValue("pport", "mediator", flags)
	if err != nil {
		return nil, fmt.Errorf("failed to get the value of 'pport' flag: %v", err)
	}
	pprofInfo, err := c.BoolValue("pprof", "mediator", flags)
	if err != nil {
		return nil, fmt.Errorf("failed to get the value of 'pprof' flag: %v", err)
	}
	crt := c.StringValue("crt", "mediator", flags)
	pk := c.StringValue("pk", "mediator", flags)
	ca := c.StringValue("ca", "mediator", flags)
	m := &MediatorConfig{port, pprofInfo, pport, crt, pk, ca}
	return m, nil
}

func fillMediatorConfigFromFile(file string) (*MediatorConfig, error) {
	b, err := readConfigurationFile(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read configuration file: %v", err)
	}

	m := new(MediatorConfig)
	yaml.Unmarshal(b, m)
	return m, nil
}

func mediatorHandler(flags map[string]string) error {
	conf, err := fillMediatorConfig(flags)
	if err != nil {
		return err
	}

	printLogo()

	// stopAll channel listens to termination and interrupt signals.
	stopAll := make(chan os.Signal, 1)
	signal.Notify(stopAll, syscall.SIGTERM, syscall.SIGINT)

	m, err := med.New(conf.Port, conf.Certificate, conf.PrivateKey, conf.CertificateAuthority)
	if err != nil {
		return fmt.Errorf("failed to start a new mediator: %v", err)
	}
	defer m.Shutdown()
	go m.Serve()

	var srv *http.Server

	if conf.Profile {
		srv = profiling.Serve(conf.ProfilePort)
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
