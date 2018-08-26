package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	pb "github.com/RomanosTrechlis/go-scribe/api"
	p "github.com/RomanosTrechlis/go-scribe/internal/util/format/print"
	"github.com/RomanosTrechlis/go-scribe/mediator"
	"github.com/RomanosTrechlis/go-scribe/scribe"
	"google.golang.org/grpc"
	"text/tabwriter"
	"bytes"
	"io"
)

type cliScribe struct {
	// should put Mediator and Scribe behind an interface
	isMediator bool
	mediator   *mediator.Mediator
	scribe     *scribe.LogScribe
}

func (cl cliScribe) GetVersion(ctx context.Context, in *pb.VersionRequest) (*pb.VersionResponse, error) {
	res := version
	if cl.isMediator && in.All {
		buf := new(bytes.Buffer)
		w := tabwriter.NewWriter(buf, 0, 0, 1, ' ', tabwriter.DiscardEmptyColumns)
		fmt.Fprint(w, "Type\tName\tVersion\n")
		fmt.Fprintf(w, "%s\t%s\t%s\n", "Mediator", "", version)
		cl.getVersionForScribes(w)
		w.Flush()
		res = string(buf.Bytes())
	}
	return &pb.VersionResponse{Version: res}, nil
}

func (cl cliScribe) getVersionForScribes(w io.Writer) {
	info := cl.mediator.GetInfo()
	for k, v := range info.Scribes {
		vr, err := getVersionFor(v)
		if err != nil {
			p.Print(fmt.Sprintf("failed to get version for %s: %v\n", k, err))
			continue
		}
		fmt.Fprintf(w, "%s\t%s\t%s\n", "Scribe", k, vr.Version)
	}
}

func getVersionFor(host string) (*pb.VersionResponse, error) {
	if strings.Contains(host, ":") {
		host = strings.Split(host, ":")[0]
	}
	conn, err := grpc.Dial(host+":4242",
		grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second))
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewCLIScribeClient(conn)
	return client.GetVersion(context.Background(), &pb.VersionRequest{})
}

func registerCLIScribeFunc(srv *grpc.Server, c cliScribe) func() {
	return func() {
		pb.RegisterCLIScribeServer(srv, c)
	}
}
