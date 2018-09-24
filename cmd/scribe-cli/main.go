package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/RomanosTrechlis/go-icls/cli"
	pb "github.com/RomanosTrechlis/go-scribe/api"
	"github.com/RomanosTrechlis/go-scribe/internal/util/fs"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"gopkg.in/yaml.v2"
)

const (
	versionShortDesc = "version command returns the version of the mediator"
	versionLongDesc  = `version command returns the version of the mediator.

It connects with gRPC to the mediator service and gets the version number.
`

	statsShortDesc = "stats command returns how many requests each scribe handled"
	statsLongDesc  = "stats command returns how many requests each scribe handled"

	respShortDesc = "resp command returns every scribe's filename responsibility"
	respLongDesc  = "resp command returns every scribe's filename responsibility"

	createShortDesc = "create command is used for creating config files"
	createLongDesc  = `create command is used as an interactive assistance for creating
various configuration files.
`
)

const (
	LOCALHOST  = "localhost"
	HELP_SHORT = "-h"
	HELP_LONG  = "--help"

	SCRIBE_DIR      = ".scribe"
	CLI_CONFIG      = "cli_config.yml"
	MEDIATOR_CONFIG = "mediator_config.yml"
	SCRIBE_CONFIG   = "scribe_config.yml"

	THE_ANSWER_TO_EVERYTHING = 4242
)

func main() {
	config := &CLIConfig{LOCALHOST}
	if getConfigFile(os.Args) {
		c, err := configFile()
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v", err)
			os.Exit(2)
		}
		config = c
	}

	c := createCommandTree(config.Host)

	if len(os.Args) == 1 {
		_, err := c.Execute(HELP_SHORT)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to execute command 'scribe-cli %s': %v", HELP_SHORT, err)
			os.Exit(2)
		}
		os.Exit(0)
	}

	cmd := strings.Join(os.Args[1:], " ")
	_, err := c.Execute(cmd)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to execute command 'scribe-cli %s': %v", cmd, err)
		os.Exit(2)
	}
	os.Exit(0)
}

func getConfigFile(args []string) bool {
	if len(args) < 2 {
		return false
	}
	if args[1] == "create" {
		return false
	}
	if args[1] == HELP_SHORT || args[1] == HELP_LONG {
		return false
	}
	return true
}

func configFile() (*CLIConfig, error) {
	homeDir, err := getUserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user's home directory: %v", err)
	}

	configFile := filepath.Join(homeDir, SCRIBE_DIR, CLI_CONFIG)
	b, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("config file doesn't exist, run 'scribe-cli create -t cli -w' to create it: %v", err)
	}

	config := new(CLIConfig)
	err = yaml.Unmarshal(b, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal file: %v", err)
	}
	return config, nil
}

func createCommandTree(host string) *cli.CLI {
	c := cli.New()
	version := c.New("version", versionShortDesc, versionLongDesc, getVersionHandler(host, c))
	version.BoolFlag("a", "all", "returns information from all the scribes", false)

	c.New("stats", statsShortDesc, statsLongDesc, getStatsHandler(host))

	c.New("resp", respShortDesc, respLongDesc, getRespHandler(host))

	create := c.New("create", createShortDesc, createLongDesc, getCreateHandler(c))
	create.StringFlag("t", "type", "cli", "prints the configuration on the stdout. Types: mediator, scribe, cli", true)
	create.BoolFlag("w", "write", "write creates the config file under .scribe directory", false)
	return c
}

func getCreateHandler(c *cli.CLI) func(flags map[string]string) error {
	return func(flags map[string]string) error {
		t := c.StringValue("t", "create", flags)
		w, err := c.BoolValue("w", "create", flags)
		if err != nil {
			w = false
		}

		switch t {
		case "cli":
			return createCLIConfig(w)
		case "mediator":
			return createMediatorConfig(w)
		case "scribe":
			return createAgentConfig(w)
		default:
			return fmt.Errorf("%s is not supported", t)
		}
		return nil
	}
}

func getUserHomeDir() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	return usr.HomeDir, nil
}

func createCLIConfig(write bool) error {
	homeDir, err := getUserHomeDir()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("What is mediator's IP address? ")
	scanner.Scan()
	addr := strings.Trim(scanner.Text(), " ")
	path := filepath.Join(homeDir, SCRIBE_DIR)
	err = fs.CreateFolderIfNotExist(path)
	if err != nil {
		return err
	}

	c := CLIConfig{addr}
	b, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	if !write {
		fmt.Fprint(os.Stdout, string(b))
		return nil
	}
	file := filepath.Join(path, CLI_CONFIG)
	err = ioutil.WriteFile(file, b, 0644)
	if err != nil {
		return err
	}
	return nil
}

func getRespHandler(host string) func(flags map[string]string) error {
	return func(flags map[string]string) error {
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, THE_ANSWER_TO_EVERYTHING),
			grpc.WithInsecure(),
			grpc.WithTimeout(1*time.Second))
		if err != nil {
			return fmt.Errorf("did not connect: %v\n", err)
		}
		defer conn.Close()

		client := pb.NewCLIScribeClient(conn)
		res, err := client.GetScribesResponsibility(context.Background(), &pb.ResponsibilityRequest{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to get response from mediator service: %v", err)
			os.Exit(2)
		}

		buf := new(bytes.Buffer)
		w := tabwriter.NewWriter(buf, 0, 0, 1, ' ', tabwriter.DiscardEmptyColumns)
		fmt.Fprint(w, "Name\tResponsibility\n")
		for _, v := range res.Result {
			fmt.Fprintf(w, "%s\t%s\n", v.Name, v.Responsibility)
		}
		w.Flush()
		fmt.Println(string(buf.Bytes()))
		return nil
	}
}

func getStatsHandler(host string) func(flags map[string]string) error {
	return func(flags map[string]string) error {
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, THE_ANSWER_TO_EVERYTHING),
			grpc.WithInsecure(),
			grpc.WithTimeout(1*time.Second))
		if err != nil {
			return fmt.Errorf("did not connect: %v\n", err)
		}
		defer conn.Close()

		client := pb.NewCLIScribeClient(conn)
		res, err := client.GetStats(context.Background(), &pb.StatsRequest{})
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to get response from mediator service: %v", err)
			os.Exit(2)
		}

		buf := new(bytes.Buffer)
		w := tabwriter.NewWriter(buf, 0, 0, 1, ' ', tabwriter.DiscardEmptyColumns)
		fmt.Fprint(w, "Name\tCount\n")
		for _, v := range res.Result {
			fmt.Fprintf(w, "%s\t%d\n", v.Name, v.Count)
		}
		w.Flush()
		fmt.Println(string(buf.Bytes()))
		return nil
	}
}

func getVersionHandler(host string, c *cli.CLI) func(flags map[string]string) error {
	return func(flags map[string]string) error {
		a, err := c.BoolValue("a", "version", flags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to get flag: %v", err)
			os.Exit(2)
		}
		conn, err := grpc.Dial(fmt.Sprintf("%s:%d", host, THE_ANSWER_TO_EVERYTHING),
			grpc.WithInsecure(),
			grpc.WithTimeout(1*time.Second))
		if err != nil {
			return fmt.Errorf("did not connect: %v\n", err)
		}
		defer conn.Close()

		client := pb.NewCLIScribeClient(conn)
		res, err := client.GetVersion(context.Background(), &pb.VersionRequest{All: a})
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to get response from mediator service: %v", err)
			os.Exit(2)
		}

		buf := new(bytes.Buffer)
		w := tabwriter.NewWriter(buf, 0, 0, 1, ' ', tabwriter.DiscardEmptyColumns)
		fmt.Fprint(w, "Type\tName\tVersion\n")
		for _, v := range res.Results {
			fmt.Fprintf(w, "%s\t%s\t%s\n", v.Type, v.Name, v.Version)
		}
		w.Flush()
		fmt.Println(string(buf.Bytes()))
		return nil
	}
}
