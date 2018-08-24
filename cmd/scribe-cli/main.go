package main

import (
	"github.com/RomanosTrechlis/go-icls/cli"
	pb "github.com/RomanosTrechlis/go-scribe/api"
	"google.golang.org/grpc"
	"time"
	"context"
	"fmt"
	"os"
	"strings"
)

const (
	shortDesc = "version command returns the version of the mediator"
	longDesc = `version command returns the version of the mediator.

It connects with gRPC to the mediator service and gets the version number.
`
)

func main() {
	c := cli.New()
	c.New("version", shortDesc, longDesc, func(flags map[string]string) error {
		conn, err := grpc.Dial("localhost:4242",
			grpc.WithInsecure(),
			grpc.WithTimeout(1*time.Second))
		if err != nil {
			return fmt.Errorf("did not connect: %v\n", err)
		}
		defer conn.Close()

		client := pb.NewCLIScribeClient(conn)
		fmt.Println(client.GetVersion(context.Background(), &pb.VersionRequest{}))
		return nil
	})

	if len(os.Args) == 1 {
		c.Execute("-h")
		return
	}

	cmd := strings.Join(os.Args[1:], " ")
	c.Execute(cmd)
}
