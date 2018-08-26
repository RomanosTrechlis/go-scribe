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
)

const (
	shortDesc = "version command returns the version of the mediator"
	longDesc  = `version command returns the version of the mediator.

It connects with gRPC to the mediator service and gets the version number.
`
)

const (
	//HOST = "localhost"
	HOST = "192.168.99.100"
)

func main() {
	c := cli.New()
	version := c.New("version", shortDesc, longDesc, func(flags map[string]string) error {
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
		fmt.Println(res.Version)
		return nil
	})
	version.BoolFlag("a", "all", "returns information from all the scribes", false)

	if len(os.Args) == 1 {
		c.Execute("-h")
		return
	}

	cmd := strings.Join(os.Args[1:], " ")
	c.Execute(cmd)
}
