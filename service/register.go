package service

import (
	"fmt"

	pb "github.com/RomanosTrechlis/go-scribe/api"
	p "github.com/RomanosTrechlis/go-scribe/internal/util/format/print"
	"golang.org/x/net/context"
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
