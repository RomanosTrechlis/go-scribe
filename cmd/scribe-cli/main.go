package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/RomanosTrechlis/go-icls/cli"
	pb "github.com/RomanosTrechlis/go-scribe/api"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"bytes"
	"text/tabwriter"
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
)

const (
	//HOST = "localhost"
	HOST = "192.168.99.100"
)

func main() {
	c := cli.New()
	version := c.New("version", versionShortDesc, versionLongDesc, getVersionHandler(c))
	version.BoolFlag("a", "all", "returns information from all the scribes", false)

	c.New("stats", statsShortDesc, statsLongDesc, getStatsHandler())

	c.New("resp", respShortDesc, respLongDesc, getRespHandler())

	if len(os.Args) == 1 {
		c.Execute("-h")
		return
	}

	cmd := strings.Join(os.Args[1:], " ")
	c.Execute(cmd)
}

func getRespHandler() func(flags map[string]string) error {
	return func (flags map[string]string) error {
		conn, err := grpc.Dial(HOST+":4242",
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

func getStatsHandler() func(flags map[string]string) error {
	return func (flags map[string]string) error {
		conn, err := grpc.Dial(HOST+":4242",
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

func getVersionHandler(c *cli.CLI) func(flags map[string]string) error {
	return func (flags map[string]string) error {
		a, err := c.BoolValue("a", "version", flags)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to get flag: %v", err)
			os.Exit(2)
		}
		conn, err := grpc.Dial(HOST+":4242",
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
