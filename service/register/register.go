package register

import (
	"fmt"

	"golang.org/x/net/context"

	pb "github.com/RomanosTrechlis/logStreamer/api"
	p "github.com/RomanosTrechlis/logStreamer/util/format/print"
)

// Register holds the subscribers
type Register struct {
	Subscribers map[string]string
}

// Register implements the corresponding protobuf service
func (r *Register) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	r.Subscribers[req.GetId()] = req.GetAddr()
	p.Print(fmt.Sprintf("Registering streamer %s from %s", req.GetId(), req.GetAddr()))
	return &pb.RegisterResponse{Res: "Success"}, nil
}
