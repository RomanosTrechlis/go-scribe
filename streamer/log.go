package streamer

import (
	"context"

	pb "github.com/RomanosTrechlis/logStreamer/api"
	"google.golang.org/grpc"
)

type logger struct {
	stream chan pb.LogRequest
}

// Log is the ptotobuf service implementation
func (l logger) Log(ctx context.Context, in *pb.LogRequest) (*pb.LogResponse, error) {
	l.stream <- *in
	return &pb.LogResponse{Res: "handling"}, nil
}

// GRPCService describes what to do with incoming protobuf requests
type GRPCService interface {
	ServiceHandler(stop chan struct{}, s *grpc.Server)
}
