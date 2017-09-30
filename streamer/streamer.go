package streamer

import (
	"fmt"
	"sync"
	"time"

	pb "github.com/RomanosTrechlis/logStreamer/api"
	"github.com/RomanosTrechlis/logStreamer/service"
	"github.com/RomanosTrechlis/logStreamer/util/format/time"
	"github.com/RomanosTrechlis/logStreamer/util/gserver"
	"google.golang.org/grpc"
)

const (
	logLayout string = "2006-01-02T15.04.05Z07.00"
	layout    string = "02012006150405"
)

// LogStreamer holds the servers and other relative information
type LogStreamer struct {
	// input stream of protobuf requests
	stream chan pb.LogRequest

	// grpcServer
	grpcServer *grpc.Server
	// GRPC server port
	grpcPort int
	// stopGrpc waits for an empty struct to stop the rpc server.
	stopGrpc chan struct{}

	// maximum log file size
	fileSize int64

	// ticker is responsible for printing
	// status updates every few seconds
	ticker     *time.Ticker
	stopTicker chan struct{}
	startTime  time.Time

	rootPath string
	counter  int64
	stopAll  chan struct{}
}

// New creates a streamer struct
func New(root string, port int, fileSize int64, crt, key, ca string) (*LogStreamer, error) {
	srv, err := gserver.New(logLayout, crt, key, ca)
	if err != nil {
		return nil, fmt.Errorf("failed to create grpc server: %v", err)
	}
	return &LogStreamer{
		stream:     make(chan pb.LogRequest),
		grpcServer: srv,
		grpcPort:   port,
		stopGrpc:   make(chan struct{}),
		fileSize:   fileSize,
		ticker:     time.NewTicker(20 * time.Second),
		stopTicker: make(chan struct{}),
		rootPath:   root,
	}, nil
}

// ServiceHandler implements the protobuf service
func (s *LogStreamer) ServiceHandler(stop chan struct{}) {
	for {
		select {
		case req := <-s.stream:
			s.counter++
			err := handleIncomingRequest(s.rootPath, req.GetPath(),
				req.GetFilename(), req.GetLine(), s.fileSize)
			if err != nil {
				fmt.Printf("hanldeIncomingRequest returned with error: %v", err)
				return
			}
		case <-stop:
			return
		}
	}
}

// Serve initializes log streamer's servers
func (s *LogStreamer) Serve() {
	fmt.Printf("%s [INFO] Log streamer is starting...\n", ftime.PrintTime(logLayout))
	s.stopAll = make(chan struct{})
	s.startTime = time.Now()
	// go func listens to stream and stop channels
	go s.ServiceHandler(s.stopGrpc)

	// rpc server
	go gserver.Serve(s.register(), fmt.Sprintf(":%d", s.grpcPort), s.grpcServer, logLayout)

	// ticker
	go s.tickerServ()
	<-s.stopAll
}

func (s *LogStreamer) register() func() {
	return func() {
		log := service.Logger{Stream: s.stream}
		pb.RegisterLogStreamerServer(s.grpcServer, log)
	}
}

// Shutdown gracefully stops log streamer from serving
func (s *LogStreamer) Shutdown() {
	s.stopAll <- struct{}{}
	fmt.Printf("\n%s [INFO] initializing shut down, please wait.\n",
		ftime.PrintTime(logLayout))
	s.stopGrpc <- struct{}{}
	s.stopTicker <- struct{}{}
	s.ticker.Stop()
	time.Sleep(1 * time.Second)
	fmt.Printf("%s [INFO] Log streamer handled %d requests during %v\n",
		ftime.PrintTime(logLayout), s.counter, time.Since(s.startTime))
	fmt.Printf("%s [INFO] Log streamer shut down\n", ftime.PrintTime(logLayout))
}

func (s *LogStreamer) tickerServ() {
	for _ = range s.ticker.C {
		select {
		case <-s.stopTicker:
			fmt.Printf("\n%s [INFO] Ticker is stopping...\n", ftime.PrintTime(logLayout))
			return
		default:
			fmt.Printf("%s [INFO] Log Streamer handled %d requests, so far.\n",
				ftime.PrintTime(logLayout), s.counter)
		}
	}
}

func handleIncomingRequest(rootPath, path, filename, line string, size int64) error {
	var mu sync.RWMutex
	mu.Lock()
	defer mu.Unlock()
	if err := writeLine(rootPath, path, filename, line, size); err != nil {
		return fmt.Errorf("failed to write line: %v", err)
	}
	return nil
}
