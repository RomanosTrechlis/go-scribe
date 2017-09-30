package mediator

import (
  "golang.org/x/net/context"
	"fmt"
	"time"

	pb "github.com/RomanosTrechlis/logStreamer/api"
	"github.com/RomanosTrechlis/logStreamer/util/gserver"
  "github.com/RomanosTrechlis/logStreamer/service"
  "github.com/RomanosTrechlis/logStreamer/util/format/time"
	"google.golang.org/grpc"
)

const logLayout string = "2006-01-02T15.04.05Z07.00"

// Mediator grpc server and other relative info
type Mediator struct {
  streamers map[string]*grpc.ClientConn

	// input stream of protobuf requests
	stream chan pb.LogRequest

	// this can be a composite with the corresponding fields of LogStreamer
	grpcServer *grpc.Server
	grpcPort   int
	stopGrpc   chan struct{}

	// these can also be composed
	startTime time.Time
	counter   int64
	stopAll   chan struct{}
}

// New creates a new mediator
func New(port int, crt, key, ca string) (*Mediator, error) {
	srv, err := gserver.New(logLayout, crt, key, ca)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc server: %v", err)
	}

	m := &Mediator{
    stream: make(chan pb.LogRequest),
    grpcServer: srv,
    grpcPort: port,
    streamers: make(map[string]*grpc.ClientConn),
  }
	return m, nil
}

// ServiceHandler implements the protobuf service
func (m *Mediator) ServiceHandler(stop chan struct{}) {
	for {
		select {
		case req := <-m.stream:
      // does something
      client := pb.NewLogStreamerClient(m.streamers["test"])
      client.Log(context.Background(), &req)
			m.counter++
		case <-stop:
			return
		}
	}
}

// Serve starts mediator server
func (m *Mediator) Serve() {
	fmt.Printf("%s [INFO] Log streamer is starting...\n", ftime.PrintTime(logLayout))
	m.stopAll = make(chan struct{})
	m.stopGrpc = make(chan struct{})
	m.startTime = time.Now()

  // todo: remove
  conn, err := grpc.Dial(":8080",
		grpc.WithInsecure(),
		grpc.WithTimeout(1*time.Second))
  if err != nil {
    return
  }
  m.streamers["test"] = conn

	// go func listens to stream and stop channels
	go m.ServiceHandler(m.stopGrpc)

	// rpc server
	go gserver.Serve(m.register(), fmt.Sprintf(":%d", m.grpcPort), m.grpcServer, logLayout)

	<-m.stopAll
}

func (m *Mediator) register() func() {
	return func() {
		med := &service.Pinger{
      Handle: func(id string) {fmt.Println(id)},
    }
		pb.RegisterPingerServer(m.grpcServer, med)
	}
}

// Shutdown gracefully stops mediator from serving
func (m *Mediator) Shutdown() {
	m.stopAll <- struct{}{}
	fmt.Printf("\n%s [INFO] initializing shut down, please wait.\n",
		ftime.PrintTime(logLayout))
	m.stopGrpc <- struct{}{}
	time.Sleep(1 * time.Second)
	fmt.Printf("%s [INFO] Log streamer handled %d requests during %v\n",
		ftime.PrintTime(logLayout), m.counter, time.Since(m.startTime))
	fmt.Printf("%s [INFO] Log streamer shut down\n", ftime.PrintTime(logLayout))
}
