package pinger

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"golang.org/x/net/context"

	pb "github.com/RomanosTrechlis/logStreamer/api"
	"google.golang.org/grpc"
)

// Pinger holds relevant info for pinging mediator
type Pinger struct {
	grpcAddr   string
	stopPinger chan struct{}
}

// New creates a new pinger
func New(addr string) *Pinger {
	return &Pinger{
		grpcAddr:   addr,
		stopPinger: make(chan struct{}, 1),
	}
}

// Ping pings mediator at fixed intervals
func (p *Pinger) Ping(id string, interval int) {
	conn, err := grpc.Dial(p.grpcAddr,
		grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	for _ = range time.Tick(time.Duration(interval) * time.Second) {
		select {
		case <-p.stopPinger:
			fmt.Printf("stopping pinger\n")
			return
		default:
			ping(conn, id)
		}
	}
}

func ping(conn *grpc.ClientConn, id string) {
	fmt.Println("pinging...")
	c := pb.NewPingerClient(conn)
	req := &pb.PingRequest{
		A:          rand.Int31(),
		B:          rand.Int31(),
		StreamerId: id,
	}
	r, err := c.Ping(context.Background(), req)
	if err != nil {
		log.Fatalf("failled: %v", err)
	}
	fmt.Printf("%d", r.Res)
}

// End stops pinger
func (p *Pinger) End() {
	p.stopPinger <- struct{}{}
}
