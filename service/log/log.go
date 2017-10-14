package log

import (
	"golang.org/x/net/context"

	pb "github.com/RomanosTrechlis/logScribe/api"
	"google.golang.org/grpc"
)

// Logger contains the stream channel
type Logger struct {
	Stream chan pb.LogRequest
}

// Log is the ptotobuf service implementation
func (l Logger) Log(ctx context.Context, in *pb.LogRequest) (*pb.LogResponse, error) {
	l.Stream <- *in
	return &pb.LogResponse{Res: "true"}, nil
}

// GRPCService describes a method dealing with protobuf incoming requests
type GRPCService interface {
	serviceHandler(stop chan struct{}, s *grpc.Server)
}
