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
	"errors"
)

type cliScribe struct {
	// should put Mediator and Scribe behind an interface
	isMediator bool
	mediator   *mediator.Mediator
	scribe     *scribe.LogScribe
}

func (cl cliScribe) GetVersion(ctx context.Context, in *pb.VersionRequest) (*pb.VersionResponse, error) {
	if !cl.isMediator {
		res := &pb.Version{
			Type:    pb.Type_SCRIBE,
			Name:    "",
			Version: version,
		}
		return &pb.VersionResponse{
			Results: []*pb.Version{res},
		}, nil
	}

	response := &pb.VersionResponse{
		Results: make([]*pb.Version, 0),
	}
	if cl.isMediator {
		res := &pb.Version{
			Type:    pb.Type_MEDIATOR,
			Name:    "",
			Version: version,
		}
		response.Results = append(response.Results, res)
	}
	if in.All {
		response = cl.getVersionForScribes(response)
	}
	return response, nil
}

func (cl cliScribe) GetStats(ctx context.Context, in *pb.StatsRequest) (*pb.StatsResponse, error) {
	return nil, errors.New("RPC not implemented yet.")
}

func (cl cliScribe) GetScribesResponsibility(ctx context.Context, in *pb.ResponsibilityRequest) (*pb.ResponsibilityResponse, error) {
	return nil, errors.New("RPC not implemented yet.")
}

func (cl cliScribe) getVersionForScribes(resp *pb.VersionResponse) *pb.VersionResponse {
	info := cl.mediator.GetInfo()
	for k, v := range info.Scribes {
		vr, err := getVersionFor(v)
		if err != nil {
			p.Print(fmt.Sprintf("failed to get version for %s: %v\n", k, err))
			continue
		}
		version := &pb.Version{
			Type: pb.Type_SCRIBE,
			Name: k,
			Version: vr.GetResults()[0].GetVersion(),
		}
		resp.Results = append(resp.Results, version)
	}
	return resp
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
