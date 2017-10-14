package ping

import (
	"context"

	pb "github.com/RomanosTrechlis/logScribe/api"
)

// Pinger holds a function that deals with incoming pings
type Pinger struct{}

// Ping implements ping protobuf service
func (p *Pinger) Ping(ctx context.Context, req *pb.PingRequest) (*pb.PingResponse, error) {
	return &pb.PingResponse{
		Res: req.GetA() * req.GetB(),
	}, nil
}
