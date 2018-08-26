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
		res = "Type\tName\tVersion\n"
		res += fmt.Sprintf("%s\t%s\t%s\n", "Mediator", "", version)
		res += cl.getVersionForScribes()
	}
	return &pb.VersionResponse{Version: res}, nil
}

func (cl cliScribe) getVersionForScribes() string {
	res := ""
	info := cl.mediator.GetInfo()
	for k, v := range info.Scribes {
		vr, err := getVersionFor(v)
		if err != nil {
			p.Print(fmt.Sprintf("failed to get version for %s: %v\n", k, err))
			continue
		}
		res += fmt.Sprintf("%s\t%s\t%s\n", "Scribe", k, vr.Version)
	}
	return res
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
